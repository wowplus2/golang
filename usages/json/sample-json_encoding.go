package main

import (
	"encoding/json"
	"fmt"
)

type Member struct {
	Name string
	Age int
	Active bool
}

func main() {
	// Go 데이터
	mem := Member{"Daniel", 41, true}

	// JSON 인코딩
	jsonBytes, err := json.Marshal(mem)
	if err != nil {
		panic(err)
	}

	// JSON 바이트 문자열로 변경
	jsonStr := string(jsonBytes)

	fmt.Println("sample-json_encoding: ")
	fmt.Println("\tBYTES: ", jsonBytes)
	fmt.Println("\tSTRING: ", jsonStr)
}
