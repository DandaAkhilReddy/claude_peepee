package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DandaAkhilReddy/claude_peepee/internal/mcp"
	"github.com/DandaAkhilReddy/claude_peepee/internal/store"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long:  `Start the Claude PP MCP server for Claude Code integration.`,
	RunE:  runServe,
}

func runServe(cmd *cobra.Command, args []string) error {
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

	instanceID := uuid.New().String()[:8]
	pid := os.Getpid()

	if err := s.RegisterInstance(instanceID, pid, workingDir); err != nil {
		return fmt.Errorf("failed to register instance: %w", err)
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		s.UnregisterInstance(instanceID)
		cancel()
		os.Exit(0)
	}()

	// Start heartbeat
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.UpdateHeartbeat(instanceID)
			}
		}
	}()

	server := mcp.NewServer(s, instanceID, workingDir)
	return server.Run()
}
