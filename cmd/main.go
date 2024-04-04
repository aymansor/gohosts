package main

import (
	"fmt"

	hosts "github.com/aymansor/gohosts"
)

func main() {
	h, err := hosts.New(hosts.WithPath("test"))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = h.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	// for _, host := range h.Entries {
	// 	fmt.Println(host)
	// }

	// fmt.Printf("\n\n%s\n", hosts.AditionalContent)

	// Add a new host entry
	// err = h.RemoveHosts(
	// 	hosts.HostEntry{IP: "192.168.0.1", Hostnames: []string{"example.com"}},
	// 	hosts.HostEntry{IP: "10.0.0.1", Hostnames: []string{"test.com"}, Comment: "Valid Entries"},
	// 	hosts.HostEntry{IP: "192.168.0.2", Hostnames: []string{"host3.com", "host1.com"}, Comment: "Multiple Hostnames"},
	// 	hosts.HostEntry{IP: "192.168.0.3", Hostnames: []string{"whitespace.com"}},
	// 	hosts.HostEntry{IP: "192.168.0.4", Hostnames: []string{"inlinecomment.com"}},
	// 	hosts.HostEntry{IP: "192.168.0.5", Hostnames: []string{"mixedwhitespace.com"}},
	// 	hosts.HostEntry{IP: "192.168.0.6", Hostnames: []string{"mixedwhitespace2.com"}},
	// )
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = h.UpdateHost("127.0.0.1", "what", "nice", "wow")
	// if err != nil {
	// 	fmt.Print(err)
	// }

	// err = hosts.RemoveHost("127.0.0.1", "example.com")
	// if err != nil {
	// 	fmt.Print(err)
	// }

	err = h.Save()
	if err != nil {
		fmt.Println(err)
	}
}
