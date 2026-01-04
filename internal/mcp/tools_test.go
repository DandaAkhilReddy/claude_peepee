package mcp

import (
	"os"
	"testing"
	"time"

	"github.com/DandaAkhilReddy/claude_pp/internal/store"
)

func setupTestServer(t *testing.T) (*Server, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "claude_pp_mcp_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	s, err := store.NewSQLiteStore(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to create store: %v", err)
	}

	server := NewServer(s, "test-instance", "/test/dir")

	cleanup := func() {
		s.Close()
		os.RemoveAll(tmpDir)
	}

	return server, cleanup
}

// ==================== REMEMBER TOOL TESTS ====================

func TestHandleRemember(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "simple fact",
			args: map[string]interface{}{
				"fact": "Test fact content",
			},
			wantErr: false,
		},
		{
			name: "fact with tags",
			args: map[string]interface{}{
				"fact": "Tagged fact",
				"tags": []interface{}{"tag1", "tag2"},
			},
			wantErr: false,
		},
		{
			name:    "missing fact",
			args:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "empty fact",
			args: map[string]interface{}{
				"fact": "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.handleRemember(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleRemember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				res := result.(map[string]interface{})
				if !res["success"].(bool) {
					t.Error("Expected success = true")
				}
				if res["id"].(int64) <= 0 {
					t.Error("Expected valid ID")
				}
			}
		})
	}
}

func TestHandleRememberSizeLimit(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	largeFact := make([]byte, maxFactSize+1)
	for i := range largeFact {
		largeFact[i] = 'a'
	}

	_, err := server.handleRemember(map[string]interface{}{
		"fact": string(largeFact),
	})
	if err == nil {
		t.Error("Expected error for oversized fact")
	}
}

func TestHandleRememberTagLimits(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Too many tags
	manyTags := make([]interface{}, maxTagCount+1)
	for i := range manyTags {
		manyTags[i] = "tag"
	}

	_, err := server.handleRemember(map[string]interface{}{
		"fact": "test",
		"tags": manyTags,
	})
	if err == nil {
		t.Error("Expected error for too many tags")
	}

	// Tag too long
	longTag := make([]byte, maxTagLength+1)
	for i := range longTag {
		longTag[i] = 'a'
	}

	_, err = server.handleRemember(map[string]interface{}{
		"fact": "test",
		"tags": []interface{}{string(longTag)},
	})
	if err == nil {
		t.Error("Expected error for oversized tag")
	}
}

// ==================== RECALL TOOL TESTS ====================

func TestHandleRecall(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Add some test facts
	server.handleRemember(map[string]interface{}{
		"fact": "PostgreSQL database",
		"tags": []interface{}{"database"},
	})
	server.handleRemember(map[string]interface{}{
		"fact": "Redis caching",
		"tags": []interface{}{"cache"},
	})

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantMin int
		wantErr bool
	}{
		{
			name:    "search by query",
			args:    map[string]interface{}{"query": "database"},
			wantMin: 1,
			wantErr: false,
		},
		{
			name:    "search by tag",
			args:    map[string]interface{}{"tags": []interface{}{"database"}},
			wantMin: 1,
			wantErr: false,
		},
		{
			name:    "empty search returns all",
			args:    map[string]interface{}{},
			wantMin: 2,
			wantErr: false,
		},
		{
			name:    "with limit",
			args:    map[string]interface{}{"limit": float64(1)},
			wantMin: 1,
			wantErr: false,
		},
		{
			name:    "no matches",
			args:    map[string]interface{}{"query": "nonexistent"},
			wantMin: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.handleRecall(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleRecall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				res := result.(map[string]interface{})
				count := res["count"].(int)
				if count < tt.wantMin {
					t.Errorf("Got %d facts, want at least %d", count, tt.wantMin)
				}
			}
		})
	}
}

// ==================== GET_CONTEXT TOOL TESTS ====================

