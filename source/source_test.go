package source

import (
	"os"
	"path/filepath"
	"testing"
)

type MockSource struct {
	data map[string]string
}

func (m *MockSource) Lookup(key string) (string, bool) {
	val, found := m.data[key]
	return val, found
}

func (m *MockSource) All() map[string]string {
	all := make(map[string]string, len(m.data))
	for key, value := range m.data {
		all[key] = value
	}
	return all
}

func TestSystemSource(t *testing.T) {
	key := "TEST_ENV_VAR_KEY"
	expectedValue := "TEST_ENV_VAR_VALUE"
	t.Setenv(key, expectedValue)
	src := &SystemSource{}

	val, found := src.Lookup(key)
	if !found {
		t.Fatalf("expected %s key to be found", key)
	}
	if val != expectedValue {
		t.Fatalf("expected %s, got %s", expectedValue, val)
	}

	emptyKey := "TEST_ENV_VAR_EMPTY"
	t.Setenv(emptyKey, "")
	if val, found := src.Lookup(emptyKey); !found || val != "" {
		t.Fatalf("expected %s to exist with an empty value, got %q (found: %t)", emptyKey, val, found)
	}

	if _, found := src.Lookup("NON_EXISTENCE_KEY"); found {
		t.Error("expected non-existent key to return false")
	}

	if _, found := src.All()[key]; !found {
		t.Errorf("expected All() to contain %s", key)
	}
}

func TestFileSource(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	content := "# comment\n\napp_port=8080\napp_env = production\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	fileSrc, err := NewFileSource(path)
	if err != nil {
		t.Fatalf("failed to initialize NewFileSource: %v", err)
	}

	if val, found := fileSrc.Lookup("app_port"); !found || val != "8080" {
		t.Errorf("expected app_port to be '8080', got '%s' (found: %t)", val, found)
	}
	if val, found := fileSrc.Lookup("app_env"); !found || val != "production" {
		t.Errorf("expected app_env to be 'production', got '%s' (found: %t)", val, found)
	}
	if _, found := fileSrc.Lookup("missing_key"); found {
		t.Error("expected missing_key to return false")
	}
	if _, found := fileSrc.All()["app_port"]; !found {
		t.Error("expected All() to contain app_port")
	}
}

func TestFileSourceMissingFile(t *testing.T) {
	if _, err := NewFileSource(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
		t.Fatal("expected an error when the file does not exist")
	}
}

func TestManagerLookupFallsBack(t *testing.T) {
	first := &MockSource{data: map[string]string{"SHARED": "first"}}
	second := &MockSource{data: map[string]string{"SHARED": "second", "ONLY_SECOND": "value"}}
	mgr := NewManager(first, second)

	if val, found := mgr.Lookup("SHARED"); !found || val != "first" {
		t.Errorf("expected the first source to win, got '%s' (found: %t)", val, found)
	}
	if val, found := mgr.Lookup("ONLY_SECOND"); !found || val != "value" {
		t.Errorf("expected a fallback to the second source, got '%s' (found: %t)", val, found)
	}
	if _, found := mgr.Lookup("MISSING"); found {
		t.Error("expected a missing key to return false")
	}

	all := mgr.All()
	if all["SHARED"] != "first" {
		t.Errorf("expected All() to keep the first source value, got '%s'", all["SHARED"])
	}
	if all["ONLY_SECOND"] != "value" {
		t.Errorf("expected All() to include the fallback source, got '%s'", all["ONLY_SECOND"])
	}
}
