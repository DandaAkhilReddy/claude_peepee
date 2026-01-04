package telemetry

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if Version != "0.1.0" {
		t.Errorf("Version = %s, want 0.1.0", Version)
	}
}

func TestGenerateAnonymousID(t *testing.T) {
	id := generateAnonymousID()

	if len(id) != 16 {
		t.Errorf("Anonymous ID length = %d, want 16", len(id))
	}

	// Should be consistent
	id2 := generateAnonymousID()
	if id != id2 {
		t.Error("Anonymous ID should be consistent for same machine")
	}

	// Should be hex characters
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("Anonymous ID should be hex, got char: %c", c)
		}
	}
}

func TestTrackFunctionsNoOp(t *testing.T) {
	// These should not panic even when disabled
	TrackInstall("test")
	TrackSetup("test")
	TrackCommand("test")
	TrackMCPTool("test")
	TrackError("test")
}

func TestClose(t *testing.T) {
	// Should not panic even when client is nil
	Close()
}

func TestInit(t *testing.T) {
	// Just verify Init doesn't panic
	Init()

	// After init, disabled should be true (no backend configured)
	if !disabled {
		t.Error("Telemetry should be disabled by default")
	}
}
