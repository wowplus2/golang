package main

import (
	"encoding/xml"
	"bytes"
	"net/http"
	"io/ioutil"
)

// Perspn -
type Person struct {
	Name string
	Age int
}

func main() {
	person := Person{"Daniel", 32}
	pbytes, _ := xml.Marshal(person)
	buff := bytes.NewBuffer(pbytes)

	resp, err := http.Post("http://httpbin.org/post", "application/xml", buff)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Response 값 체크
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
	}
}
