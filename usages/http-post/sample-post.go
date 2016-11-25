package main

import (
	"bytes"
	"net/http"
	"io/ioutil"
)

func main() {
	// 간단한 http.Post 예제
	reqBody := bytes.NewBufferString("Post plain text")
	resp, err := http.Post("http://httpbin.org/post", "text/plain", reqBody)
	//resp, err := http.Post("http://httpbin.org/post", "application/json", reqBody)
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
