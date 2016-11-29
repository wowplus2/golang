package main

import "net/http"

type testHttpHandler struct {
	http.Header
}

func main() {
	http.Handle("/", new(testHttpHandler))

	http.ListenAndServe(":5001", nil)
}

func (h *testHttpHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	str := "Your request path is " + req.URL.Path
	wr.Write([]byte(str))
}
