package cmd

import (
	"testing"
)

func TestCompletionCmd_Structure(t *testing.T) {
	cmd := NewCompletionCmd()

	if cmd.Use != "completion [bash|zsh|fish|powershell]" {
		t.Errorf("expected Use='completion [bash|zsh|fish|powershell]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}
}

func TestCompletionCmd_ValidArgs(t *testing.T) {
	cmd := NewCompletionCmd()

	expectedShells := []string{"bash", "zsh", "fish", "powershell"}

	if len(cmd.ValidArgs) != len(expectedShells) {
		t.Errorf("expected %d valid args, got %d", len(expectedShells), len(cmd.ValidArgs))
	}

	validArgs := make(map[string]bool)
	for _, arg := range cmd.ValidArgs {
		validArgs[arg] = true
	}

	for _, shell := range expectedShells {
		if !validArgs[shell] {
			t.Errorf("missing valid arg: %s", shell)
		}
	}
}
