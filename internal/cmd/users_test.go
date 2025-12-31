package cmd

import (
	"testing"
)

func TestUsersMentionsCmd_Structure(t *testing.T) {
	cmd := newUsersMentionsCmd()

	if cmd.Use != "mentions" {
		t.Errorf("expected Use=mentions, got %s", cmd.Use)
	}
}

func TestUsersMentionsCmd_Flags(t *testing.T) {
	cmd := newUsersMentionsCmd()

	flags := []string{"limit", "cursor"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}

	// Check default values
	limitFlag := cmd.Flag("limit")
	if limitFlag.DefValue != "25" {
		t.Errorf("expected limit default=25, got %s", limitFlag.DefValue)
	}
}

func TestUsersCmd_Subcommands(t *testing.T) {
	// usersCmd is a package-level var
	cmd := usersCmd

	subcommands := cmd.Commands()
	if len(subcommands) < 3 {
		t.Errorf("expected at least 3 subcommands (me, get, lookup, mentions), got %d", len(subcommands))
	}

	// Check that mentions is registered
	mentionsFound := false
	for _, sub := range subcommands {
		if sub.Use == "mentions" {
			mentionsFound = true
			break
		}
	}
	if !mentionsFound {
		t.Error("mentions subcommand not found in usersCmd")
	}
}
