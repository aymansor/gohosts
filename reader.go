package hosts

import (
	"bufio"
	"os"
)

// readHostsFile reads the hosts file at the given location and returns a slice of strings
// representing the lines of the file. It returns an error if the file cannot be read.
func readHostsFile(location string) ([]string, error) {
	file, err := os.Open(location)
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
