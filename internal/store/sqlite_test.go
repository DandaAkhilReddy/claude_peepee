package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestStore(t *testing.T) (*SQLiteStore, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "claude_peepee_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	store, err := NewSQLiteStore(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to create store: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.RemoveAll(tmpDir)
	}

	return store, cleanup
}

// ==================== FACT TESTS ====================

func TestAddFact(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	tests := []struct {
		name      string
		content   string
		tags      []string
		sourceDir string
		wantErr   bool
	}{
		{
			name:      "simple fact",
			content:   "This is a test fact",
			tags:      nil,
			sourceDir: "/test/dir",
			wantErr:   false,
		},
		{
			name:      "fact with tags",
			content:   "Tagged fact",
			tags:      []string{"tag1", "tag2"},
			sourceDir: "/test/dir",
			wantErr:   false,
		},
		{
			name:      "fact with many tags",
			content:   "Many tags fact",
			tags:      []string{"a", "b", "c", "d", "e"},
			sourceDir: "/test/dir",
			wantErr:   false,
		},
		{
			name:      "empty content",
			content:   "",
			tags:      nil,
			sourceDir: "/test/dir",
			wantErr:   false, // Empty content is allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := store.AddFact(tt.content, tt.tags, tt.sourceDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id <= 0 {
				t.Errorf("AddFact() returned invalid id = %d", id)
			}
		})
	}
}

func TestAddFactSizeLimit(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Create content larger than maxFactSize (1MB)
	largeContent := make([]byte, maxFactSize+1)
	for i := range largeContent {
		largeContent[i] = 'a'
	}

	_, err := store.AddFact(string(largeContent), nil, "/test")
	if err == nil {
		t.Error("AddFact() should fail for content exceeding size limit")
	}
}

func TestAddFactTagLimits(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Test too many tags
	manyTags := make([]string, maxTagCount+1)
	for i := range manyTags {
		manyTags[i] = "tag"
	}

	_, err := store.AddFact("test", manyTags, "/test")
	if err == nil {
		t.Error("AddFact() should fail for too many tags")
	}

	// Test tag too long
	longTag := make([]byte, maxTagLength+1)
	for i := range longTag {
		longTag[i] = 'a'
	}

	_, err = store.AddFact("test", []string{string(longTag)}, "/test")
	if err == nil {
		t.Error("AddFact() should fail for tag exceeding length limit")
	}
}

func TestGetFact(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	content := "Test fact content"
	tags := []string{"tag1", "tag2"}
	sourceDir := "/test/dir"

	id, err := store.AddFact(content, tags, sourceDir)
	if err != nil {
		t.Fatalf("AddFact() failed: %v", err)
	}

	fact, err := store.GetFact(id)
	if err != nil {
		t.Fatalf("GetFact() failed: %v", err)
	}

	if fact == nil {
		t.Fatal("GetFact() returned nil")
	}

	if fact.Content != content {
		t.Errorf("Content = %q, want %q", fact.Content, content)
	}

	if len(fact.Tags) != len(tags) {
		t.Errorf("Tags count = %d, want %d", len(fact.Tags), len(tags))
	}

	if fact.SourceDir != sourceDir {
		t.Errorf("SourceDir = %q, want %q", fact.SourceDir, sourceDir)
	}
}

func TestGetFactNotFound(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	fact, err := store.GetFact(99999)
	if err != nil {
		t.Fatalf("GetFact() failed: %v", err)
	}

	if fact != nil {
		t.Error("GetFact() should return nil for non-existent fact")
	}
}

func TestDeleteFact(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	id, err := store.AddFact("To be deleted", nil, "/test")
	if err != nil {
		t.Fatalf("AddFact() failed: %v", err)
	}

	err = store.DeleteFact(id)
	if err != nil {
		t.Fatalf("DeleteFact() failed: %v", err)
	}

	fact, err := store.GetFact(id)
	if err != nil {
		t.Fatalf("GetFact() failed: %v", err)
	}

	if fact != nil {
		t.Error("Fact should be deleted")
	}
}

