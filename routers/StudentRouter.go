package routers

import (
	"backend/controller"

	"github.com/gorilla/mux"
)

func StudentsRouter(router *mux.Router) *mux.Router {

	router.HandleFunc("/", controller.GetAllStudents).Methods("GET")
	router.HandleFunc("/{id}", controller.GetStudentByID).Methods("GET")
	router.HandleFunc("/", controller.CreateStudent).Methods("POST")
	router.HandleFunc("/{id}", controller.UpdateStudent).Methods("PUT")
	router.HandleFunc("/{id}", controller.DeleteStudent).Methods("DELETE")

	return router
}