func TestHandleGetContext(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Add facts
	server.handleRemember(map[string]interface{}{"fact": "Local fact 1"})
	server.handleRemember(map[string]interface{}{"fact": "Local fact 2"})

	result, err := server.handleGetContext(map[string]interface{}{})
	if err != nil {
		t.Fatalf("handleGetContext() error = %v", err)
	}

	res := result.(map[string]interface{})

	if res["working_dir"] != "/test/dir" {
		t.Errorf("working_dir = %v, want /test/dir", res["working_dir"])
	}

	if res["instance_id"] != "test-instance" {
		t.Errorf("instance_id = %v, want test-instance", res["instance_id"])
	}

	localFacts := res["local_facts"].([]map[string]interface{})
	if len(localFacts) != 2 {
		t.Errorf("Got %d local facts, want 2", len(localFacts))
	}
}

// ==================== LIST_INSTANCES TOOL TESTS ====================

func TestHandleListInstances(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Register the server's instance
	server.store.(*store.SQLiteStore).RegisterInstance("test-instance", 12345, "/test/dir")

	result, err := server.handleListInstances(map[string]interface{}{})
	if err != nil {
		t.Fatalf("handleListInstances() error = %v", err)
	}

	res := result.(map[string]interface{})

	if res["current_id"] != "test-instance" {
		t.Errorf("current_id = %v, want test-instance", res["current_id"])
	}

	instances := res["instances"].([]map[string]interface{})
	if len(instances) < 1 {
		t.Error("Expected at least one instance")
	}
}

// ==================== SEND_MESSAGE TOOL TESTS ====================

func TestHandleSendMessage(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Register a target instance
	server.store.(*store.SQLiteStore).RegisterInstance("target-instance", 54321, "/target/dir")

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid message",
			args: map[string]interface{}{
				"to_instance": "target-instance",
				"message":     "Hello!",
			},
			wantErr: false,
		},
		{
			name: "missing to_instance",
			args: map[string]interface{}{
				"message": "Hello!",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			args: map[string]interface{}{
				"to_instance": "target-instance",
			},
			wantErr: true,
		},
		{
			name: "non-existent instance",
			args: map[string]interface{}{
				"to_instance": "nonexistent",
				"message":     "Hello!",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.handleSendMessage(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleSendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				res := result.(map[string]interface{})
				if !res["success"].(bool) {
					t.Error("Expected success = true")
				}
			}
		})
	}
}

func TestHandleSendMessageSizeLimit(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	server.store.(*store.SQLiteStore).RegisterInstance("target", 1, "/dir")

	largeMessage := make([]byte, maxMessageSize+1)
	for i := range largeMessage {
		largeMessage[i] = 'a'
	}

	_, err := server.handleSendMessage(map[string]interface{}{
		"to_instance": "target",
		"message":     string(largeMessage),
	})
	if err == nil {
		t.Error("Expected error for oversized message")
	}
}

// ==================== GET_MESSAGES TOOL TESTS ====================

func TestHandleGetMessages(t *testing.T) {
	// Test unread messages
	t.Run("get unread messages", func(t *testing.T) {
		server, cleanup := setupTestServer(t)
		defer cleanup()

		server.store.(*store.SQLiteStore).SendMessage("sender1", "test-instance", "Message 1")
		server.store.(*store.SQLiteStore).SendMessage("sender2", "test-instance", "Message 2")

		result, err := server.handleGetMessages(map[string]interface{}{"unread_only": true})
		if err != nil {
			t.Fatalf("handleGetMessages() error = %v", err)
		}

		res := result.(map[string]interface{})
		if res["count"].(int) != 2 {
			t.Errorf("Got %d messages, want 2", res["count"].(int))
		}
	})

	// Test all messages (including read)
	t.Run("get all messages", func(t *testing.T) {
		server, cleanup := setupTestServer(t)
		defer cleanup()

		server.store.(*store.SQLiteStore).SendMessage("sender1", "test-instance", "Message 1")
		server.store.(*store.SQLiteStore).SendMessage("sender2", "test-instance", "Message 2")

		// First read them (marks as read)
		server.handleGetMessages(map[string]interface{}{"unread_only": true})

		// Now get all (should still return 2)
		result, err := server.handleGetMessages(map[string]interface{}{"unread_only": false})
		if err != nil {
			t.Fatalf("handleGetMessages() error = %v", err)
		}

		res := result.(map[string]interface{})
		if res["count"].(int) != 2 {
			t.Errorf("Got %d messages, want 2", res["count"].(int))
		}
	})
}

