package hosts

import (
	"fmt"
)

// TODO:

// CreateBackup creates a backup of the os hosts file. It returns an error
// if the hosts file location cannot be determined, or if the backup file
// cannot be created.
func CreateBackup() error {
	location, err := getHostsFileLocation()
	if err != nil {
		return err
	}

	err = CreateBackupFromLocation(location)
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	return nil
}

// CreateBackupFromLocation creates a backup of the hosts file at the given path.
// It returns an error if the backup file cannot be created.
func CreateBackupFromLocation(path string) error {
	backupLocation := getBackupLocation(path)

	if !isValidPath(backupLocation) {
		err := copyFile(path, backupLocation)
		if err != nil {
			return fmt.Errorf("failed to create backup: %v", err)
		}
	}

	return nil
}

// RestoreBackup restores the backup of the os hosts file. It returns an error
// if the hosts file location cannot be determined, or if the backup file cannot
// be restored.
func RestoreBackup() error {
	location, err := getHostsFileLocation()
	if err != nil {
		return err
	}

	err = RestoreBackupFromLocation(location)
	if err != nil {
		return fmt.Errorf("failed to restore backup: %v", err)
	}

	return nil
}

// RestoreBackupFromLocation restores the backup of the hosts file at the given path.
// It returns an error if the backup file is not found, or if the backup file cannot
// be restored.
func RestoreBackupFromLocation(path string) error {
	backupLocation := getBackupLocation(path)

	if !isValidPath(backupLocation) {
		return fmt.Errorf("backup file not found")
	}

	err := copyFile(backupLocation, path)
	if err != nil {
		return fmt.Errorf("failed to restore backup: %v", err)
	}

	return nil
}

// getBackupLocation returns the path of the backup file for the given hosts file location.
func getBackupLocation(hostsFileLocation string) string {
	return hostsFileLocation + ".bak"
}
