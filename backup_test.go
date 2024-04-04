package gohosts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateBackup(t *testing.T) {
	tempFile, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	hostsFile := &HostsFile{path: tempFile.Name()}

	err = hostsFile.CreateBackup()
	if err != nil {
		t.Errorf("Failed to create backup: %v", err)
	}

	backupFiles, err := getBackupFiles(hostsFile.path)
	if err != nil {
		t.Errorf("Failed to get backup files: %v", err)
	}
	if len(backupFiles) != 1 {
		t.Errorf("Expected 1 backup file, got %d", len(backupFiles))
	}
}

func TestCreateBackup_InValidPath(t *testing.T) {
	hostsFile := &HostsFile{path: "/invalid/path"}

	err := hostsFile.CreateBackup()
	if err == nil {
		t.Error("Expected an error for invalid path")
	}
}

func TestRestoreBackup(t *testing.T) {
	tempFile, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	hostsFile := &HostsFile{path: tempFile.Name()}

	_, err = tempFile.WriteString("Hello, World!\nThis is a test, wow!")
	if err != nil {
		t.Fatalf("Failed to write to hosts file: %v", err)
	}

	err = hostsFile.CreateBackup()
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	err = hostsFile.RestoreBackup()
	if err != nil {
		t.Errorf("Failed to restore backup: %v", err)
	}

	content, err := os.ReadFile(hostsFile.path)
	if err != nil {
		t.Fatalf("Failed to read hosts file: %v", err)
	}

	if string(content) != "Hello, World!\nThis is a test, wow!" {
		t.Errorf("Unexpected content in hosts file: %s", string(content))
	}

	// Test restoring a specific backup
	err = hostsFile.RestoreBackup(1)
	if err != nil {
		t.Errorf("Failed to restore specific backup: %v", err)
	}

	// Test restoring with an invalid rollback count
	err = hostsFile.RestoreBackup(-1)
	if err == nil {
		t.Error("Expected an error for invalid rollback count")
	}

	// Test restoring with a rollback count greater than the number of backup files
	err = hostsFile.RestoreBackup(2)
	if err == nil {
		t.Error("Expected an error for rollback count greater than backup files")
	}
}

func TestRestoreBackup_RollbackCount(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "backups")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	hostsPath := filepath.Join(tempDir, "hosts")
	backupFile1 := filepath.Join(tempDir, "hosts_gohosts_20240402000000.bak")
	backupFile2 := filepath.Join(tempDir, "hosts_gohosts_20240403000000.bak")
	err = os.WriteFile(backupFile1, []byte("Backup 1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create backup file: %v", err)
	}
	err = os.WriteFile(backupFile2, []byte("Backup 2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create backup file: %v", err)
	}

	hostsFile := &HostsFile{path: hostsPath}

	err = hostsFile.RestoreBackup(2)
	if err != nil {
		t.Errorf("Failed to restore backup: %v", err)
	}

	content, err := os.ReadFile(hostsFile.path)
	if err != nil {
		t.Fatalf("Failed to read hosts file: %v", err)
	}
	if string(content) != "Backup 1" {
		t.Errorf("Unexpected content in hosts file: %s", string(content))
	}

	// default rollback count is 1
	err = hostsFile.RestoreBackup()
	if err != nil {
		t.Errorf("Failed to restore backup: %v", err)
	}

	content, err = os.ReadFile(hostsFile.path)
	if err != nil {
		t.Fatalf("Failed to read hosts file: %v", err)
	}

	if string(content) != "Backup 2" {
		t.Errorf("Unexpected content in hosts file: %s", string(content))
	}
}

func TestRestoreBackup_InValidPath(t *testing.T) {
	hostsFile := &HostsFile{path: "/invalid/path"}

	err := hostsFile.RestoreBackup()
	if err == nil {
		t.Error("Expected an error for invalid path")
	}
}

func TestRestoreBackup_NoBackupFiles(t *testing.T) {
	tempFile, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	hostsFile := &HostsFile{path: tempFile.Name()}

	err = hostsFile.RestoreBackup()
	if err == nil {
		t.Error("Expected an error for no backup files")
	}
}

func TestGetBackupFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "backups")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	hostsPath := filepath.Join(tempDir, "hosts")
	backupFile1 := filepath.Join(tempDir, "hosts_gohosts_20240402000000.bak")
	backupFile2 := filepath.Join(tempDir, "hosts_gohosts_20240403000000.bak")
	err = os.WriteFile(backupFile1, []byte("Backup 1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create backup file: %v", err)
	}
	err = os.WriteFile(backupFile2, []byte("Backup 2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create backup file: %v", err)
	}

	backupFiles, err := getBackupFiles(hostsPath)
	if err != nil {
		t.Errorf("Failed to get backup files: %v", err)
	}
	if len(backupFiles) != 2 {
		t.Errorf("Expected 2 backup files, got %d", len(backupFiles))
	}
}

func TestCopyFile(t *testing.T) {
	srcFile, err := os.CreateTemp("", "source")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer os.Remove(srcFile.Name())
	defer srcFile.Close()

	_, err = srcFile.WriteString("Hello, World!")
	if err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}

	dstFile, err := os.CreateTemp("", "destination")
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	defer os.Remove(dstFile.Name())
	defer dstFile.Close()

	err = copyFile(srcFile.Name(), dstFile.Name())
	if err != nil {
		t.Errorf("Failed to copy file: %v", err)
	}

	content, err := os.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(content) != "Hello, World!" {
		t.Errorf("Unexpected content in destination file: %s", string(content))
	}
}

func TestCopyFile_SourceNotFound(t *testing.T) {
	dstFile, err := os.CreateTemp("", "destination")
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	defer os.Remove(dstFile.Name())
	defer dstFile.Close()

	err = copyFile("/invalid/source", dstFile.Name())
	if err == nil {
		t.Error("Expected an error for source not found")
	}
}

func TestCopyFile_DestinationNotFound(t *testing.T) {
	srcFile, err := os.CreateTemp("", "source")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer os.Remove(srcFile.Name())
	defer srcFile.Close()

	err = copyFile(srcFile.Name(), "/invalid/destination")
	if err == nil {
		t.Error("Expected an error for destination not found")
	}
}
