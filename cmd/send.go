package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/DandaAkhilReddy/claude_peepee/internal/store"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send [instance-id] [message]",
	Short: "Send a message to another instance",
	Long:  `Send a message to another running claude_pp instance.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  runSend,
}

func runSend(cmd *cobra.Command, args []string) error {
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	s, err := store.NewSQLiteStore(dataDir)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer s.Close()

	toInstance := args[0]
	message := strings.Join(args[1:], " ")

	// Verify target instance exists
	inst, err := s.GetInstance(toInstance)
	if err != nil {
		return fmt.Errorf("failed to verify instance: %w", err)
	}
	if inst == nil {
		return fmt.Errorf("instance not found: %s", toInstance)
	}

	id, err := s.SendMessage("cli", toInstance, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("Message #%d sent to %s (%s)\n", id, toInstance, inst.WorkingDir)
	return nil
}
