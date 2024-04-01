package hosts

import (
	"fmt"
	"os"
	"runtime"
)

func getHostsFileLocation() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return `C:\Windows\System32\drivers\etc\hosts`, nil
	case "linux", "darwin":
		return "/etc/hosts", nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// fileExists returns true if the file exists, false otherwise.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
