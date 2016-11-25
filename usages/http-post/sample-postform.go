package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
)

func main() {
	// 간단한 http.PostForm 예제
	resp, err := http.PostForm("http://httpbin.org/post", url.Values{"Name": {"Daniel"}, "Age": {"41"}})
	//resp, err := http.PostForm("http://121.78.237.136:9091/boards/api_get_gallery", url.Values{"bmode": {"gallery"}})
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
