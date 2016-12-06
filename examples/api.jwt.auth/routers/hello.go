package routers

import (
	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
	"github.com/wowplus2/golang/examples/api.jwt.auth/controllers"
)

func SetHelloRoutes(router *mux.Router) *mux.Router {
	router.Handle("test/hello",
		negroni.New(
			negroni.HandlerFunc(controllers.HelloController),
	)).Method("GET")

	return router
}
