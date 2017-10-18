package main

import (
	"fmt"
	"net/http"
)

func main() {
	// "/" 경로로 접속했을때 처리할 핸들러 함수 지정
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome!")
	})

	// "/about" 경로로 접속햇을때 처리할 핸들러 함수 지정
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "about...")
	})

	// 8080 port로 웹서버 구동
	http.ListenAndServe(":8080", nil)
}
