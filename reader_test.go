package hosts

import (
	"os"
	"testing"
)

func TestReadHostsFile(t *testing.T) {
	// Test case 1: Successful file read
	t.Run("Successful file read", func(t *testing.T) {
		file, err := os.CreateTemp("", "hosts")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(file.Name())

		content := "127.0.0.1 localhost\n::1 localhost"
		_, err = file.WriteString(content)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()

		lines, err := readHostsFile(file.Name())
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := []string{"127.0.0.1 localhost", "::1 localhost"}
		if len(lines) != len(expected) {
			t.Errorf("Expected %d lines, but got %d", len(expected), len(lines))
		}
		for i, line := range lines {
			if line != expected[i] {
				t.Errorf("Expected line '%s', but got '%s'", expected[i], line)
			}
		}
	})

	// Test case 2: Error when opening file
	t.Run("Error when opening file", func(t *testing.T) {
		// Call the function with a non-existent file
		_, err := readHostsFile("/path/to/nonexistent/file")
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
	})

	// Test case 3: Error when reading file
	t.Run("Error when reading file", func(t *testing.T) {
		file, err := os.CreateTemp("", "hosts")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(file.Name())

		// Make the file unreadable
		err = os.Chmod(file.Name(), 0000) // 0000: no permissions
		if err != nil {
			t.Fatal(err)
		}

		_, err = readHostsFile(file.Name())
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
	})

	// Test case 4: Empty file
	t.Run("Empty file", func(t *testing.T) {
		file, err := os.CreateTemp("", "hosts")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(file.Name())

		lines, err := readHostsFile(file.Name())
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(lines) != 0 {
			t.Errorf("Expected 0 lines, but got %d", len(lines))
		}
	})

	// Test case 5: Scanner error
	t.Run("Scanner error", func(t *testing.T) {
		file, err := os.CreateTemp("", "hosts")
		if err != nil {
			t.Fatalf("failed to create temporary file: %v", err)
		}
		defer os.Remove(file.Name())
		defer file.Close()

		content := make([]byte, 65536+1) // 64KB + 1 byte to exceed default buffer size.
		for i := range content {
			content[i] = 'a' // Fill with a repeating character.
		}

		_, err = file.Write(content)
		if err != nil {
			t.Fatalf("failed to write to temporary file: %v", err)
		}

		_, err = readHostsFile(file.Name())
		if err == nil {
			t.Error("expected an error, but got nil")
		}
	})
}
