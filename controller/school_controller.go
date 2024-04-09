package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "backend/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://dhairyarajbabbar:qvo0dkslzZ3UNbYc@feemanagement.lmkfmxp.mongodb.net/?retryWrites=true&w=majority"
const dbName = "FeeManagement"

var schoolCollection *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")
	schoolCollection = client.Database(dbName).Collection("school")
	fmt.Println("Collection instance is ready")
}

func GetAllSchools(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var schools []model.School
	filter := bson.M{}
	cur, err := schoolCollection.Find(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to retrieve schools"}`))
		return
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var school model.School
		if err := cur.Decode(&school); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to decode schools"}`))
			return
		}
		schools = append(schools, school)
	}
	if len(schools) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "No schools found"}`))
		return
	}
	schoolsJSON, err := json.Marshal(schools)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal schools to JSON"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(schoolsJSON)
}
func GetSchoolByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	schoolID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid school ID"}`))
		return
	}
	var school model.School
	filter := bson.M{"_id": schoolID}
	err = schoolCollection.FindOne(context.Background(), filter).Decode(&school)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "School not found"}`))
		return
	}
	schoolJSON, err := json.Marshal(school)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal school to JSON"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(schoolJSON)
}

func CreateSchool(w http.ResponseWriter, r *http.Request) {
	var school model.School
	err := json.NewDecoder(r.Body).Decode(&school)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request body"}`))
		return
	}
	result, err := schoolCollection.InsertOne(context.Background(), school)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to create school"}`))
		return

	}
	createdID := result.InsertedID.(primitive.ObjectID).Hex()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": createdID})
}

func UpdateSchool(w http.ResponseWriter, r *http.Request) {
	var updatedSchool model.School
	err := json.NewDecoder(r.Body).Decode(&updatedSchool)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	schoolID := vars["id"]
	objID, err := primitive.ObjectIDFromHex(schoolID)
	if err != nil {
		http.Error(w, "Invalid school ID", http.StatusBadRequest)
		return
	}
	result, err := schoolCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": updatedSchool})
	if err != nil {
		log.Fatal(err)
	}
	if result.ModifiedCount == 0 {
		http.Error(w, "School not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteSchool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schoolID := vars["id"]
	objID, err := primitive.ObjectIDFromHex(schoolID)
	if err != nil {
		http.Error(w, "Invalid school ID", http.StatusBadRequest)
		return
	}
	result, err := schoolCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		log.Fatal(err)
	}
	if result.DeletedCount == 0 {
		http.Error(w, "School not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
