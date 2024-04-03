package main

import (
	"fmt"

	hosts "github.com/aymansor/gohosts"
)

func main() {
	h, err := hosts.New(
		hosts.WithPath("s"),
	)
	if err != nil {
		fmt.Printf("unexpected error: %v\n", err)
		return
	}

	err = h.Load()
	if err != nil {
		fmt.Printf("failed to load hosts file: %v\n", err)
	}

	for _, entry := range h.Entries {
		fmt.Printf("%s %s %s\n", entry.IP, entry.Hostnames, entry.Comment)
	}
}
