package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

func main() {
	// GET Method 방식 호출
	resp, err := http.Get("http://csharp.news")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// 호출 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(data))
}
