package ui

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

// captureOutput captures stdout during a function call
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestSuccess(t *testing.T) {
	output := captureOutput(t, func() {
		Success("test message %s", "arg")
	})
	if output == "" {
		t.Error("expected output from Success")
	}
	// The output contains ANSI codes, just check it's not empty
}

func TestWarning(t *testing.T) {
	output := captureOutput(t, func() {
		Warning("test warning %d", 42)
	})
	if output == "" {
		t.Error("expected output from Warning")
	}
}

func TestInfo(t *testing.T) {
	output := captureOutput(t, func() {
		Info("test info")
	})
	if output == "" {
		t.Error("expected output from Info")
	}
}

func TestError(t *testing.T) {
	// Error writes to stderr, need different capture
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w

	Error("test error %v", "details")

	_ = w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if output == "" {
		t.Error("expected output from Error")
	}
}

func TestBold(t *testing.T) {
	result := Bold("test")
	if result == "" {
		t.Error("Bold should return non-empty string")
	}
	// Result contains ANSI codes or the plain text depending on terminal
}

func TestDim(t *testing.T) {
	result := Dim("test")
	if result == "" {
		t.Error("Dim should return non-empty string")
	}
}

func TestColorize(t *testing.T) {
	result := Colorize("test", Green)
	if result == "" {
		t.Error("Colorize should return non-empty string")
	}

	result2 := Colorize("test", Red)
	if result2 == "" {
		t.Error("Colorize should return non-empty string")
	}
}

func TestStatusColor(t *testing.T) {
	tests := []struct {
		status   string
		expected interface{}
	}{
		{"active", Green},
		{"valid", Green},
		{"published", Green},
		{"success", Green},
		{"expired", Red},
		{"error", Red},
		{"failed", Red},
		{"pending", Yellow},
		{"processing", Yellow},
		{"unknown", Gray},
		{"", Gray},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			color := StatusColor(tt.status)
			// Just verify it returns a color without error
			if color == nil {
				t.Errorf("StatusColor(%q) returned nil", tt.status)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		days     float64
		expected string
	}{
		{-1, "unknown"},
		{0.5, "12 hours"},
		{0.25, "6 hours"},
		{1, "1 days"},
		{3, "3 days"},
		{6.5, "6 days"},
		{7, "1.0 weeks"},
		{14, "2.0 weeks"},
		{10.5, "1.5 weeks"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatDuration(tt.days)
			if result != tt.expected {
				t.Errorf("FormatDuration(%v) = %q, want %q", tt.days, result, tt.expected)
			}
		})
	}
}

func TestIsTerminal(t *testing.T) {
	// Just verify it returns without error
	_ = IsTerminal()
}

func TestFormatRelativeTimeWithNow(t *testing.T) {
	// Test the public FormatRelativeTime function
	future := time.Now().Add(time.Hour)
	result := FormatRelativeTime(future)
	if result == "" {
		t.Error("FormatRelativeTime should return non-empty string")
	}

	past := time.Now().Add(-time.Hour)
	result2 := FormatRelativeTime(past)
	if result2 == "" {
		t.Error("FormatRelativeTime should return non-empty string")
	}
}

func TestFormatRelativeTime(t *testing.T) {
	// Use a fixed reference time for deterministic tests
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		offset   time.Duration
		expected string
	}{
		// Future times
		{"30 seconds from now", 30 * time.Second, "now"},
		{"5 minutes from now", 5 * time.Minute, "in 5m"},
		{"45 minutes from now", 45 * time.Minute, "in 45m"},
		{"1 hour from now", time.Hour, "in 1h"},
		{"1.5 hours from now", 90 * time.Minute, "in 1h 30m"},
		{"3 hours from now", 3 * time.Hour, "in 3h"},
		{"1 day from now", 24 * time.Hour, "in 1d"},
		{"1 day 12 hours from now", 36 * time.Hour, "in 1d 12h"},
		{"3 days from now", 72 * time.Hour, "in 3d"},
		{"1 week from now", 7 * 24 * time.Hour, "in 1w"},
		{"10 days from now", 10 * 24 * time.Hour, "in 1w 3d"},

		// Past times
		{"-30 seconds ago", -30 * time.Second, "now"},
		{"5 minutes ago", -5 * time.Minute, "5m ago"},
		{"45 minutes ago", -45 * time.Minute, "45m ago"},
		{"1 hour ago", -time.Hour, "1h ago"},
		{"1.5 hours ago", -90 * time.Minute, "1h 30m ago"},
		{"3 hours ago", -3 * time.Hour, "3h ago"},
		{"1 day ago", -24 * time.Hour, "1d ago"},
		{"1 day 12 hours ago", -36 * time.Hour, "1d 12h ago"},
		{"3 days ago", -72 * time.Hour, "3d ago"},
		{"1 week ago", -7 * 24 * time.Hour, "1w ago"},
		{"10 days ago", -10 * 24 * time.Hour, "1w 3d ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetTime := now.Add(tt.offset)
			result := formatRelativeTimeFrom(targetTime, now)
			if result != tt.expected {
				t.Errorf("formatRelativeTimeFrom(%v) = %q, want %q", tt.offset, result, tt.expected)
			}
		})
	}
}
