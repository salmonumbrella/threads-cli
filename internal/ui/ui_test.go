package ui

import (
	"bytes"
	"testing"
	"time"

	"github.com/salmonumbrella/threads-cli/internal/iocontext"
	"github.com/salmonumbrella/threads-cli/internal/outfmt"
)

func TestPrinterSuccess(t *testing.T) {
	var buf bytes.Buffer
	io := &iocontext.IO{Out: &buf, ErrOut: &buf}
	p := New(io, outfmt.ColorNever)

	p.Success("test message %s", "arg")
	if buf.String() == "" {
		t.Error("expected output from Success")
	}
}

func TestPrinterWarning(t *testing.T) {
	var buf bytes.Buffer
	io := &iocontext.IO{Out: &buf, ErrOut: &buf}
	p := New(io, outfmt.ColorNever)

	p.Warning("test warning %d", 42)
	if buf.String() == "" {
		t.Error("expected output from Warning")
	}
}

func TestPrinterInfo(t *testing.T) {
	var buf bytes.Buffer
	io := &iocontext.IO{Out: &buf, ErrOut: &buf}
	p := New(io, outfmt.ColorNever)

	p.Info("test info")
	if buf.String() == "" {
		t.Error("expected output from Info")
	}
}

func TestPrinterError(t *testing.T) {
	var errBuf bytes.Buffer
	io := &iocontext.IO{Out: &bytes.Buffer{}, ErrOut: &errBuf}
	p := New(io, outfmt.ColorNever)

	p.Error("test error %v", "details")
	if errBuf.String() == "" {
		t.Error("expected output from Error")
	}
}

func TestPrinterBold(t *testing.T) {
	var buf bytes.Buffer
	io := &iocontext.IO{Out: &buf, ErrOut: &buf}
	p := New(io, outfmt.ColorNever)

	result := p.Bold("test")
	if result == "" {
		t.Error("Bold should return non-empty string")
	}
}

func TestPrinterDim(t *testing.T) {
	var buf bytes.Buffer
	io := &iocontext.IO{Out: &buf, ErrOut: &buf}
	p := New(io, outfmt.ColorNever)

	result := p.Dim("test")
	if result == "" {
		t.Error("Dim should return non-empty string")
	}
}

func TestPrinterColorize(t *testing.T) {
	var buf bytes.Buffer
	io := &iocontext.IO{Out: &buf, ErrOut: &buf}
	p := New(io, outfmt.ColorNever)

	result := p.Colorize("test", p.Green)
	if result == "" {
		t.Error("Colorize should return non-empty string")
	}
}

func TestPrinterStatusColor(t *testing.T) {
	var buf bytes.Buffer
	io := &iocontext.IO{Out: &buf, ErrOut: &buf}
	p := New(io, outfmt.ColorNever)

	tests := []struct {
		status string
	}{
		{"active"},
		{"valid"},
		{"published"},
		{"success"},
		{"expired"},
		{"error"},
		{"failed"},
		{"pending"},
		{"processing"},
		{"unknown"},
		{""},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			color := p.StatusColor(tt.status)
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
	_ = IsTerminal()
}

func TestFormatRelativeTimeWithNow(t *testing.T) {
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
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		offset   time.Duration
		expected string
	}{
		{"30 seconds from now", 30 * time.Second, "now"},
		{"5 minutes from now", 5 * time.Minute, "in 5m"},
		{"45 minutes from now", 45 * time.Minute, "in 45m"},
		{"1 hour from now", time.Hour, "in 1h"},
		{"1.5 hours from now", 90 * time.Minute, "in 1h 30m"},
		{"3 hours from now", 3 * time.Hour, "in 3h"},
		{"1 day from now", 24 * time.Hour, "in 1d"},
		{"2 days from now", 48 * time.Hour, "in 2d"},
		{"1 week from now", 7 * 24 * time.Hour, "in 1w"},
		{"2 weeks from now", 14 * 24 * time.Hour, "in 2w"},
		{"30 seconds ago", -30 * time.Second, "now"},
		{"5 minutes ago", -5 * time.Minute, "5m ago"},
		{"45 minutes ago", -45 * time.Minute, "45m ago"},
		{"1 hour ago", -1 * time.Hour, "1h ago"},
		{"1.5 hours ago", -90 * time.Minute, "1h 30m ago"},
		{"3 hours ago", -3 * time.Hour, "3h ago"},
		{"1 day ago", -24 * time.Hour, "1d ago"},
		{"2 days ago", -48 * time.Hour, "2d ago"},
		{"1 week ago", -7 * 24 * time.Hour, "1w ago"},
		{"2 weeks ago", -14 * 24 * time.Hour, "2w ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRelativeTimeFrom(now.Add(tt.offset), now)
			if result != tt.expected {
				t.Errorf("formatRelativeTimeFrom(%v) = %q, want %q", tt.offset, result, tt.expected)
			}
		})
	}
}
