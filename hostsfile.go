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

	// Create a temporary file to write the modified hosts file instead of modifying the original file
	// directly. This is done to prevent data loss in case of an error while writing the file.
	tempFile, err := os.CreateTemp("", "hosts-")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Start by writing the additional content to the temporary file
	// This is the content that was not parsed as host entries
	writer := bufio.NewWriter(tempFile)
	_, err = writer.WriteString(h.AditionalContent)
	if err != nil {
		return fmt.Errorf("failed to write additional content to temporary file: %v", err)
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
		// Then write the host entry to the temporary file
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write entry to temporary file: %v", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush writer: %v", err)
	}

	// Finally replace the original hosts file with the temporary file
	err = os.Rename(tempFile.Name(), h.path)
	if err != nil {
		return fmt.Errorf("failed to replace hosts file: %v", err)
	}

	return nil
}
