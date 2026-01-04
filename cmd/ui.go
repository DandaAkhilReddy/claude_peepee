package cmd

import (
	"fmt"
	"os"

	"github.com/DandaAkhilReddy/claude_peepee/internal/store"
	"github.com/DandaAkhilReddy/claude_peepee/internal/ui"
	"github.com/spf13/cobra"
)

var uiPort int

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the web UI for managing facts",
	Long: `Start a local web server that provides a user-friendly interface
for viewing and managing your Claude PP knowledge base.

The UI allows you to:
  - Browse and search all stored facts
  - View running Claude Code instances
  - Add new facts with tags
  - Delete old or irrelevant facts

Example:
  claude_pp ui           # Start on default port 8420
  claude_pp ui -p 3000   # Start on port 3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dataDir := getDataDir()

		// Ensure data directory exists
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}

		s, err := store.NewSQLiteStore(dataDir)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer s.Close()

		server := ui.NewServer(s, uiPort)
		return server.Start()
	},
}

func init() {
	uiCmd.Flags().IntVarP(&uiPort, "port", "p", 8420, "Port to run the web UI on")
}
