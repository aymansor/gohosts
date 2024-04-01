package hosts

import (
	"net"
	"os"
	"testing"
)

func TestAddHost(t *testing.T) {
	// Create a temporary hosts file for testing
	tempFile, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Test case 1: Valid host entry
	host := HostEntry{
		IP:        net.ParseIP("192.168.0.1"),
		Hostnames: []string{"example.com", "example"},
		Comment:   "Test host",
	}

	err = AddHostsFromLocation(tempFile.Name(), host)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	host.IP = nil

	// Test case 2: Invalid IP address
	err = AddHostsFromLocation(tempFile.Name(), host)
	if err == nil {
		t.Error("expected an error, but got nil")
	}

	// Test case 3: Invalid hostname
	host.IP = net.ParseIP("192.168.0.2")
	host.Hostnames = []string{"invalid_hostname"}

	err = AddHostsFromLocation(tempFile.Name(), host)
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}

func TestRemoveHost(t *testing.T) {
	// Create a temporary hosts file for testing
	tempFile, err := os.CreateTemp("", "hosts")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	host := HostEntry{
		IP:        net.ParseIP("192.168.0.1"),
		Hostnames: []string{"example.com", "example"},
		Comment:   "Test host",
	}

	err = AddHostsFromLocation(tempFile.Name(), host)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test case 1: Remove existing host entry
	t.Run("Remove existing host entry", func(t *testing.T) {
		host := HostEntry{
			IP:        net.ParseIP("192.168.0.1"),
			Hostnames: []string{"example.com", "example"},
		}

		err := RemoveHostsFromLocation(tempFile.Name(), host)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Test case 2: Remove non-existent host entry
	t.Run("Remove non-existent host entry", func(t *testing.T) {
		host := HostEntry{
			IP:        net.ParseIP("192.168.0.3"),
			Hostnames: []string{"nonexistent.com"},
		}

		err := RemoveHostsFromLocation(tempFile.Name(), host)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
