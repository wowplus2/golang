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
	// 테스트용 데이터
	xmlBytes, _ := xml.Marshal(Member{"Nick", 45, true})

	// XML 디코딩
	var mem Member
	err := xml.Unmarshal(xmlBytes, &mem)
	if err != nil {
		panic(err)
	}

	// mem 구조체 필드 엑세스
	fmt.Println("sample-xml_decoding: ")
	fmt.Println("\tBYTES :", xmlBytes)
	fmt.Println("\tmem.Name :", mem.Name)
	fmt.Println("\tmem.Age :", mem.Age)
	fmt.Println("\tmem.Active :", mem.Active)

}