package main

import "net/http"

func main() {
	http.HandleFunc("/hello", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello Golang World~"))
	})

	http.ListenAndServe(":5001", nil)
}
