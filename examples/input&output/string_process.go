package main

import (
	"strings"
	"strconv"
	"bufio"
	"os"
	"fmt"
)

func cleanString(stream string, seperator string) []int {
	// Trims the stream and then splits
	trimmed_stream := strings.TrimSpace(stream)
	split_arr := strings.Split(trimmed_stream, seperator)
	// convert strings to integers and store them in a slice
	clean_arr := make([]int, len(split_arr))

	for k, v := range split_arr {
		clean_arr[k], _ = strconv.Atoi(v)
	}

	return clean_arr
}

func main() {
	ir := bufio.NewReader(os.Stdin)
	in, _ := ir.ReadString('\n')
	noOfTestCase, _ := strconv.Atoi(strings.TrimSpace(in))

	for i := 0; i < noOfTestCase; i++ {
		linp, _ := ir.ReadString('\n')
		fmt.Println(cleanString(linp, " "))
	}
}
