package main

import (
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
)


// Person -
type Person struct {
	Name string
	Age int
}

func main() {
	person := Person{"Daniel", 41}
	pbytes, _ := json.Marshal(person)
	buff := bytes.NewBuffer(pbytes)
	resp, err := http.Post("http://httpbin.org/post", "application/json", buff)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Response 체크
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
	}
}
