package main

import (
	"strconv"
	"fmt"
)

func main() {
	str := "34"
	floatStr := "12.34"
	truth := "false"

	// string to int conversion
	intNum, _ := strconv.Atoi(str)

	// Use ParseFloat and ParseBool for other data types
	floatNum, _ := strconv.ParseFloat(floatStr, 64)
	boolVal, _ := strconv.ParseBool(truth)

	fmt.Println(intNum, floatNum, boolVal)
}
