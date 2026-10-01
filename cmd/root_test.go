package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSourceManagerPrefersInputFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("FROM_FILE=yes\n"), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	restoreInput := setInputFile(t, path)
	defer restoreInput()
	t.Setenv("FROM_FILE", "from-env")

	manager, err := newSourceManager()
	if err != nil {
		t.Fatalf("newSourceManager returned error: %v", err)
	}

	if val, found := manager.Lookup("FROM_FILE"); !found || val != "yes" {
		t.Errorf("expected the input file to win, got %q (found: %t)", val, found)
	}
}

func TestNewSourceManagerUsesSystemEnvironment(t *testing.T) {
	restoreInput := setInputFile(t, "")
	defer restoreInput()
	t.Setenv("XENV_SYSTEM_ONLY", "present")

	manager, err := newSourceManager()
	if err != nil {
		t.Fatalf("newSourceManager returned error: %v", err)
	}

	if val, found := manager.Lookup("XENV_SYSTEM_ONLY"); !found || val != "present" {
		t.Errorf("expected the system environment, got %q (found: %t)", val, found)
	}
}

func TestNewSourceManagerMissingFileFallsBackToSystem(t *testing.T) {
	restoreInput := setInputFile(t, filepath.Join(t.TempDir(), "does-not-exist"))
	defer restoreInput()
	t.Setenv("XENV_SYSTEM_ONLY", "present")

	manager, err := newSourceManager()
	if err != nil {
		t.Fatalf("newSourceManager returned error: %v", err)
	}

	if _, found := manager.Lookup("XENV_SYSTEM_ONLY"); !found {
		t.Error("expected a missing input file to fall back to the system environment")
	}
}

func setInputFile(t *testing.T, path string) func() {
	t.Helper()
	previous := inputFile
	inputFile = path
	return func() { inputFile = previous }
}
