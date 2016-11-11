package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
)


func main() {
	// connection to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8200")

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		txt, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, txt + "\n")
		// listen for reply
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + msg)
	}
}
