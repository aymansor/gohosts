package hosts

// import (
// 	"net"
// 	"os"
// 	"strings"
// 	"testing"
// )

// // TODO: Add more test cases and clean up the tests (what a miss down here)
// func TestAddHost(t *testing.T) {
// 	tempFile, err := os.CreateTemp("", "hosts")
// 	if err != nil {
// 		t.Fatalf("failed to create temporary file: %v", err)
// 	}
// 	defer os.Remove(tempFile.Name())
// 	defer tempFile.Close()

// 	// Test case 1: Valid host entry
// 	host := HostEntry{
// 		IP:        net.ParseIP("192.168.0.1"),
// 		Hostnames: []string{"example.com", "example"},
// 		Comment:   "Test host",
// 	}

// 	err = AddHostsFromLocation(tempFile.Name(), host)
// 	if err != nil {
// 		t.Errorf("unexpected error: %v", err)
// 	}

// 	host.IP = nil

// 	// Test case 2: Invalid IP address
// 	err = AddHostsFromLocation(tempFile.Name(), host)
// 	if err == nil {
// 		t.Error("expected an error, but got nil")
// 	}

// 	// Test case 3: Invalid hostname
// 	host.IP = net.ParseIP("192.168.0.2")
// 	host.Hostnames = []string{"invalid_hostname"}

// 	err = AddHostsFromLocation(tempFile.Name(), host)
// 	if err == nil {
// 		t.Error("expected an error, but got nil")
// 	}
// }

// func TestRemoveHost(t *testing.T) {
// 	tempFile, err := os.CreateTemp("", "hosts")
// 	if err != nil {
// 		t.Fatalf("failed to create temporary file: %v", err)
// 	}
// 	defer os.Remove(tempFile.Name())
// 	defer tempFile.Close()

// 	host := HostEntry{
// 		IP:        net.ParseIP("192.168.0.1"),
// 		Hostnames: []string{"example.com", "example"},
// 		Comment:   "Test host",
// 	}

// 	err = AddHostsFromLocation(tempFile.Name(), host)
// 	if err != nil {
// 		t.Errorf("unexpected error: %v", err)
// 	}

// 	// Test case 1: Remove existing host entry
// 	t.Run("Remove existing host entry", func(t *testing.T) {
// 		err := RemoveHostsFromLocation(tempFile.Name(), host)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}
// 	})

// 	// Test case 2: Remove non-existent host entry
// 	t.Run("Remove non-existent host entry", func(t *testing.T) {
// 		host = HostEntry{
// 			IP:        net.ParseIP("192.168.0.3"),
// 			Hostnames: []string{"nonexistent.com"},
// 		}

// 		err := RemoveHostsFromLocation(tempFile.Name(), host)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}
// 	})
// }

// func TestAddAndRemoveHosts(t *testing.T) {
// 	tempFile, err := os.CreateTemp("", "hosts")
// 	if err != nil {
// 		t.Fatalf("failed to create temporary file: %v", err)
// 	}
// 	defer os.Remove(tempFile.Name())
// 	defer tempFile.Close()

// 	hostsFile := `# Test hosts file for unit testing purposes
// 	::1 localhost
// 	ip6-localhost ip6-loopback
// 	127.0.0.1 example.com example # Test host

// 	# Empty line

// 	# Another test host
// 	192.168.0.1 example
// 	`

// 	_, err = tempFile.WriteString(hostsFile)
// 	if err != nil {
// 		t.Fatalf("failed to write to temporary file: %v", err)
// 	}

// 	// Test case 1: Add a new host entry
// 	host := HostEntry{
// 		IP:        net.ParseIP("127.0.0.1"),
// 		Hostnames: []string{"newhost.com", "newhost"},
// 		Comment:   "New host entry",
// 	}

// 	err = AddHostsFromLocation(tempFile.Name(), host)
// 	if err != nil {
// 		t.Errorf("unexpected error: %v", err)
// 	}

// 	// Verify that the new host entry was added
// 	entries, err := readHostsFile(tempFile.Name())
// 	if err != nil {
// 		t.Errorf("unexpected error: %v", err)
// 	}

// 	// parse the entries
// 	parsedEntries, err := parseHostsFile(entries)
// 	if err != nil {
// 		t.Errorf("unexpected error: %v", err)
// 	}

// 	found := false

// 	// check if the new host entry is present
// 	for _, entry := range parsedEntries {
// 		if entry.IP.Equal(host.IP) && sameHostnames(entry.Hostnames, host.Hostnames) {
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		t.Error("new host entry not found in the hosts file")
// 	}
// 	// Test case 2: Remove an existing host entry
// 	host = HostEntry{
// 		IP:        net.ParseIP("127.0.0.1"),
// 		Hostnames: []string{"example.com", "example"},
// 	}

// 	err = RemoveHostsFromLocation(tempFile.Name(), host)
// 	if err != nil {
// 		t.Errorf("unexpected error: %v", err)
// 	}

// 	expectedHostFile := `# Test hosts file for unit testing purposes
// 	::1 localhost
// 	ip6-localhost ip6-loopback

// 	# Empty line

// 	# Another test host
// 	192.168.0.1 example
// 	127.0.0.1 newhost.com newhost #New host entry
// 	`

// 	// Read the contents of the temporary file and compare it with the expected contents
// 	contents, err := os.ReadFile(tempFile.Name())
// 	if err != nil {
// 		t.Fatalf("failed to read temporary file: %v", err)
// 	}

// 	if strings.TrimSpace(string(contents)) != strings.TrimSpace(expectedHostFile) {
// 		t.Errorf("expected file contents do not match actual contents \n expected: |%s| \n actual: |%s|", strings.TrimSpace(expectedHostFile), strings.TrimSpace(string(contents)))
// 	}
// }
