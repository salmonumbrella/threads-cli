package cmd

import (
	"testing"
)

func TestLocationsCmd_Structure(t *testing.T) {
	cmd := newLocationsCmd()

	if cmd.Use != "locations" {
		t.Errorf("expected Use=locations, got %s", cmd.Use)
	}

	// Check aliases
	expectedAliases := []string{"location", "loc"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	// Check subcommands
	subcommands := cmd.Commands()
	if len(subcommands) != 2 {
		t.Errorf("expected 2 subcommands, got %d", len(subcommands))
	}

	// Verify subcommand names
	names := make(map[string]bool)
	for _, sub := range subcommands {
		names[sub.Use] = true
	}
	if !names["search [query]"] {
		t.Error("missing 'search' subcommand")
	}
	if !names["get [location-id]"] {
		t.Error("missing 'get' subcommand")
	}
}

func TestLocationsSearchCmd_Flags(t *testing.T) {
	cmd := newLocationsSearchCmd()

	flags := []string{"lat", "lng"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}
}

func TestLocationsGetCmd_RequiresArg(t *testing.T) {
	cmd := newLocationsGetCmd()

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}
}
