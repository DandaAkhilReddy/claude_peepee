package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DandaAkhilReddy/claude_pp/internal/telemetry"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claude_pp",
	Short: "Claude PP - Persistent Memory for Claude Code",
	Long: `Claude PP is an MCP server that provides AI coding tools with
persistent memory and multi-instance communication.

Features:
  - Persistent Memory: Store facts, decisions, and context across sessions
  - Multi-Instance Discovery: Find and message other Claude Code instances
  - Automatic Context: Load relevant information based on working directory`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "claude_pp" {
			telemetry.Init()
			telemetry.TrackCommand(cmd.Name())
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		telemetry.Close()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(rememberCmd)
	rootCmd.AddCommand(recallCmd)
	rootCmd.AddCommand(instancesCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(messagesCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(uiCmd)
}

// getDataDir returns the data directory path
func getDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, ".claude_pp")
}
