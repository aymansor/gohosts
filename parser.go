package gohosts

import (
	"bufio"
	"net"
	"os"
	"strings"
)

func (h *HostsFile) readHosts() ([]string, error) {
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

func (h *HostsFile) parseHosts(lines []string) ([]HostEntry, error) {
	var entries []HostEntry

	for _, line := range lines {
		originalLine := line // Keep the original line for additional content purposes
		line = strings.TrimSpace(line)

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Identify if the line is a comment and potentially a valid but inactive entry
		isActive := true
		if line[0] == '#' {
			trimmedLine := strings.TrimSpace(line[1:])
			// If after trimming it looks like a valid entry (has space), and the first field
			// is a valid IP, then it's an inactive host entry
			if strings.Contains(trimmedLine, " ") && net.ParseIP(strings.Fields(trimmedLine)[0]) != nil {
				line = trimmedLine
				isActive = false
			} else {
				// Otherwise, it's just a comment line, add it to the additional content and skip
				h.AditionalContent += originalLine + "\n"
				continue
			}
		}

		var hostnames []string
		var comment string

		commentIndex := strings.Index(line, "#")
		// If there's a comment, separate it from the line
		if commentIndex != -1 {
			comment = strings.TrimSpace(line[commentIndex+1:])
			line = strings.TrimSpace(line[:commentIndex])
		}

		parts := strings.Fields(line)
		// If there's no hostname, skip
		if len(parts) < 2 {
			continue
		}

		// The first part should be the IP address
		ip := net.ParseIP(parts[0])
		if ip == nil {
			continue
		}

		// Finally, the rest of the parts are the hostnames
		hostnames = parts[1:]

		entry := HostEntry{
			IP:        ip.String(),
			Hostnames: hostnames,
			Comment:   comment,
			Active:    isActive,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
