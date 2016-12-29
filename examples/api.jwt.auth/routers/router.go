package routers

import "github.com/gorilla/mux"

func InitRoutes() *mux.Reouter {
	router := mux.NewRouter()
	router = SetHelloRoutes(router)
	router = SetAuthenticationRoutes(router)

	return router
}
