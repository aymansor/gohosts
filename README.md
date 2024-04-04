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
- [x] Improve the API to make it easy to use.
- [ ] Document common use cases and finish the examples in the README.
- [ ] Add documentation.
- [ ] Add support to disallow duplicate entries in the hosts file.
- [x] Fix the issue where a new entry isn't added on a new line if the last line doesn't end with a newline
- [ ] Add examples demonstrating how to use the library.
- [x] Add an option to enable/disable an entry.

## Installation

To install GoHosts, use the following command:

```bash
go get github.com/aymansor/gohosts
```

## Usage

TODO