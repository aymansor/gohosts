package hosts

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func (h *HostsFile) readLines() ([]string, error) {
	file, err := os.Open(h.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func parseLines(lines []string) ([]HostEntry, error) {
	var entries []HostEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		// TODO: Is this check necessary?
		if len(parts) < 2 {
			continue
		}

		ip := net.ParseIP(parts[0])
		if ip == nil {
			continue
		}

		hostnames := make([]string, 0)
		comment := ""

		for i := 1; i < len(parts); i++ {
			if strings.HasPrefix(parts[i], "#") {
				comment = strings.TrimSpace(strings.TrimPrefix(strings.Join(parts[i:], " "), "#"))
				break
			}
			// TODO: this allows for the possibility of having an empty hostname, check if this is okay
			if isValidHostname(parts[i]) {
				hostnames = append(hostnames, parts[i])
			}
		}

		entry := HostEntry{
			IP:        ip,
			Hostnames: hostnames,
			Comment:   comment,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// parseHostEntry parses a single line from the hosts file and returns a HostEntry.
func parseHostEntry(line string) (HostEntry, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return HostEntry{}, fmt.Errorf("invalid host entry: %s", line)
	}

	ip := net.ParseIP(fields[0])
	if ip == nil {
		return HostEntry{}, fmt.Errorf("invalid IP address: %s", fields[0])
	}

	var hostnames []string
	var comment string
	for i := 1; i < len(fields); i++ {
		if strings.HasPrefix(fields[i], "#") {
			comment = strings.TrimSpace(strings.Join(fields[i:], " "))
			break
		}
		hostnames = append(hostnames, fields[i])
	}

	return HostEntry{IP: ip, Hostnames: hostnames, Comment: comment}, nil
}
