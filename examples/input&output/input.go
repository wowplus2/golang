package main

import (
	"bufio"
	"os"
	"fmt"
)

func main() {
	ir := bufio.NewReader(os.Stdin)

	inp, _ := ir.ReadString('\n')
	fmt.Print(inp)
}
