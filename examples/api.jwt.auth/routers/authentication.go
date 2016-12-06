package routers

import (
	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
	"github.com/wowplus2/golang/examples/api.jwt.auth/controllers"
)

func SetAuthenticationRoutes(router *mux.Route) *mux.Router {
	router.HandleFunc("/token-auth", controllers.Login).Methods("POST")
	router.Handle("/refresh-token-auth",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RefreshToken)),
	).Methods("GET")
	router.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.Logout)),
	).Methods("GET")

	return router
}
