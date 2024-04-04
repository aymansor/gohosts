package hosts

import (
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	tempFile, err := os.CreateTemp("", "hosts_test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	h := &HostsFile{
		path: tempFile.Name(),
	}

	// Test adding a valid entry
	err = h.Add("127.0.0.1", []string{"localhost"}, "Test comment")
	if err != nil {
		t.Errorf("Error adding valid entry: %v", err)
	}

	// Test adding an invalid IP address
	err = h.Add("invalid_ip", []string{"localhost"}, "")
	if err == nil {
		t.Error("Expected an error for invalid IP address")
	}

	// Test adding an entry with no hostnames
	err = h.Add("127.0.0.1", []string{}, "")
	if err == nil {
		t.Error("Expected an error for entry with no hostnames")
	}

	// Test adding an entry with an invalid hostname
	err = h.Add("127.0.0.1", []string{"invalid_hostname!"}, "")
	if err == nil {
		t.Error("Expected an error for entry with an invalid hostname")
	}
}

func TestAddBatch(t *testing.T) {
	tempFile, err := os.CreateTemp("", "hosts_test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	h := &HostsFile{
		path: tempFile.Name(),
	}

	// Test adding multiple valid entries
	entries := []HostEntry{
		{IP: "127.0.0.1", Hostnames: []string{"localhost"}, Comment: "Entry 1"},
		{IP: "192.168.0.1", Hostnames: []string{"router"}, Comment: "Entry 2"},
	}
	err = h.AddBatch(entries...)
	if err != nil {
		t.Errorf("Error adding multiple valid entries: %v", err)
	}

	// Test adding an entry with an invalid IP address
	invalidEntries := []HostEntry{
		{IP: "invalid_ip", Hostnames: []string{"localhost"}, Comment: ""},
	}
	err = h.AddBatch(invalidEntries...)
	if err == nil {
		t.Error("Expected an error for entry with an invalid IP address")
	}

	expectedEntries := []HostEntry{
		{IP: "127.0.0.1", Hostnames: []string{"localhost"}, Comment: "Entry 1", Active: true},
		{IP: "192.168.0.1", Hostnames: []string{"router"}, Comment: "Entry 2", Active: true},
	}

	if !compareEntries(h.Entries, expectedEntries) {
		t.Errorf("Entries do not match expected entries. \nGot: %v, \nExpected: %v", h.Entries, expectedEntries)
	}
}

func TestRemove(t *testing.T) {
	tempFile, err := os.CreateTemp("", "hosts_test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	h := &HostsFile{
		Entries: []HostEntry{
			{IP: "127.0.0.1", Hostnames: []string{"localhost", "local"}, Comment: "Test entry", Active: false},
			{IP: "192.168.0.1", Hostnames: []string{"router"}, Comment: "Router entry"},
		},
		path: tempFile.Name(),
	}

	// Test removing a specific hostname from an entry
	err = h.Remove("127.0.0.1", []string{"local"})
	if err != nil {
		t.Errorf("Error removing specific hostname: %v", err)
	}

	// Test removing the entire entry
	err = h.Remove("192.168.0.1", []string{"router"})
	if err != nil {
		t.Errorf("Error removing entire entry: %v", err)
	}

	// Test removing a non-existent entry
	err = h.Remove("10.0.0.1", []string{"nonexistent"})
	if err == nil {
		t.Error("Expected an error for removing a non-existent entry")
	}

	// Test removing an entry with an invalid IP address
	err = h.Remove("invalid_ip", []string{"localhost"})
	if err == nil {
		t.Error("Expected an error for removing an entry with an invalid IP address")
	}

	// Test removing an entry with no hostnames
	err = h.Remove("127.0.0.1", []string{})
	if err == nil {
		t.Error("Expected an error for removing an entry with no hostnames")
	}

	// Test removing an entry with an invalid hostname
	err = h.Remove("127.0.0.1", []string{"invalid_hostname!"})
	if err == nil {
		t.Error("Expected an error for removing an entry with an invalid hostname")
	}

	expectedEntries := []HostEntry{
		{IP: "127.0.0.1", Hostnames: []string{"localhost"}, Comment: "Test entry", Active: false},
	}

	if !compareEntries(h.Entries, expectedEntries) {
		t.Errorf("Entries do not match expected entries. \nGot: %v, \nExpected: %v", h.Entries, expectedEntries)
	}
}

func TestRemoveBatch(t *testing.T) {
	tempFile, err := os.CreateTemp("", "hosts_test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	h := &HostsFile{
		Entries: []HostEntry{
			{IP: "127.0.0.1", Hostnames: []string{"localhost", "local"}, Comment: "Test entry"},
			{IP: "192.168.0.1", Hostnames: []string{"router"}, Comment: "Router entry"},
		},
		path: tempFile.Name(),
	}

	// Test removing multiple entries
	entriesToRemove := []HostEntry{
		{IP: "127.0.0.1", Hostnames: []string{"local"}},
		{IP: "192.168.0.1", Hostnames: []string{"router"}},
	}
	err = h.RemoveBatch(entriesToRemove...)
	if err != nil {
		t.Errorf("Error removing multiple entries: %v", err)
	}

	// Test removing an entry with an invalid IP address
	invalidEntriesToRemove := []HostEntry{
		{IP: "invalid_ip", Hostnames: []string{"localhost"}},
	}
	err = h.RemoveBatch(invalidEntriesToRemove...)
	if err == nil {
		t.Error("Expected an error for removing an entry with an invalid IP address")
	}

	expectedEntries := []HostEntry{
		{IP: "127.0.0.1", Hostnames: []string{"localhost"}, Comment: "Test entry"},
	}

	if !compareEntries(h.Entries, expectedEntries) {
		t.Errorf("Entries do not match expected entries. \nGot: %v, \nExpected: %v", h.Entries, expectedEntries)
	}
}
