package main

import (
	"encoding/xml"
	"fmt"
)

type Member struct {
	Name string
	Age int
	Active bool
}

func main() {
	mem := Member{"Alex", 10, false}

	xmlBytes, err := xml.Marshal(mem)
	if err != nil {
		panic(err)
	}

	xmlStr := string(xmlBytes)

	fmt.Println("sample-xml_encoding:")
	fmt.Println("\tBYTES :", xmlBytes)
	fmt.Println("\tSTRING :", xmlStr)
}