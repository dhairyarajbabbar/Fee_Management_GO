package main

import (
	"fmt"
	"log"
	"net/http"

	routers "backend/routers"

	"github.com/gorilla/mux"
)

func main() {
	// Set up MongoDB connection
	// client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://dhairyarajbabbar:qvo0dkslzZ3UNbYc@feemanagement.lmkfmxp.mongodb.net/?retryWrites=true&w=majority"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Disconnect(ctx)
	// database := client.Database("your_database_name")
	// schoolsCollection := database.Collection("schools")
	// studentsCollection := database.Collection("students")
	// feesCollection := database.Collection("fees")
	// paymentsCollection := database.Collection("payments")

	// // Set up routes
	// router := mux.NewRouter()

	// Define HTTP handler functions for your routes
	// router.HandleFunc("/schools", func(w http.ResponseWriter, r *http.Request) {
	// 	var schools []School
	// 	cursor, err := schoolsCollection.Find(ctx, bson.M{})
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	defer cursor.Close(ctx)
	// 	if err := cursor.All(ctx, &schools); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	json.NewEncoder(w).Encode(schools)
	// }).Methods("GET")

	// serverAddr := ":8080"
	// server := &http.Server{
	// 	Addr:         serverAddr,
	// 	Handler:      router,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	// log.Printf("Server listening on %s", serverAddr)
	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("MongoDB API")

	// r :=
	// http.Handle("/schools", router.SchoolsRouter())
	// // http.HandleFunc("/schools", router.SchoolsRouter().ServeHTTP)
	// // sr :=
	// http.Handle("/students", router.StudentsRouter())

	// // pr :=
	// http.Handle("/payments", router.PaymentsRouter())

	// // fr :=
	// http.Handle("/fees", router.FeesRouter())

	// http.HandleFunc("/schools", func(w http.ResponseWriter, r *http.Request) {
	// 	router.SchoolsRouter().ServeHTTP(w, r)
	// })

	// http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
	// 	router.StudentsRouter().ServeHTTP(w, r)
	// })

	// http.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
	// 	router.PaymentsRouter().ServeHTTP(w, r)
	// })

	// http.HandleFunc("/fees", func(w http.ResponseWriter, r *http.Request) {
	// 	router.FeesRouter().ServeHTTP(w, r)
	// })
	// r := router.Router()
	router := mux.NewRouter()
	// routers.SchoolsRouter(router)
	// routers.FeesRouter(router)
	// routers.PaymentsRouter(router)

	studentrouter := router.PathPrefix("/student").Subrouter()
	routers.StudentsRouter(studentrouter)

	schoolrouter := router.PathPrefix("/school").Subrouter()
	routers.SchoolsRouter(schoolrouter)

	feesrouter := router.PathPrefix("/fee").Subrouter()
	routers.FeesRouter(feesrouter)

	paymentsrouter := router.PathPrefix("/payment").Subrouter()
	routers.PaymentsRouter(paymentsrouter)

	fmt.Println("Listening at port 4000 ...")
	log.Fatal(http.ListenAndServe("localhost:4000", router))

	fmt.Println("Listening at port 4000 ...")

}
