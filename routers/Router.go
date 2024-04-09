package routers

import (
	controller "backend/controller"

	"github.com/gorilla/mux"
)

func Router(router *mux.Router) *mux.Router {

	router.HandleFunc("/schools", controller.GetAllSchools).Methods("GET")
	router.HandleFunc("/schools{id}", controller.GetSchoolByID).Methods("GET")
	router.HandleFunc("/schools", controller.CreateSchool).Methods("POST")
	router.HandleFunc("/schools/{id}", controller.UpdateSchool).Methods("PUT")
	router.HandleFunc("/schools/{id}", controller.DeleteSchool).Methods("DELETE")

	router.HandleFunc("/students", controller.GetAllStudents).Methods("GET")
	router.HandleFunc("/students/{id}", controller.GetStudentByID).Methods("GET")
	router.HandleFunc("/students", controller.CreateStudent).Methods("POST")
	router.HandleFunc("/students/{id}", controller.UpdateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", controller.DeleteStudent).Methods("DELETE")

	router.HandleFunc("/payments", controller.GetAllPayments).Methods("GET")
	router.HandleFunc("/payments/{id}", controller.GetPaymentByID).Methods("GET")
	// router.HandleFunc("/payments", controller.CreateCashPayment).Methods("POST")
	router.HandleFunc("/payments/{id}", controller.UpdatePayment).Methods("PUT")
	router.HandleFunc("/payments/{id}", controller.DeletePayment).Methods("DELETE")

	// router.HandleFunc("/fees", controller.GetAllFees).Methods("GET")
	// router.HandleFunc("/fees/{id}", controller.GetFeeByID).Methods("GET")
	// router.HandleFunc("/fees/", controller.CreateFee).Methods("POST")
	// router.HandleFunc("/fees/{id}", controller.UpdateFee).Methods("PUT")
	// router.HandleFunc("/fees/{id}", controller.DeleteFee).Methods("DELETE")

	return router
}
