package cmd

import "testing"

func TestRepliesCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewRepliesCmd(f)

	if cmd.Use != "replies" {
		t.Errorf("expected Use=replies, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	subcommands := cmd.Commands()
	expectedCount := 5
	if len(subcommands) != expectedCount {
		t.Errorf("expected %d subcommands, got %d", expectedCount, len(subcommands))
	}
}

func TestRepliesCmd_Subcommands(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewRepliesCmd(f)

	expectedSubs := map[string]bool{
		"list":         true,
		"create":       true,
		"hide":         true,
		"unhide":       true,
		"conversation": true,
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

func TestRepliesListCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newRepliesListCmd(f)

	if cmd.Use != "list [post-id]" {
		t.Errorf("expected Use='list [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestRepliesListCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := newRepliesListCmd(f)

	limitFlag := cmd.Flag("limit")
	if limitFlag == nil {
		t.Fatal("missing limit flag")
	}

	if limitFlag.DefValue != "25" {
		t.Errorf("expected limit default=25, got %s", limitFlag.DefValue)
	}
}

func TestRepliesCreateCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newRepliesCreateCmd(f)

	if cmd.Use != "create [post-id]" {
		t.Errorf("expected Use='create [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestRepliesCreateCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := newRepliesCreateCmd(f)

	textFlag := cmd.Flag("text")
	if textFlag == nil {
		t.Fatal("missing text flag")
	}

	if textFlag.Shorthand != "t" {
		t.Errorf("expected text flag shorthand='t', got %s", textFlag.Shorthand)
	}
}

func TestRepliesHideCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newRepliesHideCmd(f)

	if cmd.Use != "hide [reply-id]" {
		t.Errorf("expected Use='hide [reply-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestRepliesUnhideCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newRepliesUnhideCmd(f)

	if cmd.Use != "unhide [reply-id]" {
		t.Errorf("expected Use='unhide [reply-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestRepliesConversationCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newRepliesConversationCmd(f)

	if cmd.Use != "conversation [post-id]" {
		t.Errorf("expected Use='conversation [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestRepliesConversationCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := newRepliesConversationCmd(f)

	limitFlag := cmd.Flag("limit")
	if limitFlag == nil {
		t.Fatal("missing limit flag")
	}

	if limitFlag.DefValue != "25" {
		t.Errorf("expected limit default=25, got %s", limitFlag.DefValue)
	}
}
