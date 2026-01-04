package cmd

import (
	"fmt"

	"github.com/DandaAkhilReddy/claude_pp/internal/store"
	"github.com/spf13/cobra"
)

var messagesAll bool

var messagesCmd = &cobra.Command{
	Use:   "messages [instance-id]",
	Short: "View messages for an instance",
	Long:  `View messages sent to a specific instance.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runMessages,
}

func init() {
	messagesCmd.Flags().BoolVarP(&messagesAll, "all", "a", false, "Show all messages (including read)")
}

func runMessages(cmd *cobra.Command, args []string) error {
	dataDir := getDataDir()

	s, err := store.NewSQLiteStore(dataDir)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer s.Close()

	instanceID := args[0]
	unreadOnly := !messagesAll

	messages, err := s.GetMessages(instanceID, unreadOnly)
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}

	if len(messages) == 0 {
		if unreadOnly {
			fmt.Println("No unread messages.")
		} else {
			fmt.Println("No messages found.")
		}
		return nil
	}

	for _, msg := range messages {
		status := "unread"
		if msg.ReadAt != nil {
			status = fmt.Sprintf("read at %s", msg.ReadAt.Format("15:04"))
		}

		fmt.Printf("\n#%d [%s] from %s (%s)\n", msg.ID, msg.CreatedAt.Format("2006-01-02 15:04"), msg.FromInstance, status)
		fmt.Printf("%s\n", msg.Content)
	}

	return nil
}
