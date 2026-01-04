package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/DandaAkhilReddy/claude_pp/internal/store"
	"github.com/spf13/cobra"
)

var rememberTags []string

var rememberCmd = &cobra.Command{
	Use:   "remember [fact]",
	Short: "Store a fact or decision",
	Long:  `Store a fact, decision, or piece of context that should persist across Claude Code sessions.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runRemember,
}

func init() {
	rememberCmd.Flags().StringSliceVarP(&rememberTags, "tags", "t", nil, "Tags to categorize this fact")
}

func runRemember(cmd *cobra.Command, args []string) error {
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	s, err := store.NewSQLiteStore(dataDir)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer s.Close()

	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	fact := strings.Join(args, " ")

	id, err := s.AddFact(fact, rememberTags, workingDir)
	if err != nil {
		return fmt.Errorf("failed to store fact: %w", err)
	}

	fmt.Printf("Stored fact #%d\n", id)
	return nil
}
