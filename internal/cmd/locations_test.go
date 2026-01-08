package cmd

import "testing"

func TestLocationsCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewLocationsCmd(f)

	if cmd.Use != "locations" {
		t.Errorf("expected Use=locations, got %s", cmd.Use)
	}

	expectedAliases := []string{"location", "loc"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	subcommands := cmd.Commands()
	if len(subcommands) != 2 {
		t.Errorf("expected 2 subcommands, got %d", len(subcommands))
	}

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
	f := newTestFactory(t)
	cmd := newLocationsSearchCmd(f)

	flags := []string{"lat", "lng"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}
}

func TestLocationsGetCmd_RequiresArg(t *testing.T) {
	f := newTestFactory(t)
	cmd := newLocationsGetCmd(f)

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}
}
