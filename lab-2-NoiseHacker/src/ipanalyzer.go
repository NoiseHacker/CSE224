package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"math"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 2 && len(os.Args) != 3 {
		log.Fatalf("Usage: %s cidr_block [ip_address]", os.Args[0])
	}

	// os.Args[1] contains the cidr_block
	// os.Args[2] optionally contains the IP address to test

	// Replace the line below and start coding your logic from here.

	// Parse and validates IP and CIDR notation
	_, ipNet, err := net.ParseCIDR(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}

	if len(os.Args) == 2{
		fmt.Printf("Analyzing network: %s\n\n", os.Args[1])

		// Compute Broadcast Address
		broadcast := make(net.IP, len(ipNet.IP))
		copy(broadcast, ipNet.IP)
		for i := range broadcast {
			broadcast[i] = ipNet.IP[i] | ^ipNet.Mask[i]
		}
		// Compute Number of Usable Hosts
		ones, bits := ipNet.Mask.Size()
		if bits != 32 {
			fmt.Println("Invalid IPv4 Mask format")
			return
		}
		diff := float64 (32 - ones)
		hosts := int (math.Pow(2, diff) - 2)

		fmt.Printf("Network address: %s\n", ipNet.IP)
		fmt.Printf("Broadcast address: %s\n", broadcast)
		fmt.Printf("Subnet mask: %s\n", net.IP(ipNet.Mask))
		fmt.Printf("Number of usable hosts: %d\n", hosts)
	} else {
		// Check if a provided IP address is in subnet
		providedIP := net.ParseIP(os.Args[2])
		if providedIP == nil {
			fmt.Println("Invalid IP format")
			return
		}
		fmt.Println(ipNet.Contains(providedIP))
	}
}
