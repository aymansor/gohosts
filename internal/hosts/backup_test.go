package hosts

import (
	"os"
	"testing"
)

// TODO: Add more test cases and clean up the tests
func TestCreateBackup(t *testing.T) {
	// Create a temporary hosts file for testing
	tempFile, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write sample content to the temporary hosts file
	sampleContent := "127.0.0.1 localhost"
	err = os.WriteFile(tempFile.Name(), []byte(sampleContent), 0644)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	// Create a backup
	err = CreateBackupFromLocation(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to create backup: %v", err)
	}

	// Check if the backup file exists
	backupLocation := getBackupLocation(tempFile.Name())
	if !fileExists(backupLocation) {
		t.Error("backup file not found")
	}
}

func TestRestoreBackup(t *testing.T) {
	// Create a temporary hosts file for testing
	tempFile, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write sample content to the temporary hosts file
	sampleContent := "127.0.0.1 localhost"
	err = os.WriteFile(tempFile.Name(), []byte(sampleContent), 0644)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	// Create a backup
	err = CreateBackupFromLocation(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to create backup: %v", err)
	}

	// Modify the hosts file
	modifiedContent := "127.0.0.1 example.com"
	err = os.WriteFile(tempFile.Name(), []byte(modifiedContent), 0644)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	// Restore the backup
	err = RestoreBackupFromLocation(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to restore backup: %v", err)
	}

	// Check if the hosts file content is restored
	restoredContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to read temporary file: %v", err)
	}

	if string(restoredContent) != sampleContent {
		t.Error("hosts file content not restored correctly")
	}
}
