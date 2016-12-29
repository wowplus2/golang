package main

import (
	"os"
	"fmt"
	"net"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s dotted-ip-addr\n", os.Args[0])
		os.Exit(0)
	}

	dotAddr := os.Args[1]
	addr := net.ParseIP(dotAddr)
	if addr == nil {
		fmt.Println("Invalid address")
		os.Exit(1)
	}

	mask := addr.DefaultMask()
	network := addr.Mask(mask)
	ones, bits := mask.Size()

	fmt.Println("Address is ", addr.String())
	fmt.Println("\tDefault mask length is ", bits)
	fmt.Println("\tLeading ones count is ", ones)
	fmt.Println("\tMask is ", mask.String())
	fmt.Println("\tNetwork is ", network.String())
	os.Exit(0)
}
