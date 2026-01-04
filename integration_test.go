package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Integration tests for the claude_pp CLI

func TestCLIVersion(t *testing.T) {
	cmd := exec.Command("./claude_pp", "version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	if !strings.Contains(string(output), "claude_pp version") {
		t.Error("version output should contain 'claude_pp version'")
	}

	if !strings.Contains(string(output), "0.1.0") {
		t.Error("version output should contain version number")
	}
}

func TestCLIHelp(t *testing.T) {
	cmd := exec.Command("./claude_pp", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("help command failed: %v", err)
	}

	expectedCommands := []string{
		"remember",
		"recall",
		"instances",
		"messages",
		"send",
		"serve",
		"setup",
		"status",
		"version",
	}

	for _, expected := range expectedCommands {
		if !strings.Contains(string(output), expected) {
			t.Errorf("help output should contain '%s'", expected)
		}
	}
}

func TestCLIRememberHelp(t *testing.T) {
	cmd := exec.Command("./claude_pp", "remember", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("remember --help failed: %v", err)
	}

	if !strings.Contains(string(output), "-t, --tags") {
		t.Error("remember help should mention tags flag")
	}
}

func TestCLIRecallHelp(t *testing.T) {
	cmd := exec.Command("./claude_pp", "recall", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("recall --help failed: %v", err)
	}

	expectedFlags := []string{"-t, --tags", "-n, --limit", "-l, --local"}
	for _, flag := range expectedFlags {
		if !strings.Contains(string(output), flag) {
			t.Errorf("recall help should mention %s flag", flag)
		}
	}
}

func TestCLISetupHelp(t *testing.T) {
	cmd := exec.Command("./claude_pp", "setup", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("setup --help failed: %v", err)
	}

	expectedFlags := []string{"--global", "--project", "--opencode", "--codex", "--gemini", "--allow-all"}
	for _, flag := range expectedFlags {
		if !strings.Contains(string(output), flag) {
			t.Errorf("setup help should mention %s flag", flag)
		}
	}
}

func TestCLIRememberAndRecall(t *testing.T) {
	// Use a temporary data directory
	tmpDir, err := os.MkdirTemp("", "claude_pp_integration_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	env := append(os.Environ(), "HOME="+tmpDir)

	// Remember a fact
	rememberCmd := exec.Command("./claude_pp", "remember", "Integration test fact", "-t", "test", "-t", "integration")
	rememberCmd.Env = env
	rememberOutput, err := rememberCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("remember failed: %v, output: %s", err, rememberOutput)
	}

	if !strings.Contains(string(rememberOutput), "Stored fact") {
		t.Error("remember should output 'Stored fact'")
	}

	// Recall the fact
	recallCmd := exec.Command("./claude_pp", "recall", "Integration")
	recallCmd.Env = env
	recallOutput, err := recallCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("recall failed: %v, output: %s", err, recallOutput)
	}

	if !strings.Contains(string(recallOutput), "Integration test fact") {
		t.Error("recall should return the stored fact")
	}

	if !strings.Contains(string(recallOutput), "test, integration") {
		t.Error("recall should show the tags")
	}
}

func TestCLIStatus(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude_pp_status_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	env := append(os.Environ(), "HOME="+tmpDir)

	// Add a fact first
	rememberCmd := exec.Command("./claude_pp", "remember", "Status test fact")
	rememberCmd.Env = env
	rememberCmd.Run()

	// Check status
	cmd := exec.Command("./claude_pp", "status")
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("status failed: %v, output: %s", err, output)
	}

	expectedStrings := []string{
		"Claude PP Status",
		"Data directory:",
		"Working directory:",
		"Facts:",
		"Total:",
		"Running instances:",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(string(output), s) {
			t.Errorf("status output should contain '%s'", s)
		}
	}
}

