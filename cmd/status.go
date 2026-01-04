package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/DandaAkhilReddy/claude_peepee/internal/store"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status information",
	Long:  `Show information about stored facts and running instances.`,
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	dataDir := getDataDir()

	s, err := store.NewSQLiteStore(dataDir)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer s.Close()

	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Get fact counts
	total, local, err := s.CountFacts(workingDir)
	if err != nil {
		return fmt.Errorf("failed to count facts: %w", err)
	}

	// Cleanup and get instances
	if err := s.CleanupStaleInstances(5 * time.Minute); err != nil {
		return fmt.Errorf("failed to cleanup stale instances: %w", err)
	}

	instances, err := s.GetInstances()
	if err != nil {
		return fmt.Errorf("failed to get instances: %w", err)
	}

	fmt.Printf("Claude PP Status\n")
	fmt.Printf("================\n\n")
	fmt.Printf("Data directory: %s\n", dataDir)
	fmt.Printf("Working directory: %s\n\n", workingDir)
	fmt.Printf("Facts:\n")
	fmt.Printf("  Total: %d\n", total)
	fmt.Printf("  Local (this directory): %d\n\n", local)
	fmt.Printf("Running instances: %d\n", len(instances))

	for _, inst := range instances {
		fmt.Printf("  - %s: %s\n", inst.ID, inst.WorkingDir)
	}

	return nil
}
