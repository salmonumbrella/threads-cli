package cmd

import "testing"

func TestAuthCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewAuthCmd(f)

	if cmd.Use != "auth" {
		t.Errorf("expected Use=auth, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	expectedSubs := map[string]bool{
		"login":   true,
		"token":   true,
		"refresh": true,
		"status":  true,
		"list":    true,
		"remove":  true,
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

func TestAuthLoginCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := newAuthLoginCmd(f)

	flags := []struct {
		name      string
		shorthand string
	}{
		{"name", "n"},
		{"client-id", ""},
		{"client-secret", ""},
		{"redirect-uri", ""},
		{"scopes", ""},
	}

	for _, flag := range flags {
		fl := cmd.Flag(flag.name)
		if fl == nil {
			t.Errorf("missing flag: %s", flag.name)
			continue
		}
		if flag.shorthand != "" && fl.Shorthand != flag.shorthand {
			t.Errorf("flag %s expected shorthand %q, got %q", flag.name, flag.shorthand, fl.Shorthand)
		}
	}

	nameFlag := cmd.Flag("name")
	if nameFlag.DefValue != "default" {
		t.Errorf("expected name default='default', got %s", nameFlag.DefValue)
	}
}

func TestAuthTokenCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newAuthTokenCmd(f)

	if cmd.Use != "token [access-token]" {
		t.Errorf("expected Use='token [access-token]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestAuthTokenCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := newAuthTokenCmd(f)

	flags := []string{"name", "client-id", "client-secret"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}
}

func TestAuthRefreshCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newAuthRefreshCmd(f)

	if cmd.Use != "refresh" {
		t.Errorf("expected Use=refresh, got %s", cmd.Use)
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestAuthStatusCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newAuthStatusCmd(f)

	if cmd.Use != "status" {
		t.Errorf("expected Use=status, got %s", cmd.Use)
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestAuthListCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newAuthListCmd(f)

	if cmd.Use != "list" {
		t.Errorf("expected Use=list, got %s", cmd.Use)
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestAuthRemoveCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newAuthRemoveCmd(f)

	if cmd.Use != "remove [account]" {
		t.Errorf("expected Use='remove [account]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}