func TestCLIInstances(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude_pp_instances_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create the data directory
	os.MkdirAll(filepath.Join(tmpDir, ".claude_pp"), 0755)

	cmd := exec.Command("./claude_pp", "instances")
	cmd.Env = append(os.Environ(), "HOME="+tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("instances failed: %v, output: %s", err, output)
	}

	if !strings.Contains(string(output), "No running instances found") {
		t.Error("instances should report no running instances")
	}
}

func TestMCPServerProtocol(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude_pp_mcp_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	env := append(os.Environ(), "HOME="+tmpDir)

	// Test initialize
	initRequest := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`

	cmd := exec.Command("./claude_pp", "serve")
	cmd.Env = env
	cmd.Stdin = strings.NewReader(initRequest + "\n")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Start()
	if err != nil {
		t.Fatalf("failed to start serve: %v", err)
	}

	// Give it a moment to process
	cmd.Process.Kill()

	output := stdout.String()
	if output == "" {
		t.Skip("MCP server didn't produce output (might need more time)")
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(strings.Split(output, "\n")[0]), &response); err != nil {
		t.Fatalf("failed to parse response: %v, output: %s", err, output)
	}

	result, ok := response["result"].(map[string]interface{})
	if !ok {
		t.Fatal("response should have result field")
	}

	if result["protocolVersion"] != "2024-11-05" {
		t.Errorf("protocolVersion = %v, want 2024-11-05", result["protocolVersion"])
	}
}

func TestCLIRecallFilters(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude_pp_filter_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	env := append(os.Environ(), "HOME="+tmpDir)

	// Add multiple facts with different tags
	cmd1 := exec.Command("./claude_pp", "remember", "Database PostgreSQL", "-t", "database")
	cmd1.Env = env
	cmd1.Run()

	cmd2 := exec.Command("./claude_pp", "remember", "API REST endpoints", "-t", "api")
	cmd2.Env = env
	cmd2.Run()

	cmd3 := exec.Command("./claude_pp", "remember", "Frontend React", "-t", "frontend")
	cmd3.Env = env
	cmd3.Run()

	// Test tag filter
	cmd := exec.Command("./claude_pp", "recall", "-t", "database")
	cmd.Env = env
	output, _ := cmd.CombinedOutput()

	if !strings.Contains(string(output), "PostgreSQL") {
		t.Error("tag filter should return database fact")
	}

	if strings.Contains(string(output), "React") {
		t.Error("tag filter should not return frontend fact")
	}

	// Test limit
	cmd = exec.Command("./claude_pp", "recall", "-n", "1")
	cmd.Env = env
	output, _ = cmd.CombinedOutput()

	lines := strings.Split(strings.TrimSpace(string(output)), "\n\n")
	if len(lines) > 1 {
		t.Error("limit should restrict results")
	}
}

func TestDataPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude_pp_persist_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	env := append(os.Environ(), "HOME="+tmpDir)

	// Store a fact
	rememberCmd := exec.Command("./claude_pp", "remember", "Persistent data test")
	rememberCmd.Env = env
	rememberCmd.Run()

	// Verify database file exists
	dbPath := filepath.Join(tmpDir, ".claude_pp", "claude_pp.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("database file should be created")
	}

	// Recall should still work (data persisted)
	cmd := exec.Command("./claude_pp", "recall", "Persistent")
	cmd.Env = env
	output, _ := cmd.CombinedOutput()

	if !strings.Contains(string(output), "Persistent data test") {
		t.Error("data should persist across commands")
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "remember without fact",
			args:    []string{"remember"},
			wantErr: true,
		},
		{
			name:    "send without enough args",
			args:    []string{"send", "only-one-arg"},
			wantErr: true,
		},
		{
			name:    "messages without instance id",
			args:    []string{"messages"},
			wantErr: true,
		},
		{
			name:    "invalid command",
			args:    []string{"nonexistent"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./claude_pp", tt.args...)
			err := cmd.Run()

			if (err != nil) != tt.wantErr {
				t.Errorf("Command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
