package telemetry

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"sync"
)

const (
	Version = "0.1.0"
)

var (
	anonymousID string
	initOnce    sync.Once
	disabled    bool
)

// Init initializes telemetry (respecting user privacy preferences)
func Init() {
	initOnce.Do(func() {
		// Check for opt-out
		if os.Getenv("CLAUDE_PP_NO_TELEMETRY") == "1" || os.Getenv("DO_NOT_TRACK") == "1" {
			disabled = true
			return
		}

		// Generate anonymous ID from machine info
		anonymousID = generateAnonymousID()

		// Telemetry is disabled by default - can be enabled later with proper backend
		disabled = true
	})
}

// generateAnonymousID creates a stable anonymous ID for this machine
func generateAnonymousID() string {
	homeDir, _ := os.UserHomeDir()
	hostname, _ := os.Hostname()
	data := homeDir + ":" + hostname
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}

// TrackInstall tracks installation events
func TrackInstall(method string) {
	if disabled {
		return
	}
	// Placeholder for future telemetry implementation
}

// TrackSetup tracks setup events
func TrackSetup(target string) {
	if disabled {
		return
	}
	// Placeholder for future telemetry implementation
}

// TrackCommand tracks CLI command usage
func TrackCommand(command string) {
	if disabled {
		return
	}
	// Placeholder for future telemetry implementation
}

// TrackMCPTool tracks MCP tool usage
func TrackMCPTool(tool string) {
	if disabled {
		return
	}
	// Placeholder for future telemetry implementation
}

// TrackError tracks errors (anonymized)
func TrackError(context string) {
	if disabled {
		return
	}
	// Placeholder for future telemetry implementation
}

// Close flushes and closes the telemetry client
func Close() {
	// Placeholder for future telemetry implementation
}
