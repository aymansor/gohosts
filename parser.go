package hosts

import (
	"net"
	"regexp"
	"strings"
)

// HostEntry represents a single entry in a hosts file.
type HostEntry struct {
	IP        net.IP
	Hostnames []string
	Comment   string
}

// parseHostsFile takes a slice of strings representing the lines of a hosts file
// and parses each line into a HostEntry. It returns a slice of HostEntry.
// It ignores empty lines and lines starting with a "#" (comments).
func parseHostsFile(lines []string) ([]HostEntry, error) {
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

func isValidHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 255 {
		return false
	}

	for _, label := range strings.Split(hostname, ".") {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		if !regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`).MatchString(label) {
			return false
		}
	}

	return true
}
