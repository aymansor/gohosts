package gohosts

import (
	"net"
	"os"
	"strings"
	"testing"
)

const TestHostsData = `
# Additional Comments
#
# This is an additional comment
# Another additional comment
#
# And one more
# Blank Lines


# Basic Functionality

192.168.0.1    example.com # Valid Entries
10.0.0.1       test.com

192.168.0.2    host1.com host2.com host3.com # Multiple Hostnames

  192.168.0.3    whitespace.com      # Leading and Trailing Whitespace
  
192.168.0.4    inlinecomment.com   # This is an inline comment

192.168.0.5    mixedwhitespace.com # Mixed Whitespace
192.168.0.6	mixedwhitespace2.com

999.999.999.999    invalidip.com # Error Handling
abcd               invalidip2.com  # Invalid IP Addresses

192.168.0.7 # Missing Hostname

missingip.com # Missing IP Address

192.168.0.8!           extraneouschar.com # Extraneous Characters
192.168.0.9$    extraneouschar2.com

192.168.0.10    ipv4address.com             # IPv4 Addresses

2001:0db8:85a3:0000:0000:8a2e:0370:7334    ipv6address.com # IPv6 Addresses

127.0.0.1       localhost # Localhost Entries
::1             localhost

192.168.0.11    domainname.com # Domain Names

192.168.0.13    *.wildcard.com # Wildcard Entries

192.168.0.15    café.com # Non-ASCII Characters

192.168.0.16    escapechar.com # Escape Characters

# 127.0.0.1   commentedentry.com
# 127.0.0.1   commentedentry.com  # This is a commented entry with an inline comment

`

