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
	// 테스트용 JSON 데이터
	jsonBytes, _ := json.Marshal(Member{"Piljun", 4, true})

	// JSON 디코딩
	var mem Member
	err := json.Unmarshal(jsonBytes, &mem)
	if err != nil {
		panic(err)
	}

	// mem 구조체 필드 엑세스
	fmt.Println("sample-json_decoding")
	fmt.Println("\tBYTES :", jsonBytes)
	fmt.Println("\tmem.Name :", mem.Name)
	fmt.Println("\tmem.Age :", mem.Age)
	fmt.Println("\tmem.Active :", mem.Active)
}
