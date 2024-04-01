# GoHosts

GoHosts is a Golang library for parsing and managing hosts files.

## Features

- Cross-platform support (Windows, Linux, macOS)
- Add and remove host entries
- Parse the contents of the hosts file
- Create a backup of the hosts file
- Restore the hosts file from a backup

## To-Do List
- [ ] Write more unit tests to cover more cases (improve tests)
- [ ] Add proper error handling and improve input validation.
- [ ] Improve the API to make it easy to use.
- [ ] Document common use cases and finish the examples in the README.
- [ ] Add documentation.
- [ ] Add support to disallow duplicate entries in the hosts file.
- [ ] Fix the issue where a new entry isn't added on a new line if the last line doesn't end with a newline
- [ ] Add examples demonstrating how to use the library.
- [ ] Add an option to enable/disable an entry.

## Installation

To install GoHosts, use the following command:

```bash
go get github.com/aymansor/gohosts
```

## Usage

Import the GoHosts package in your Golang code:

```go
import gohosts "github.com/aymansor/gohosts"
```

### Adding a Host Entry

To add a new host entry to the hosts file, use the AddHost function:

```go
entry := hosts.HostEntry{
    IP:        net.ParseIP("192.168.0.1"),
    Hostnames: []string{"example.com", "www.example.com"},
    Comment:   "Example host entry", // optional
}

err := hosts.AddHost(entry)
if err != nil {
    // Handle the error
}
```

### Removing a Host Entry

```go
entry := hosts.HostEntry{
    IP:        net.ParseIP("192.168.0.1"),
    Hostnames: []string{"example.com", "www.example.com"},
}

err := hosts.RemoveHosts(entry)
if err != nil {
    // Handle the error
}
```

### Parsing the Hosts File

To parse the contents of the hosts file, use the ParseHostsFile function:

```go
TODO: add code snippet
```

### Creating a Backup

To create a backup of the hosts file, use the CreateBackup function:

```go
err := hosts.CreateBackup()
if err != nil {
    // Handle the error
}
```

### Restoring from a Backup

To restore the hosts file from a backup, use the RestoreBackup function:

```go
err := hosts.RestoreBackup()
if err != nil {
    // Handle the error
}
```