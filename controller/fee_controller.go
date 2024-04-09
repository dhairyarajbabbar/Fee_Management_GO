package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	model "backend/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var feeCollection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	feeCollection = client.Database(dbName).Collection("fee")
}

func GetAllFees(w http.ResponseWriter, r *http.Request) {
	var fees []model.Fee
	cursor, err := feeCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var fee model.Fee
		if err := cursor.Decode(&fee); err != nil {
			log.Fatal(err)
		}
		fees = append(fees, fee)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fees)
}

func CreateNewFee(w http.ResponseWriter, r *http.Request) {

}

// func GetFeeByID(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	feeID := vars["id"]
// 	objID, err := primitive.ObjectIDFromHex(feeID)
// 	if err != nil {
// 		http.Error(w, "Invalid fee ID", http.StatusBadRequest)
// 		return
// 	}
// 	var fee model.Fee
// 	err = feeCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&fee)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			http.Error(w, "Fee not found", http.StatusNotFound)
// 			return
// 		}
// 		log.Fatal(err)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(fee)
// }

// func CreateFee(w http.ResponseWriter, r *http.Request) {
// 	var fee model.Fee
// 	err := json.NewDecoder(r.Body).Decode(&fee)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	result, err := feeCollection.InsertOne(context.Background(), fee)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(result.InsertedID)
// }

// func UpdateFee(w http.ResponseWriter, r *http.Request) {
// 	var updatedFee model.Fee
// 	err := json.NewDecoder(r.Body).Decode(&updatedFee)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	vars := mux.Vars(r)
// 	feeID := vars["id"]
// 	objID, err := primitive.ObjectIDFromHex(feeID)
// 	if err != nil {
// 		http.Error(w, "Invalid fee ID", http.StatusBadRequest)
// 		return
// 	}
// 	_, err = feeCollection.ReplaceOne(context.Background(), bson.M{"_id": objID}, updatedFee)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	w.WriteHeader(http.StatusOK)
// }

// func DeleteFee(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	feeID := vars["id"]
// 	objID, err := primitive.ObjectIDFromHex(feeID)
// 	if err != nil {
// 		http.Error(w, "Invalid fee ID", http.StatusBadRequest)
// 		return
// 	}
// 	_, err = feeCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	w.WriteHeader(http.StatusOK)
// }
