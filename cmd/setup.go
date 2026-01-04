package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var (
	setupGlobal   bool
	setupProject  bool
	setupOpenCode bool
	setupCodex    bool
	setupGemini   bool
	setupAllowAll bool
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure claude_peepee for your AI coding tools",
	Long: `Configure claude_peepee as an MCP server for Claude Code, OpenCode, Codex CLI, or Gemini CLI.

By default, configures Claude Code globally. Use flags to configure other tools or project-specific settings.`,
	RunE: runSetup,
}

func init() {
	setupCmd.Flags().BoolVarP(&setupGlobal, "global", "g", false, "Configure globally (default for Claude Code)")
	setupCmd.Flags().BoolVarP(&setupProject, "project", "p", false, "Configure for current project only")
	setupCmd.Flags().BoolVar(&setupOpenCode, "opencode", false, "Configure for OpenCode")
	setupCmd.Flags().BoolVar(&setupCodex, "codex", false, "Configure for Codex CLI")
	setupCmd.Flags().BoolVar(&setupGemini, "gemini", false, "Configure for Gemini CLI")
	setupCmd.Flags().BoolVar(&setupAllowAll, "allow-all", false, "Pre-approve all claude_peepee commands")
}

func runSetup(cmd *cobra.Command, args []string) error {
	binaryPath, err := getBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to get binary path: %w", err)
	}

	if setupOpenCode {
		return setupForOpenCode(binaryPath)
	}
	if setupCodex {
		return setupForCodex(binaryPath)
	}
	if setupGemini {
		return setupForGemini(binaryPath)
	}

	// Default: Claude Code
	return setupForClaudeCode(binaryPath)
}

func getBinaryPath() (string, error) {
	// Try to find claude_peepee in PATH
	path, err := exec.LookPath("claude_peepee")
	if err == nil {
		return path, nil
	}

	// Fall back to current executable
	return os.Executable()
}

func setupForClaudeCode(binaryPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	var configPath string
	if setupProject {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		configPath = filepath.Join(cwd, ".claude", "settings.json")
	} else {
		configPath = filepath.Join(homeDir, ".claude", "settings.json")
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read existing config
	config := make(map[string]interface{})
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &config)
	}

	// Add MCP server config
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
	}

	mcpServers["claude_peepee"] = map[string]interface{}{
		"command": binaryPath,
		"args":    []string{"serve"},
	}
	config["mcpServers"] = mcpServers

	// Add permissions if requested
	if setupAllowAll {
		permissions, ok := config["permissions"].(map[string]interface{})
		if !ok {
			permissions = make(map[string]interface{})
		}

		allow, ok := permissions["allow"].([]interface{})
		if !ok {
			allow = []interface{}{}
		}

		// Add claude_peepee tool permissions
		tools := []string{
			"mcp__claude_peepee__remember",
			"mcp__claude_peepee__recall",
			"mcp__claude_peepee__get_context",
			"mcp__claude_peepee__list_instances",
			"mcp__claude_peepee__send_message",
			"mcp__claude_peepee__get_messages",
		}

		for _, tool := range tools {
			found := false
			for _, existing := range allow {
				if existing == tool {
					found = true
					break
				}
			}
			if !found {
				allow = append(allow, tool)
			}
		}

		permissions["allow"] = allow
		config["permissions"] = permissions
	}

	// Write config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	scope := "globally"
	if setupProject {
		scope = "for this project"
	}
	fmt.Printf("Claude PeePee configured %s for Claude Code.\n", scope)
	fmt.Printf("Config file: %s\n", configPath)

	// Create/update CLAUDE.md
	if err := updateClaudeMD(); err != nil {
		fmt.Printf("Warning: failed to update CLAUDE.md: %v\n", err)
	}

	return nil
}

func setupForOpenCode(binaryPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".opencode", "opencode.json")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read existing config
	config := make(map[string]interface{})
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &config)
	}

	// Add MCP server config
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
	}

	mcpServers["claude_peepee"] = map[string]interface{}{
		"command": binaryPath,
		"args":    []string{"serve"},
	}
	config["mcpServers"] = mcpServers

	// Write config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Println("Claude PeePee configured for OpenCode.")
	fmt.Printf("Config file: %s\n", configPath)

	return nil
}

func setupForCodex(binaryPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".codex", "config.toml")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read existing config
	config := make(map[string]interface{})
	if data, err := os.ReadFile(configPath); err == nil {
		toml.Unmarshal(data, &config)
	}

	// Add MCP server config
	mcpServers, ok := config["mcp_servers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
	}

	mcpServers["claude_peepee"] = map[string]interface{}{
		"command": binaryPath,
		"args":    []string{"serve"},
	}
	config["mcp_servers"] = mcpServers

	// Write config
	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Println("Claude PeePee configured for Codex CLI.")
	fmt.Printf("Config file: %s\n", configPath)

	return nil
}

func setupForGemini(binaryPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".gemini", "settings.json")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read existing config
	config := make(map[string]interface{})
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &config)
	}

	// Add MCP server config
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
	}

	mcpServers["claude_peepee"] = map[string]interface{}{
		"command": binaryPath,
		"args":    []string{"serve"},
	}
	config["mcpServers"] = mcpServers

	// Write config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Println("Claude PeePee configured for Gemini CLI.")
	fmt.Printf("Config file: %s\n", configPath)

	return nil
}

func updateClaudeMD() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	claudeMDPath := filepath.Join(cwd, "CLAUDE.md")
	content := `# Claude PeePee - Persistent Memory

This project uses Claude PeePee for persistent memory across Claude Code sessions.

## Available Tools

- **mcp__claude_peepee__remember**: Store facts, decisions, and context
- **mcp__claude_peepee__recall**: Search for previously stored facts
- **mcp__claude_peepee__get_context**: Load all relevant context for this directory
- **mcp__claude_peepee__list_instances**: Find other running Claude Code instances
- **mcp__claude_peepee__send_message**: Send messages to other instances
- **mcp__claude_peepee__get_messages**: Retrieve messages from other instances

## Workflow Recommendations

1. At the start of each session, use ` + "`get_context`" + ` to load previous context
2. Store important decisions and architectural notes using ` + "`remember`" + `
3. When working in a monorepo, use ` + "`list_instances`" + ` and messaging to coordinate
4. Periodically check for messages from other instances
`

	// Check if file exists and contains our content
	if data, err := os.ReadFile(claudeMDPath); err == nil {
		if strings.Contains(string(data), "Claude PeePee") {
			return nil // Already configured
		}
		// Append to existing file
		content = string(data) + "\n\n" + content
	}

	return os.WriteFile(claudeMDPath, []byte(content), 0644)
}
