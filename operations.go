package hosts

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (h *HostsFile) AddEntry(host HostEntry) error {
	if host.IP == nil {
		return fmt.Errorf("invalid IP address")
	}

	for _, hostname := range host.Hostnames {
		if !isValidHostname(hostname) {
			return fmt.Errorf("invalid hostname: %s", hostname)
		}
	}

	file, err := os.OpenFile(h.path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	entry := fmt.Sprintf("%s %s", host.IP, strings.Join(host.Hostnames, " "))
	if host.Comment != "" {
		entry += " #" + host.Comment
	}
	_, err = file.WriteString(entry + "\n")
	return err
}

func (h *HostsFile) RemoveEntry(host HostEntry) error {
	if host.IP == nil {
		return fmt.Errorf("invalid IP address")
	}

	// TODO: find better way
	lines, err := h.readLines()
	if err != nil {
		return err
	}

	var newLines []string
	for _, line := range lines {
		entry, err := parseHostEntry(line)
		if err != nil {
			newLines = append(newLines, line)
			continue
		}

		if !entry.IP.Equal(host.IP) {
			newLines = append(newLines, line)
			continue
		}

		if !sameHostnames(entry.Hostnames, host.Hostnames) {
			newLines = append(newLines, line)
		}
	}

	file, err := os.Create(h.path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range newLines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
