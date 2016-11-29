package main

import (
	"net/http"
	"io/ioutil"
	"path/filepath"
)

type testStaticHandler struct {
	http.Handler
}

func getContentType(path string) string {
	var conttype string
	ext := filepath.Ext(path)

	switch ext {
	case ".html":
		conttype = "text/html"
	case ".css":
		conttype = "text/css"
	case ".js":
		conttype = "application/javascript"
	case ".png":
		conttype = "image/png"
	case ".gif":
		conttype = "image/gif"
	case ".jpg":
		conttype = "image/jpeg"
	default:
		conttype = "text/plain"
	}

	return conttype
}

func main() {
	http.Handle("/", new(testStaticHandler))
	http.ListenAndServe(":5001", nil)
}

func (h *testStaticHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	path := "wwwroot" + req.URL.Path
	cont, err := ioutil.ReadFile(path)
	if err != nil {
		wr.WriteHeader(404)
		wr.Write([]byte(http.StatusText(404)))
		return
	}

	contentType := getContentType(path)
	wr.Header().Add("Content-Type", contentType)
	wr.Write(cont)
}
