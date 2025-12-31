package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestConfigDir(t *testing.T) {
	dir := ConfigDir()
	if dir == "" {
		t.Error("ConfigDir should not return empty string")
	}

	// Should contain app name
	if !strings.Contains(dir, appName) {
		t.Errorf("ConfigDir should contain app name %q, got %q", appName, dir)
	}

	// Platform-specific checks
	if runtime.GOOS == "darwin" {
		if !strings.Contains(dir, "Library/Application Support") {
			t.Errorf("on macOS, ConfigDir should use Library/Application Support, got %q", dir)
		}
	}
}

func TestConfigDir_XDGOverride(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("XDG override not used on macOS")
	}

	// Set custom XDG_CONFIG_HOME - t.Setenv restores automatically
	testDir := "/tmp/xdg-test-config"
	t.Setenv("XDG_CONFIG_HOME", testDir)

	dir := ConfigDir()
	expected := filepath.Join(testDir, appName)
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

func TestConfigDir_DefaultFallback(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("fallback not used on macOS")
	}

	// Unset XDG_CONFIG_HOME to test fallback - set to empty string
	t.Setenv("XDG_CONFIG_HOME", "")

	dir := ConfigDir()
	home := os.Getenv("HOME")
	expected := filepath.Join(home, ".config", appName)
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

func TestDataDir(t *testing.T) {
	dir := DataDir()
	if dir == "" {
		t.Error("DataDir should not return empty string")
	}

	// Should contain app name
	if !strings.Contains(dir, appName) {
		t.Errorf("DataDir should contain app name %q, got %q", appName, dir)
	}

	// Platform-specific checks
	if runtime.GOOS == "darwin" {
		if !strings.Contains(dir, "Library/Application Support") {
			t.Errorf("on macOS, DataDir should use Library/Application Support, got %q", dir)
		}
	}
}

func TestDataDir_XDGOverride(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("XDG override not used on macOS")
	}

	// Set custom XDG_DATA_HOME
	testDir := "/tmp/xdg-test-data"
	t.Setenv("XDG_DATA_HOME", testDir)

	dir := DataDir()
	expected := filepath.Join(testDir, appName)
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

func TestDataDir_DefaultFallback(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("fallback not used on macOS")
	}

	// Unset XDG_DATA_HOME to test fallback
	t.Setenv("XDG_DATA_HOME", "")

	dir := DataDir()
	home := os.Getenv("HOME")
	expected := filepath.Join(home, ".local", "share", appName)
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

func TestCacheDir(t *testing.T) {
	dir := CacheDir()
	if dir == "" {
		t.Error("CacheDir should not return empty string")
	}

	// Should contain app name
	if !strings.Contains(dir, appName) {
		t.Errorf("CacheDir should contain app name %q, got %q", appName, dir)
	}

	// Platform-specific checks
	if runtime.GOOS == "darwin" {
		if !strings.Contains(dir, "Library/Caches") {
			t.Errorf("on macOS, CacheDir should use Library/Caches, got %q", dir)
		}
	}
}

func TestCacheDir_XDGOverride(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("XDG override not used on macOS")
	}

	// Set custom XDG_CACHE_HOME
	testDir := "/tmp/xdg-test-cache"
	t.Setenv("XDG_CACHE_HOME", testDir)

	dir := CacheDir()
	expected := filepath.Join(testDir, appName)
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

func TestCacheDir_DefaultFallback(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("fallback not used on macOS")
	}

	// Unset XDG_CACHE_HOME to test fallback
	t.Setenv("XDG_CACHE_HOME", "")

	dir := CacheDir()
	home := os.Getenv("HOME")
	expected := filepath.Join(home, ".cache", appName)
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	// Create a temp directory for testing
	tmpDir := t.TempDir()

	if runtime.GOOS == "darwin" {
		// On macOS, set HOME to temp directory
		t.Setenv("HOME", tmpDir)
		expectedDir := filepath.Join(tmpDir, "Library", "Application Support", appName)

		err := EnsureConfigDir()
		if err != nil {
			t.Fatalf("EnsureConfigDir failed: %v", err)
		}

		// Check directory was created
		info, err := os.Stat(expectedDir)
		if err != nil {
			t.Fatalf("expected directory to exist at %q: %v", expectedDir, err)
		}
		if !info.IsDir() {
			t.Errorf("expected %q to be a directory", expectedDir)
		}
	} else {
		// On Linux, use XDG_CONFIG_HOME
		t.Setenv("XDG_CONFIG_HOME", tmpDir)
		expectedDir := filepath.Join(tmpDir, appName)

		err := EnsureConfigDir()
		if err != nil {
			t.Fatalf("EnsureConfigDir failed: %v", err)
		}

		// Check directory was created
		info, err := os.Stat(expectedDir)
		if err != nil {
			t.Fatalf("expected directory to exist at %q: %v", expectedDir, err)
		}
		if !info.IsDir() {
			t.Errorf("expected %q to be a directory", expectedDir)
		}
	}
}

func TestAppNameConstant(t *testing.T) {
	if appName != "threads-cli" {
		t.Errorf("expected appName 'threads-cli', got %q", appName)
	}
}
