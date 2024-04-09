package routers

import (
	controller "backend/controller"

	"github.com/gorilla/mux"
)

func SchoolsRouter(router *mux.Router) *mux.Router {

	router.HandleFunc("/", controller.GetAllSchools).Methods("GET")
	router.HandleFunc("/{id}", controller.GetSchoolByID).Methods("GET")
	router.HandleFunc("/", controller.CreateSchool).Methods("POST")
	router.HandleFunc("/{id}", controller.UpdateSchool).Methods("PUT")
	router.HandleFunc("/{id}", controller.DeleteSchool).Methods("DELETE")

	return router
}
