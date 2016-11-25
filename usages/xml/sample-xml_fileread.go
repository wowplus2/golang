package main

import (
	"os"
	"io/ioutil"
	"encoding/xml"
	"fmt"
)

type Member struct {
	Name string
	Age int
	Active bool
}

type Members struct {
	Member []Member
}


func main() {
	// xml 파일 오픈
	fp, err := os.Open("C:\\Temp\\test.xml")
	if err != nil {
		panic(err)
	}

	defer fp.Close()

	// xml 파일 읽기
	data, err := ioutil.ReadAll(fp)

	// xml 디코딩
	var mems Members
	xmlerr := xml.Unmarshal(data, &mems)
	if xmlerr != nil {
		panic(err)
	}

	fmt.Println("sample-xml_fileread: ")
	fmt.Println("\tMembers :", mems)
}
