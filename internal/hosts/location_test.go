package hosts

import (
	"testing"
)

func TestGetHostsFileLocation(t *testing.T) {
	location, err := getHostsFileLocation()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fileExists(location) {
		t.Errorf("hosts file not found at location: %s", location)
	}
}

func TestFileExists(t *testing.T) {
	if fileExists("nonexistent") {
		t.Error("fileExists returned true for nonexistent file")
	}
}

func TestFileExistsWithEmptyFilename(t *testing.T) {
	if fileExists("") {
		t.Error("fileExists returned true for empty filename")
	}
}