func TestReadHosts(t *testing.T) {
	h := &HostsFile{path: "testdata/hosts"}

	_, err := h.readHosts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRaedHosts_InvalidPath(t *testing.T) {
	h := &HostsFile{path: "testdata/invalid"}

	_, err := h.readHosts()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestReadHosts_ScanError(t *testing.T) {
	tempFile, err := os.CreateTemp("", "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := make([]byte, 65536+1) // 64KB + 1 byte to exceed default buffer size.
	for i := range content {
		content[i] = 'a'
	}

	_, err = tempFile.Write(content)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	h := &HostsFile{path: tempFile.Name()}
	_, err = h.readHosts()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestParseHosts(t *testing.T) {
	h := &HostsFile{path: "testdata/hosts"}

	// lines, err := h.readHosts()
	// if err != nil {
	// 	t.Fatalf("unexpected error: %v", err)
	// }

	// split TestHostsData into lines
	lines := make([]string, 0)
	for _, line := range strings.Split(TestHostsData, "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}

	entries, err := h.parseHosts(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []HostEntry{
		{IP: "192.168.0.1", Hostnames: []string{"example.com"}, Comment: "Valid Entries", Active: true},
		{IP: "10.0.0.1", Hostnames: []string{"test.com"}, Comment: "", Active: true},

		{IP: "192.168.0.2", Hostnames: []string{"host1.com", "host2.com", "host3.com"}, Comment: "Multiple Hostnames", Active: true},

		{IP: "192.168.0.3", Hostnames: []string{"whitespace.com"}, Comment: "Leading and Trailing Whitespace", Active: true},

		{IP: "192.168.0.4", Hostnames: []string{"inlinecomment.com"}, Comment: "This is an inline comment", Active: true},

		{IP: "192.168.0.5", Hostnames: []string{"mixedwhitespace.com"}, Comment: "Mixed Whitespace", Active: true},
		{IP: "192.168.0.6", Hostnames: []string{"mixedwhitespace2.com"}, Comment: "", Active: true},

		{IP: "192.168.0.10", Hostnames: []string{"ipv4address.com"}, Comment: "IPv4 Addresses", Active: true},

		{IP: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", Hostnames: []string{"ipv6address.com"}, Comment: "IPv6 Addresses", Active: true},

		{IP: "127.0.0.1", Hostnames: []string{"localhost"}, Comment: "Localhost Entries", Active: true},
		{IP: "::1", Hostnames: []string{"localhost"}, Comment: "", Active: true},

		{IP: "192.168.0.11", Hostnames: []string{"domainname.com"}, Comment: "Domain Names", Active: true},

		{IP: "192.168.0.13", Hostnames: []string{"*.wildcard.com"}, Comment: "Wildcard Entries", Active: true},

		{IP: "192.168.0.15", Hostnames: []string{"café.com"}, Comment: "Non-ASCII Characters", Active: true},

		{IP: "192.168.0.16", Hostnames: []string{"escapechar.com"}, Comment: "Escape Characters", Active: true},

		{IP: "127.0.0.1", Hostnames: []string{"commentedentry.com"}, Comment: "", Active: false},
		{IP: "127.0.0.1", Hostnames: []string{"commentedentry.com"}, Comment: "This is a commented entry with an inline comment", Active: false},
	}

	if len(entries) != len(expected) {
		t.Fatalf("expected %d entries, got %d", len(expected), len(entries))
	}

	for i, entry := range entries {
		if string(net.ParseIP(entry.IP)) != string(net.ParseIP(expected[i].IP)) {
			t.Errorf("expected IP %s, got %s", expected[i].IP, entry.IP)
		}

		if len(entry.Hostnames) != len(expected[i].Hostnames) {
			t.Errorf("expected %d hostnames, got %d", len(expected[i].Hostnames), len(entry.Hostnames))
		}

		for j, hostname := range entry.Hostnames {
			if hostname != expected[i].Hostnames[j] {
				t.Errorf("expected hostname %s, got %s", expected[i].Hostnames[j], hostname)
			}
		}

		if entry.Comment != expected[i].Comment {
			t.Errorf("expected Comment %s, got %s", expected[i].Comment, entry.Comment)
		}

		if entry.Active != expected[i].Active {
			t.Errorf("expected Active %t, got %t", expected[i].Active, entry.Active)
		}
	}

	expectedAdditionalContent := "# Additional Comments\n#\n# This is an additional comment\n# Another additional comment\n#\n# And one more\n# Blank Lines\n# Basic Functionality\n"
	// 	expectedAdditionalContent :=
	// 		`# Additional Comments
	// #
	// # This is an additional comment
	// # Another additional comment
	// #
	// # And one more
	// # Blank Lines
	// # Basic Functionality
	// `

	if h.AditionalContent != expectedAdditionalContent {
		t.Errorf("\n%s\n\n%s", expectedAdditionalContent, h.AditionalContent)
	}
}

func TestParseHosts_Empty(t *testing.T) {
	emptyFile, err := os.CreateTemp("", "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(emptyFile.Name())

	h := &HostsFile{path: emptyFile.Name()}

	lines, err := h.readHosts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := h.parseHosts(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseHosts_OneLine(t *testing.T) {
	tempFile, err := os.CreateTemp("", "one")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("127.0.0.1 localhost")
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	h := &HostsFile{path: tempFile.Name()}
	lines, err := h.readHosts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := h.parseHosts(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].IP != "127.0.0.1" {
		t.Errorf("expected IP 127.0.0.1, got %s", entries[0].IP)
	}

	if len(entries[0].Hostnames) != 1 {
		t.Fatalf("expected 1 hostname, got %d", len(entries[0].Hostnames))
	}

	if entries[0].Hostnames[0] != "localhost" {
		t.Errorf("expected hostname localhost, got %s", entries[0].Hostnames[0])
	}

	if entries[0].Comment != "" {
		t.Errorf("expected empty comment, got %s", entries[0].Comment)
	}

	if entries[0].Active != true {
		t.Errorf("expected active entry, got inactive")
	}

	if h.AditionalContent != "" {
		t.Errorf("expected empty additional content, got %s", h.AditionalContent)
	}

}
