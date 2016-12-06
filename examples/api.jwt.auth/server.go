package main

import (
	"github.com/wowplus2/golang/examples/api.jwt.auth/routers"
	"github.com/codegangsta/negroni"
	"net/http"
)

func main() {
	setting.Init()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(":5000", n)
}
