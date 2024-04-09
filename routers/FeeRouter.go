package routers

import (
	"backend/controller"

	"github.com/gorilla/mux"
)

func FeesRouter(router *mux.Router) *mux.Router {

	router.HandleFunc("/", controller.GetAllFees).Methods("GET")
	// router.HandleFunc("/{id}", controller.GetFeeByID).Methods("GET")
	// router.HandleFunc("/", controller.CreateFee).Methods("POST")
	// router.HandleFunc("/{id}", controller.UpdateFee).Methods("PUT")
	// router.HandleFunc("/{id}", controller.DeleteFee).Methods("DELETE")

	return router
}
