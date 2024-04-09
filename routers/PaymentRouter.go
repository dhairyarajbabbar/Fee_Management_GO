package routers

import (
	"backend/controller"

	"github.com/gorilla/mux"
)

func PaymentsRouter(router *mux.Router) *mux.Router {

	router.HandleFunc("/", controller.GetAllPayments).Methods("GET")
	router.HandleFunc("/{date}", controller.GetPaymentsForDate).Methods("GET")
	router.HandleFunc("/{id}", controller.GetPaymentByID).Methods("GET")
	// router.HandleFunc("/cash/{id}", controller.CreateCashPayment).Methods("POST")
	// router.HandleFunc("/cash/partial/{id}", controller.CreatePartialPayment).Methods("POST")
	router.HandleFunc("/{id}", controller.UpdatePayment).Methods("PUT")
	router.HandleFunc("/{id}", controller.DeletePayment).Methods("DELETE")

	return router
}
