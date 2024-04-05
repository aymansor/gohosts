package gohosts

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// HostEntry represents a single entry in a hosts file.
type HostEntry struct {
	IP        string
	Hostnames []string
	Comment   string
	Active    bool
}

// HostsFile represents a hosts file.
type HostsFile struct {
	path             string
	Entries          []HostEntry
	AditionalContent string
}

// HostsOption is a functional option for configuring a HostsFile.
type HostsOption func(*HostsFile)

// WithPath is a HostsOption that sets the path of the hosts file to be used.
func WithPath(path string) HostsOption {
	return func(h *HostsFile) {
		h.path = path
	}
}

// New creates a new HostsFile with the provided options.
// If no options are provided, the system hosts file is used.
func New(opts ...HostsOption) (*HostsFile, error) {
	defaultPath, err := getHostsFileLocation()
	if err != nil {
		return nil, err
	}

	h := &HostsFile{
		path: defaultPath,
	}

	for _, opt := range opts {
		opt(h)
	}

	if !isValidPath(h.path) {
		return nil, fmt.Errorf("hosts file does not exist: %s", h.path)
	}

	return h, nil
}

// Load reads the hosts file and parses its content.
func (h *HostsFile) Load() error {
	lines, err := h.readHosts()
	if err != nil {
		return err
	}

	entries, err := h.parseHosts(lines)
	if err != nil {
		return err
	}

	h.Entries = entries

	return nil
}

// Save writes the hosts file with the modified content. It creates a backup of the original hosts file
// before writing the modified content.
func (h *HostsFile) Save() error {
	// Before doing anything, create a backup of the hosts file
	err := h.CreateBackup()
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	// Open the hosts file for writing
	file, err := os.OpenFile(h.path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open hosts file: %v", err)
	}
	defer file.Close()

	// Start by writing the additional content to the file
	// This is the content that was not parsed as host entries
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(h.AditionalContent)
	if err != nil {
		// If an error occurs while writing, restore the backup
		if restoreErr := h.RestoreBackup(); restoreErr != nil {
			return fmt.Errorf("failed to write additional content to hosts file: %v, and failed to restore backup: %v", err, restoreErr)
		}
		return fmt.Errorf("failed to write additional content to hosts file: %v", err)
	}

	for _, entry := range h.Entries {
		// TODO: make a constant for the spacing between the columns
		// TODO: also maybe pretty print the entries
		line := fmt.Sprintf("%s     %s", entry.IP, strings.Join(entry.Hostnames, " "))
		if entry.Comment != "" {
			line += "     # " + entry.Comment
		}
		if !entry.Active {
			line = "# " + line
		}
		// Then write the host entry to the file
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			// If an error occurs while writing, restore the backup
			if restoreErr := h.RestoreBackup(); restoreErr != nil {
				return fmt.Errorf("failed to write entry to hosts file: %v, and failed to restore backup: %v", err, restoreErr)
			}
			return fmt.Errorf("failed to write entry to hosts file: %v", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		// If an error occurs while flushing, restore the backup
		if restoreErr := h.RestoreBackup(); restoreErr != nil {
			return fmt.Errorf("failed to flush writer: %v, and failed to restore backup: %v", err, restoreErr)
		}
		return fmt.Errorf("failed to flush writer: %v", err)
	}

	return nil
}
