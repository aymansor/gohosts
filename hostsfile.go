package hosts

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type HostEntry struct {
	IP        string
	Hostnames []string
	Comment   string
	Active    bool
}

type HostsFile struct {
	path             string
	Entries          []HostEntry
	AditionalContent string
}

type HostsOption func(*HostsFile)

func WithPath(path string) HostsOption {
	return func(h *HostsFile) {
		h.path = path
	}
}

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

func (h *HostsFile) Save() error {
	backup := fmt.Sprintf("%s_%s_%s.bak", h.path, "gohosts", time.Now().Format("20060102150405"))
	err := copyFile(h.path, backup)
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
