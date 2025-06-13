package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <hostname>\n", os.Args[0])
	}

	hostname := os.Args[1]

	ips, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs for %s: %v\n", hostname, err)
	} else {
		fmt.Printf("IP addresses for %s:\n", hostname)
		for _, ip := range ips {
			fmt.Println(ip.String())
		}
	}
}
