package gohosts

import (
	"os"
	"testing"
)

func TestNewHosts_Default(t *testing.T) {
	h, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h.path == "" {
		t.Error("expected a valid path, but got an empty string")
	}
}

func TestNewHosts_WithPath(t *testing.T) {
	tempPath, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempPath.Name())
	defer tempPath.Close()

	h, err := New(WithPath(tempPath.Name()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h.path != tempPath.Name() {
		t.Errorf("expected path %s, but got %s", tempPath.Name(), h.path)
	}

}

func TestNewHosts_WithInvalidPath(t *testing.T) {
	h, err := New(WithPath("non-existent/path"))
	if err == nil {
		t.Error("expected an error, but got nil")
	}

	if h != nil {
		t.Errorf("expected nil, but got %v", h)
	}
}

func TestLoad(t *testing.T) {
	h, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = h.Load()
	if err != nil {
		t.Errorf("failed to load hosts file: %v", err)
	}
}

func TestLoad_Entries(t *testing.T) {
	tempPath, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempPath.Name())
	defer tempPath.Close()

	h, err := New(WithPath(tempPath.Name()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = h.Load()
	if err != nil {
		t.Errorf("failed to load hosts file: %v", err)
	}

	if len(h.Entries) != 0 {
		t.Errorf("expected 0 entries, but got %d", len(h.Entries))
	}
}

func TestLoad_CorruptedFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "courpted")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := make([]byte, 65536+1) // 64KB + 1 byte to exceed default buffer size.
	for i := range content {
		content[i] = 'a'
	}

	_, err = tempFile.Write(content)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	h, err := New(WithPath(tempFile.Name()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = h.Load()
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}

func TestSave(t *testing.T) {
	tempPath, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempPath.Name())
	defer tempPath.Close()

	h, err := New(WithPath(tempPath.Name()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = h.Add("127.0.0.1", []string{"localhost"}, "test entry")
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	err = h.Save()
	if err != nil {
		t.Errorf("failed to save hosts file: %v", err)
	}

	err = h.Load()
	if err != nil {
		t.Errorf("failed to load hosts file: %v", err)
	}

	if len(h.Entries) != 1 {
		t.Errorf("expected 1 entry, but got %d", len(h.Entries))
	}
}

// TODO: Add more tests for Save method.
