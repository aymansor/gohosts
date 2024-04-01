package hosts

import (
	"net"
	"reflect"
	"testing"
)

func TestParseHostsFile(t *testing.T) {
	// Test case 1: Successful parsing
	t.Run("Successful parsing", func(t *testing.T) {
		input := []string{
			"# Comment",
			"127.0.0.1 localhost",
			"::1 localhost",
			"192.168.0.1 example.com example # Comment",
			"invalid line",
		}

		expected := []HostEntry{
			{IP: net.ParseIP("127.0.0.1"), Hostnames: []string{"localhost"}},
			{IP: net.ParseIP("::1"), Hostnames: []string{"localhost"}},
			{IP: net.ParseIP("192.168.0.1"), Hostnames: []string{"example.com", "example"}, Comment: "Comment"},
		}

		entries, err := parseHostsFile(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(entries, expected) {
			t.Errorf("parsed entries do not match expected entries")
		}
	})

	// Test case 2: Valid hostnames
	t.Run("Valid hostnames", func(t *testing.T) {
		input := []string{
			"127.0.0.1 localhost",
			"::1 localhost",
			"192.168.0.1 example.com valid-hostname",
		}
		expected := []HostEntry{
			{IP: net.ParseIP("127.0.0.1"), Hostnames: []string{"localhost"}},
			{IP: net.ParseIP("::1"), Hostnames: []string{"localhost"}},
			{IP: net.ParseIP("192.168.0.1"), Hostnames: []string{"example.com", "valid-hostname"}},
		}
		entries, err := parseHostsFile(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(entries, expected) {
			t.Errorf("parsed entries do not match expected entries")
		}
	})

	// Test case 3: Invalid hostnames
	t.Run("Invalid hostnames", func(t *testing.T) {
		{
			input := []string{
				"127.0.0.1 invalid_hostname",
				"192.168.0.1 too.long.hostname.exceeding.63.characters.in.a.single.label.asjkdfnasdklnfklasdjnfkjsdanfkjsdanfkjsdanfjkasdnkfjnasdkfjnasdksdj",
				"127.0.0.1 toolonghostnameexceeding255charactersmeoqwfmnoqiwenfoieqwnfonewofnqweofinwweqopfnweqofnqweadkjgpoadnfiopguhriuwehiguhweriughriweuhguierhwgiuerhwiguherwiughreuiwhguierwhguiwerhguierhwgiuhweriguheruiwghuierwhgiuwerhgiuerwhgiuherwiguherwiughreuiwghiuerwhguierw",
			}
			expected := []HostEntry{
				{IP: net.ParseIP("127.0.0.1"), Hostnames: []string{}},
				{IP: net.ParseIP("192.168.0.1"), Hostnames: []string{}},
				{IP: net.ParseIP("127.0.0.1"), Hostnames: []string{}},
			}
			entries, err := parseHostsFile(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(entries, expected) {
				t.Errorf("parsed entries do not match expected entries")
			}
		}
	})

	// test case 4: Empty hostname
	t.Run("Empty hostname", func(t *testing.T) {
		input := []string{
			"127.0.0.1",
		}

		entry, err := parseHostsFile(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(entry) != 0 {
			t.Errorf("expected no entries, got %d", len(entry))
		}
	})

	// Test case 5: Invalid IP
	t.Run("Invalid IP", func(t *testing.T) {
		input := []string{
			"234534 example.com",
		}

		entries, err := parseHostsFile(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(entries) != 0 {
			t.Errorf("expected no entries, got %d", len(entries))
		}
	})

	// Test case 6: Empty file
	t.Run("Empty file", func(t *testing.T) {
		input := []string{}

		entries, err := parseHostsFile(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(entries) != 0 {
			t.Errorf("expected no entries, got %d", len(entries))
		}
	})
}
