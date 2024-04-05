# GoHosts

GoHosts is a Golang library for parsing and managing hosts files.

## Features

- Cross-platform support (Windows, Linux, macOS)
- Add and remove host entries
- Parse the contents of the hosts file
- Create a backup of the hosts file
- Restore the hosts file from a backup

## Installation

To install GoHosts, use the following command:

```bash
go get github.com/aymansor/gohosts
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/aymansor/gohosts"
)

func main() {
    // The default hosts file path is the system hosts file, no need to specify it
    hosts, err := gohosts.New()
    if err != nil {
        panic(err)
    }

    // Load the hosts file
    err = hosts.Load()
    if err != nil {
        panic(err)
    }

    // Access the hosts entries parsed from the hosts file
    for _, host := range hosts.Entries {
        fmt.Printf("IP: %s, Hostnames: %v, Comment: %s\n", host.IP, host.Hostnames, host.Comment)
    }

    // Add a snigle host
    err = hosts.Add("127.0.0.1", []string{"example.com"}, "optional comment")
    if err != nil {
        panic(err)
    }

    // Add multiple hosts
    err = hosts.AddBatch(
        gohosts.HostEntry{IP: "192.168.0.1", Hostnames: []string{"host1.com", "host2.com"}},
        gohosts.HostEntry{IP: "192.168.0.2", Hostnames: []string{"router"}, Comment: "wow"},
    )
    if err != nil {
        panic(err)
    }

    // Remove a host
    err = hosts.Remove("127.0.0.1", []string{"example.com"})
    if err != nil {
        panic(err)
    }

    // Remove multiple hosts
    err = hosts.RemoveBatch(
        gohosts.HostEntry{IP: "192.168.0.1", Hostnames: []string{"host1.com", "host2.com"}},
        gohosts.HostEntry{IP: "192.168.0.2", Hostnames: []string{"router"}, Comment: "wow"},
    )
    if err != nil {
        panic(err)
    }

    // Finally save the changes
    err = hosts.Save()
    if err != nil {
        panic(err)
    }
}

```
