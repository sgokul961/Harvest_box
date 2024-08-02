package routes

import "github.com/gorilla/mux"

func Routes() *mux.Router {
	router := mux.NewRouter()
	HarvestBoxRoutes(router)
	return router
}