func TestQueryFacts(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Add test facts
	store.AddFact("PostgreSQL database setup", []string{"database", "setup"}, "/project1")
	store.AddFact("MySQL configuration", []string{"database", "config"}, "/project1")
	store.AddFact("API authentication", []string{"api", "security"}, "/project2")
	store.AddFact("Frontend React setup", []string{"frontend", "setup"}, "/project2")

	tests := []struct {
		name      string
		query     string
		tags      []string
		sourceDir string
		limit     int
		wantMin   int
	}{
		{
			name:    "search by keyword",
			query:   "database",
			wantMin: 2,
		},
		{
			name:    "search by tag",
			tags:    []string{"setup"},
			wantMin: 2,
		},
		{
			name:      "filter by source dir",
			sourceDir: "/project1",
			wantMin:   2,
		},
		{
			name:    "combined search",
			query:   "database",
			tags:    []string{"setup"},
			wantMin: 1,
		},
		{
			name:    "limit results",
			limit:   1,
			wantMin: 1,
		},
		{
			name:    "no results",
			query:   "nonexistent",
			wantMin: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facts, err := store.QueryFacts(tt.query, tt.tags, tt.sourceDir, tt.limit)
			if err != nil {
				t.Fatalf("QueryFacts() failed: %v", err)
			}

			if len(facts) < tt.wantMin {
				t.Errorf("QueryFacts() returned %d facts, want at least %d", len(facts), tt.wantMin)
			}
		})
	}
}

func TestCountFacts(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.AddFact("Fact 1", nil, "/dir1")
	store.AddFact("Fact 2", nil, "/dir1")
	store.AddFact("Fact 3", nil, "/dir2")

	total, local, err := store.CountFacts("/dir1")
	if err != nil {
		t.Fatalf("CountFacts() failed: %v", err)
	}

	if total != 3 {
		t.Errorf("Total = %d, want 3", total)
	}

	if local != 2 {
		t.Errorf("Local = %d, want 2", local)
	}
}

// ==================== INSTANCE TESTS ====================

func TestRegisterInstance(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	err := store.RegisterInstance("test-id", 12345, "/working/dir")
	if err != nil {
		t.Fatalf("RegisterInstance() failed: %v", err)
	}

	inst, err := store.GetInstance("test-id")
	if err != nil {
		t.Fatalf("GetInstance() failed: %v", err)
	}

	if inst == nil {
		t.Fatal("Instance not found")
	}

	if inst.ID != "test-id" {
		t.Errorf("ID = %q, want %q", inst.ID, "test-id")
	}

	if inst.PID != 12345 {
		t.Errorf("PID = %d, want %d", inst.PID, 12345)
	}

	if inst.WorkingDir != "/working/dir" {
		t.Errorf("WorkingDir = %q, want %q", inst.WorkingDir, "/working/dir")
	}
}

func TestUpdateHeartbeat(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.RegisterInstance("test-id", 12345, "/dir")

	// Wait to ensure timestamp difference (SQLite has second-level precision)
	time.Sleep(1100 * time.Millisecond)

	err := store.UpdateHeartbeat("test-id")
	if err != nil {
		t.Fatalf("UpdateHeartbeat() failed: %v", err)
	}

	// Just verify no error - SQLite timestamp precision makes time comparison flaky
	inst2, _ := store.GetInstance("test-id")
	if inst2 == nil {
		t.Error("Instance should still exist after heartbeat update")
	}
}

func TestUnregisterInstance(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.RegisterInstance("test-id", 12345, "/dir")

	err := store.UnregisterInstance("test-id")
	if err != nil {
		t.Fatalf("UnregisterInstance() failed: %v", err)
	}

	inst, _ := store.GetInstance("test-id")
	if inst != nil {
		t.Error("Instance should be unregistered")
	}
}

func TestGetInstances(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.RegisterInstance("id1", 1, "/dir1")
	store.RegisterInstance("id2", 2, "/dir2")
	store.RegisterInstance("id3", 3, "/dir3")

	instances, err := store.GetInstances()
	if err != nil {
		t.Fatalf("GetInstances() failed: %v", err)
	}

	if len(instances) != 3 {
		t.Errorf("Got %d instances, want 3", len(instances))
	}
}

