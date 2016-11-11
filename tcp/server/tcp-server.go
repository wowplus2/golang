package main

import (
	"net"
	"fmt"
	"bufio"
	"strings"	// only needed below for sample processing
)

func main() {
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8200")

	// accept connection on port
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline(\n)
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(msg))
		// sample process for string received
		new_msg := strings.ToUpper(msg)
		// send new string back to client
		conn.Write([]byte(new_msg + "\n"))
	}

}
