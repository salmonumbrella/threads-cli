package cmd

import "testing"

func TestWebhooksCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewWebhooksCmd(f)

	if cmd.Use != "webhooks" {
		t.Errorf("expected Use=webhooks, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}
}

func TestWebhooksCmd_Subcommands(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewWebhooksCmd(f)

	expectedSubs := map[string]bool{
		"subscribe": true,
		"list":      true,
		"delete":    true,
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
