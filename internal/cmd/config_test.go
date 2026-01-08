package cmd

import "testing"

func TestConfigCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewConfigCmd(f)

	if cmd.Use != "config" {
		t.Errorf("expected Use=config, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}
}

func TestConfigCmd_Subcommands(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewConfigCmd(f)

	expectedSubs := map[string]bool{
		"path":  true,
		"list":  true,
		"get":   true,
		"set":   true,
		"unset": true,
	}

	for _, sub := range cmd.Commands() {
		name := sub.Name()
		if !expectedSubs[name] {
			t.Errorf("unexpected subcommand: %s", name)
		}
		delete(expectedSubs, name)
	}

	for name := range expectedSubs {
		t.Errorf("missing subcommand: %s", name)
	}
}
