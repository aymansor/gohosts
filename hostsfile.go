package hosts

import (
	"fmt"
	"net"
)

type HostEntry struct {
	IP        net.IP
	Hostnames []string
	Comment   string
}

type HostsFile struct {
	path    string
	Entries []HostEntry
}

type HostsOption func(*HostsFile)

func WithPath(path string) HostsOption {
	return func(h *HostsFile) {
		h.path = path
	}
}

func New(opts ...HostsOption) (*HostsFile, error) {
	defaultPath, err := getHostsFileLocation()
	if err != nil {
		return nil, err
	}

	h := &HostsFile{
		path: defaultPath,
	}

	for _, opt := range opts {
		opt(h)
	}

	if !isValidPath(h.path) {
		return nil, fmt.Errorf("hosts file does not exist: %s", h.path)
	}

	return h, nil
}

func (h *HostsFile) Load() error {
	lines, err := h.readLines()
	if err != nil {
		return err
	}

	entries, err := parseLines(lines)
	if err != nil {
		return err
	}

	h.Entries = entries

	return nil
}
