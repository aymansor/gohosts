package gohosts

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// getHostsFileLocation returns the path of the hosts file for the current operating system.
func getHostsFileLocation() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Windows\\System32\\drivers\\etc\\hosts", nil
	case "linux", "darwin":
		return "/etc/hosts", nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// isValidPath checks if the provided path is a valid file path.
func isValidPath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

// isValidHostname checks if the provided hostname is a valid domain name.
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

// contains checks if a slice contains all the provided items.
func contains(slice []string, items ...string) bool {
	for _, item := range items {
		found := false
		for _, s := range slice {
			if s == item {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// removeStrings removes the provided values from the slice.
func removeStrings(slice []string, values ...string) []string {
	var result []string
	for _, s := range slice {
		if !contains(values, s) {
			result = append(result, s)
		}
	}
	return result
}

// compareEntrie checks if two HostEntry are equal.
func compareEntrie(a, b HostEntry) bool {
	if a.IP != b.IP {
		return false
	}
	if len(a.Hostnames) != len(b.Hostnames) {
		return false
	}
	for i := range a.Hostnames {
		if a.Hostnames[i] != b.Hostnames[i] {
			return false
		}
	}
	if a.Comment != b.Comment {
		return false
	}
	if a.Active != b.Active {
		return false
	}

	return true
}

// compareEntries checks if two slices of HostEntry are equal.
func compareEntries(a, b []HostEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !compareEntrie(a[i], b[i]) {
			return false
		}
	}
	return true
}
