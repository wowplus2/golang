package main

import "net/http"

func main() {
	http.Handle("/", http.FileServer(http.Dir("wwwroot")))
	// http.Handle("/static", http.FileServer(http.Dir("wwwroot")))
	http.ListenAndServe(":5001", nil)
}
