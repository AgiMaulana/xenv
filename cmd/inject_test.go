package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInjectCommand_ResolvesAndPassesSecret(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("TEST_TOKEN=super-secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin,
		"-i", envFile,
		"inject", "TEST_TOKEN", "--",
		"printenv", "TEST_TOKEN",
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("inject failed: %v\nstderr: %s", err, string(out))
	}

	got := strings.TrimSpace(string(out))
	if got != "super-secret" {
		t.Fatalf("expected 'super-secret', got %q", got)
	}
}

func TestInjectCommand_MissingKeyFails(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("EXISTING=ok\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin,
		"-i", envFile,
		"inject", "EXISTING",
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing command, got none")
	}
	_ = out
}

func TestInjectCommand_MissingKeyInSources(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("EXISTING=ok\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin,
		"-i", envFile,
		"inject", "EXISTING", "MISSING_KEY", "--",
		"printenv", "EXISTING",
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing key, got none")
	}
	if !strings.Contains(string(out), "not found") {
		t.Fatalf("error should mention missing key, got: %s", out)
	}
}

func TestInjectCommand_SecretsNotInParentsEnv(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("TEST_TOKEN=injected-only\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin,
		"-i", envFile,
		"inject", "TEST_TOKEN", "--",
		"true",
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("inject should succeed: %v", err)
	}

	if val, set := os.LookupEnv("TEST_TOKEN"); set {
		t.Fatalf("TEST_TOKEN leaked to test process after inject: %q", val)
	}
}

func TestInjectCommand_WithoutDashDashSeparator(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("MY_VAR=hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin,
		"-i", envFile,
		"inject", "MY_VAR",
		"printenv", "MY_VAR",
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("inject failed: %v\nstderr: %s", err, string(out))
	}

	got := strings.TrimSpace(string(out))
	if got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestInjectCommand_MultipleKeys(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("FIRST=one\nSECOND=two\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin,
		"-i", envFile,
		"inject", "FIRST", "SECOND", "--",
		"sh", "-c", `printf '%s %s' "$FIRST" "$SECOND"`,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("inject failed: %v\nstderr: %s", err, string(out))
	}

	got := strings.TrimSpace(string(out))
	if got != "one two" {
		t.Fatalf("expected 'one two', got %q", got)
	}
}

func TestMain(m *testing.M) {
	// Build the test binary once before all tests.
	dir, err := os.MkdirTemp("", "xenv-test-bin-*")
	if err != nil {
		os.Stderr.WriteString("failed to create temp dir: " + err.Error() + "\n")
		os.Exit(1)
	}
	bin = filepath.Join(dir, "xenv-test")
	build := exec.Command("go", "build", "-o", bin, "xenv")
	if out, err := build.CombinedOutput(); err != nil {
		os.Stderr.WriteString("failed to build test binary: " + err.Error() + "\n")
		os.Stderr.Write(out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

var bin string
