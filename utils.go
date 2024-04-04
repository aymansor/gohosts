package hosts

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
)

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

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

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

func removeStrings(slice []string, values ...string) []string {
	var result []string
	for _, s := range slice {
		if !contains(values, s) {
			result = append(result, s)
		}
	}
	return result
}

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
