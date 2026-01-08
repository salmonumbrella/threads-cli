package cmd

import "testing"

func TestSearchCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewSearchCmd(f)

	if cmd.Use != "search [query]" {
		t.Errorf("expected Use='search [query]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}
}

func TestSearchCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewSearchCmd(f)

	flags := []string{"limit", "cursor", "media-type", "since", "until", "mode", "type"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}

	limitFlag := cmd.Flag("limit")
	if limitFlag.DefValue != "25" {
		t.Errorf("expected limit default=25, got %s", limitFlag.DefValue)
	}
}
