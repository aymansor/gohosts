package hosts

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// AddHost adds a new entry to the os hosts file with the given host. It returns an error
// if the os hosts file location cannot be determined.
func AddHost(host HostEntry) error {
	location, err := getHostsFileLocation()
	if err != nil {
		return err
	}

	return AddHostsFromLocation(location, host)
}

// AddHostsFromLocation appends a new entry to the hosts file at the given path with the given host.
// It returns an error if the IP address is invalid or if any of the hostnames are invalid.
func AddHostsFromLocation(path string, host HostEntry) error {
	if host.IP == nil {
		return fmt.Errorf("invalid IP address")
	}

	for _, hostname := range host.Hostnames {
		if !isValidHostname(hostname) {
			return fmt.Errorf("invalid hostname: %s", hostname)
		}
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
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

// RemoveHosts removes the specified entry from the os hosts file. It returns an error if the os hosts file location cannot be determined.
func RemoveHosts(host HostEntry) error {
	location, err := getHostsFileLocation()
	if err != nil {
		return err
	}

	return RemoveHostsFromLocation(location, host)
}

// RemoveHostsFromLocation removes the specified entry from the hosts file at the given path.
// It returns an error if the IP address is invalid, fails to read the hosts file, or fails to write the updated hosts file.
func RemoveHostsFromLocation(path string, host HostEntry) error {
	if host.IP == nil {
		return fmt.Errorf("invalid IP address")
	}

	lines, err := readHostsFile(path)
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

	file, err := os.Create(path)
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

// sameHostnames checks if two slices of hostnames are equal.
func sameHostnames(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
