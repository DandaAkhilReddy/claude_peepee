package cmd

import (
	"fmt"
	"time"

	"github.com/DandaAkhilReddy/claude_peepee/internal/store"
	"github.com/spf13/cobra"
)

var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List running instances",
	Long:  `List all running claude_peepee instances across different directories.`,
	RunE:  runInstances,
}

func runInstances(cmd *cobra.Command, args []string) error {
	dataDir := getDataDir()

	s, err := store.NewSQLiteStore(dataDir)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer s.Close()

	// Clean up stale instances (older than 5 minutes)
	if err := s.CleanupStaleInstances(5 * time.Minute); err != nil {
		return fmt.Errorf("failed to cleanup stale instances: %w", err)
	}

	instances, err := s.GetInstances()
	if err != nil {
		return fmt.Errorf("failed to get instances: %w", err)
	}

	if len(instances) == 0 {
		fmt.Println("No running instances found.")
		return nil
	}

	fmt.Printf("Running instances (%d):\n\n", len(instances))
	for _, inst := range instances {
		fmt.Printf("ID: %s\n", inst.ID)
		fmt.Printf("  PID: %d\n", inst.PID)
		fmt.Printf("  Dir: %s\n", inst.WorkingDir)
		fmt.Printf("  Started: %s\n", inst.StartedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Last heartbeat: %s\n\n", inst.LastHeartbeat.Format("2006-01-02 15:04:05"))
	}

	return nil
}
