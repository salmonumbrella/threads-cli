package cmd

import (
	"testing"
)

func TestRateLimitCmd_Structure(t *testing.T) {
	cmd := newRateLimitCmd()

	if cmd.Use != "ratelimit" {
		t.Errorf("expected Use=ratelimit, got %s", cmd.Use)
	}

	// Check aliases
	expectedAliases := []string{"rate", "limits"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	// Check subcommands
	subcommands := cmd.Commands()
	if len(subcommands) != 2 {
		t.Errorf("expected 2 subcommands, got %d", len(subcommands))
	}
}

func TestRateLimitStatusCmd_Structure(t *testing.T) {
	cmd := newRateLimitStatusCmd()

	if cmd.Use != "status" {
		t.Errorf("expected Use=status, got %s", cmd.Use)
	}
}

func TestRateLimitPublishingCmd_Structure(t *testing.T) {
	cmd := newRateLimitPublishingCmd()

	if cmd.Use != "publishing" {
		t.Errorf("expected Use=publishing, got %s", cmd.Use)
	}
}
