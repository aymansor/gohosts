package gohosts

import (
	"fmt"
	"net"
)

// Add appends a new host entry to the hosts file.
func (h *HostsFile) Add(ip string, hostname []string, comment string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address: %s", ip)
	}

	if len(hostname) == 0 {
		return fmt.Errorf("no hostnames provided")
	}

	for _, hostname := range hostname {
		if !isValidHostname(hostname) {
			return fmt.Errorf("invalid hostname: %s", hostname)
		}
	}

	entry := HostEntry{
		IP:        ip,
		Hostnames: hostname,
		Comment:   comment,
		Active:    true,
	}
	h.Entries = append(h.Entries, entry)
	return nil
}

// AddBatch appends multiple host entries to the hosts file.
func (h *HostsFile) AddBatch(entries ...HostEntry) error {
	for _, entry := range entries {
		err := h.Add(entry.IP, entry.Hostnames, entry.Comment)
		if err != nil {
			return err
		}
	}

	return nil
}

// Remove deletes a host entry from the hosts file.
func (h *HostsFile) Remove(ip string, hostname []string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address: %s", ip)
	}

	if len(hostname) == 0 {
		return fmt.Errorf("no hostnames provided")
	}

	for _, hostname := range hostname {
		if !isValidHostname(hostname) {
			return fmt.Errorf("invalid hostname: %s", hostname)
		}
	}

	for i, entry := range h.Entries {
		if entry.IP == ip && contains(entry.Hostnames, hostname...) {
			if len(entry.Hostnames) == 1 || len(entry.Hostnames) == len(hostname) {
				// Remove the entire entry if it's the only hostname
				h.Entries = append(h.Entries[:i], h.Entries[i+1:]...)
			} else {
				// Remove the specific hostname from the entry
				entry.Hostnames = removeStrings(entry.Hostnames, hostname...)
				h.Entries[i] = entry
			}

			return nil
		}
	}
	return fmt.Errorf("host entry not found: IP=%s, Hostname=%s", ip, hostname)
}

// RemoveBatch deletes multiple host entries from the hosts file.
func (h *HostsFile) RemoveBatch(entries ...HostEntry) error {
	for _, removedEntry := range entries {
		// TODO: maybe add a comment match as well, I don't know seem useless, who knows
		err := h.Remove(removedEntry.IP, removedEntry.Hostnames)
		if err != nil {
			return err
		}
	}

	return nil
}