func TestGetMessagesMarksAsRead(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	server.store.(*store.SQLiteStore).SendMessage("sender", "test-instance", "Test message")

	// First call should return the message
	result1, _ := server.handleGetMessages(map[string]interface{}{"unread_only": true})
	res1 := result1.(map[string]interface{})
	if res1["count"].(int) != 1 {
		t.Error("Expected 1 unread message initially")
	}

	// Second call should return 0 (message was marked as read)
	result2, _ := server.handleGetMessages(map[string]interface{}{"unread_only": true})
	res2 := result2.(map[string]interface{})
	if res2["count"].(int) != 0 {
		t.Error("Expected 0 unread messages after first retrieval")
	}
}

// ==================== INTEGRATION TESTS ====================

func TestToolsWorkflow(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Remember some facts
	_, err := server.handleRemember(map[string]interface{}{
		"fact": "Project uses PostgreSQL",
		"tags": []interface{}{"database", "architecture"},
	})
	if err != nil {
		t.Fatalf("Remember failed: %v", err)
	}

	// 2. Recall the fact
	recallResult, err := server.handleRecall(map[string]interface{}{
		"query": "PostgreSQL",
	})
	if err != nil {
		t.Fatalf("Recall failed: %v", err)
	}
	recallRes := recallResult.(map[string]interface{})
	if recallRes["count"].(int) != 1 {
		t.Error("Expected to recall 1 fact")
	}

	// 3. Get context
	contextResult, err := server.handleGetContext(map[string]interface{}{})
	if err != nil {
		t.Fatalf("GetContext failed: %v", err)
	}
	contextRes := contextResult.(map[string]interface{})
	if contextRes["local_count"].(int) != 1 {
		t.Error("Expected 1 local fact in context")
	}
}

func TestInstanceCommunicationWorkflow(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// Register instances
	server.store.(*store.SQLiteStore).RegisterInstance("instance-a", 1, "/dir/a")
	server.store.(*store.SQLiteStore).RegisterInstance("instance-b", 2, "/dir/b")

	// List instances
	listResult, err := server.handleListInstances(map[string]interface{}{})
	if err != nil {
		t.Fatalf("ListInstances failed: %v", err)
	}
	listRes := listResult.(map[string]interface{})
	if listRes["count"].(int) < 2 {
		t.Error("Expected at least 2 instances")
	}

	// Send message from server to instance-b
	_, err = server.handleSendMessage(map[string]interface{}{
		"to_instance": "instance-b",
		"message":     "Hello from test!",
	})
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	// Create a server for instance-b to receive the message
	serverB := NewServer(server.store, "instance-b", "/dir/b")

	msgResult, err := serverB.handleGetMessages(map[string]interface{}{})
	if err != nil {
		t.Fatalf("GetMessages failed: %v", err)
	}
	msgRes := msgResult.(map[string]interface{})
	if msgRes["count"].(int) != 1 {
		t.Errorf("Expected 1 message, got %d", msgRes["count"].(int))
	}
}

// ==================== CONSTANTS TESTS ====================

func TestConstants(t *testing.T) {
	if maxFactSize != 1024*1024 {
		t.Errorf("maxFactSize = %d, want 1MB", maxFactSize)
	}

	if maxMessageSize != 64*1024 {
		t.Errorf("maxMessageSize = %d, want 64KB", maxMessageSize)
	}

	if maxTagLength != 100 {
		t.Errorf("maxTagLength = %d, want 100", maxTagLength)
	}

	if maxTagCount != 50 {
		t.Errorf("maxTagCount = %d, want 50", maxTagCount)
	}

	if staleTimeout != 5*time.Minute {
		t.Errorf("staleTimeout = %v, want 5 minutes", staleTimeout)
	}
}
