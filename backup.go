package gohosts

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	BackupFileInfix = "gohosts"
)

// CreateBackup creates a backup of the hosts file with the format <path>_<BackupFileInfix>_<timestamp>.bak
func (h *HostsFile) CreateBackup() error {
	backup := fmt.Sprintf("%s_%s_%s.bak", h.path, BackupFileInfix, time.Now().Format("20060102150405"))
	err := copyFile(h.path, backup)
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	return nil
}

// RestoreBackup restores the hosts file from the latest backup file or a specific backup file
// based on the rollback count if provided (default is 1)
func (h *HostsFile) RestoreBackup(rollback ...int) error {
	var rollbackCount int
	// Only one argument is expected
	if len(rollback) > 0 {
		rollbackCount = rollback[0]
	} else {
		// Default rollback count is 1
		rollbackCount = 1
	}

	if rollbackCount < 1 {
		return fmt.Errorf("rollback count must be greater than 0")
	}

	backupFiles, err := getBackupFiles(h.path)
	if err != nil {
		return err
	}

	if rollbackCount > len(backupFiles) {
		return fmt.Errorf("rollback count is greater than the number of backup files")
	}

	// The backup file to restore
	backupFile := backupFiles[len(backupFiles)-rollbackCount]

	err = copyFile(filepath.Join(filepath.Dir(h.path), backupFile), h.path)
	if err != nil {
		return fmt.Errorf("failed to restore backup: %v", err)
	}

	return nil
}

// getBackupFiles returns a list of backup files for the hosts file
func getBackupFiles(path string) ([]string, error) {
	files, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	var backupFiles []string
	for _, file := range files {
		// Check if the file is a backup file with the format <path>_<BackupFileInfix>_<timestamp>.bak
		if strings.HasPrefix(file.Name(), filepath.Base(path)+"_"+BackupFileInfix+"_") && strings.HasSuffix(file.Name(), ".bak") {
			backupFiles = append(backupFiles, file.Name())
		}
	}

	// Mayeb unnecessary, if no backup files are found, then it's just an empty slice
	// if len(backupFiles) == 0 {
	// 	return nil, fmt.Errorf("no backup files found")
	// }

	return backupFiles, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