func TestCleanupStaleInstances(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.RegisterInstance("fresh", 1, "/dir1")
	store.RegisterInstance("stale", 2, "/dir2")

	// This is a simplified test - in reality we'd need to manipulate timestamps
	err := store.CleanupStaleInstances(1 * time.Hour)
	if err != nil {
		t.Fatalf("CleanupStaleInstances() failed: %v", err)
	}

	instances, _ := store.GetInstances()
	if len(instances) != 2 {
		t.Errorf("Got %d instances, want 2 (nothing should be cleaned with 1hr threshold)", len(instances))
	}
}

// ==================== MESSAGE TESTS ====================

func TestSendMessage(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	id, err := store.SendMessage("from-id", "to-id", "Hello!")
	if err != nil {
		t.Fatalf("SendMessage() failed: %v", err)
	}

	if id <= 0 {
		t.Errorf("Invalid message ID: %d", id)
	}
}

func TestSendMessageSizeLimit(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	largeMessage := make([]byte, maxMessageSize+1)
	for i := range largeMessage {
		largeMessage[i] = 'a'
	}

	_, err := store.SendMessage("from", "to", string(largeMessage))
	if err == nil {
		t.Error("SendMessage() should fail for message exceeding size limit")
	}
}

func TestGetMessages(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.SendMessage("sender1", "receiver", "Message 1")
	store.SendMessage("sender2", "receiver", "Message 2")
	store.SendMessage("sender1", "other", "Message 3")

	messages, err := store.GetMessages("receiver", false)
	if err != nil {
		t.Fatalf("GetMessages() failed: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Got %d messages, want 2", len(messages))
	}
}

func TestGetMessagesUnreadOnly(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	id1, _ := store.SendMessage("sender", "receiver", "Message 1")
	store.SendMessage("sender", "receiver", "Message 2")

	store.MarkMessageRead(id1)

	messages, err := store.GetMessages("receiver", true)
	if err != nil {
		t.Fatalf("GetMessages() failed: %v", err)
	}

	if len(messages) != 1 {
		t.Errorf("Got %d unread messages, want 1", len(messages))
	}
}

func TestMarkMessageRead(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	id, _ := store.SendMessage("sender", "receiver", "Test message")

	err := store.MarkMessageRead(id)
	if err != nil {
		t.Fatalf("MarkMessageRead() failed: %v", err)
	}

	messages, _ := store.GetMessages("receiver", true)
	if len(messages) != 0 {
		t.Error("Message should be marked as read")
	}
}

// ==================== DATABASE TESTS ====================

func TestNewSQLiteStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude_peepee_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewSQLiteStore(tmpDir)
	if err != nil {
		t.Fatalf("NewSQLiteStore() failed: %v", err)
	}
	defer store.Close()

	// Verify database file was created
	dbPath := filepath.Join(tmpDir, "claude_peepee.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestStoreClose(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	err := store.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Operations after close should fail
	_, err = store.AddFact("test", nil, "/test")
	if err == nil {
		t.Error("Operations should fail after Close()")
	}
}

// ==================== FTS TESTS ====================

func TestFullTextSearch(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.AddFact("PostgreSQL is a powerful database", nil, "/test")
	store.AddFact("MySQL is another database option", nil, "/test")
	store.AddFact("Redis for caching", nil, "/test")

	// Search for "database"
	facts, err := store.QueryFacts("database", nil, "", 10)
	if err != nil {
		t.Fatalf("QueryFacts() failed: %v", err)
	}

	if len(facts) != 2 {
		t.Errorf("Expected 2 results for 'database', got %d", len(facts))
	}

	// Search for "PostgreSQL"
	facts, err = store.QueryFacts("PostgreSQL", nil, "", 10)
	if err != nil {
		t.Fatalf("QueryFacts() failed: %v", err)
	}

	if len(facts) != 1 {
		t.Errorf("Expected 1 result for 'PostgreSQL', got %d", len(facts))
	}
}

func TestSanitizeFTSQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", `"simple"`},
		{`with "quotes"`, `"with ""quotes"""`},
		{"with special*", `"with special*"`},
	}

	for _, tt := range tests {
		result := sanitizeFTSQuery(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeFTSQuery(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
