package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	model "backend/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var studentCollection *mongo.Collection

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
	studentCollection = client.Database(dbName).Collection("student")
}

// MarkPastDueFees marks the fees with past-due generation dates as unpaid
func markPastDueFees(ctx context.Context, studentID primitive.ObjectID) error {
	var student model.Student
	if err := studentCollection.FindOne(ctx, bson.M{"_id": studentID}).Decode(&student); err != nil {
		return err
	}
	now := time.Now()
	updateRequired := false
	for year, fees := range student.FeeStatuses {
		for i, fee := range fees {
			fmt.Println("hello")
			if fee.GenerateDate.Before(now) && fee.Status != "paid" {
				fmt.Println("yes")
				student.FeeStatuses[year][i].Status = "unpaid"
				updateRequired = true
			}
		}
	}
	if updateRequired {
		_, err := studentCollection.UpdateOne(
			ctx,
			bson.M{"_id": studentID},
			bson.M{"$set": bson.M{"fee_statuses": student.FeeStatuses}},
		)
		fmt.Println("done")
		if err != nil {
			return err // Handle update error
		}
	}
	return nil
}

func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		SchoolID string `json:"school_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Convert school_id string to primitive.ObjectID
	schoolID, err := primitive.ObjectIDFromHex(requestBody.SchoolID)
	if err != nil {
		http.Error(w, "Invalid school_id", http.StatusBadRequest)
		return
	}
	// Retrieve all students belonging to the specified school
	cursor, err := studentCollection.Find(context.Background(), bson.M{"school_id": schoolID})
	if err != nil {
		http.Error(w, "Failed to retrieve students", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())
	// Iterate through the cursor and collect students
	var students []model.Student
	for cursor.Next(context.Background()) {
		var student model.Student
		if err := cursor.Decode(&student); err != nil {
			http.Error(w, "Failed to decode student", http.StatusInternalServerError)
			return
		}
		// Mark past due fees for the current student
		if err := markPastDueFees(context.Background(), student.ID); err != nil {
			http.Error(w, "Failed to mark past due fees", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}
	if err := cursor.Err(); err != nil {
		http.Error(w, "Failed to iterate students", http.StatusInternalServerError)
		return
	}
	// Encode and return the students as JSON
	if err := json.NewEncoder(w).Encode(students); err != nil {
		http.Error(w, "Failed to encode students", http.StatusInternalServerError)
		return
	}
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Extract studentID from URL parameters
	vars := mux.Vars(r)
	studentID := vars["id"]
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	// Update past due fees before fetching student data
	if err := markPastDueFees(r.Context(), objID); err != nil {
		log.Println(err) // Logging the error might be useful for debugging
		http.Error(w, "Failed to update fee status", http.StatusInternalServerError)
		return
	}
	// Fetch updated student data
	var student model.Student
	if err := studentCollection.FindOne(r.Context(), bson.M{"_id": objID}).Decode(&student); err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(student); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func getNextRollNumber(schoolID primitive.ObjectID) (string, error) { // Check if the school exists
	filter := bson.M{"_id": schoolID}
	var existingSchool model.School
	err := schoolCollection.FindOne(context.Background(), filter).Decode(&existingSchool)
	if err != nil {
		return "", err
	}
	nextRollNumber := len(existingSchool.StudentIDs) + 1
	return fmt.Sprintf("%d", nextRollNumber), nil
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type CreateStudentRequest struct {
		model.Student
		FeeAmount      float64   `json:"fee_amount"`
		EnrollmentDate time.Time `json:"enrollment_date"`
		SchoolID       string    `json:"school_id"`
	}
	// Parse the enrollment date from the JSON request
	var req CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	student := req.Student // Extract the student from the request
	schoolID, err := primitive.ObjectIDFromHex(req.SchoolID)
	if err != nil {
		http.Error(w, "Invalid school_id", http.StatusBadRequest)
		return
	}
	student.SchoolID = schoolID
	student.EnrollmentDate = req.EnrollmentDate
	rollNumber, err := getNextRollNumber(student.SchoolID)
	if err != nil {
		http.Error(w, "Failed to get roll number", http.StatusInternalServerError)
		return
	}
	student.RollNumber = rollNumber
	fmt.Println(req)
	// Set the EnrollmentDate to now if not provided
	if student.EnrollmentDate.IsZero() {
		student.EnrollmentDate = time.Now()
	}
	currentYear, startMonth, _ := student.EnrollmentDate.Date()
	sessionYearStart := currentYear
	sessionYearEnd := currentYear + 1
	if startMonth < time.April {
		sessionYearStart--
		sessionYearEnd--
	}
	session := fmt.Sprintf("%d-%d", sessionYearStart, sessionYearEnd)
	var fees []model.Fee
	// If the start month is after March, continue the loop for the rest of the year
	if startMonth > time.March {
		for m := startMonth; m <= time.December; m++ {
			currentMonth := m
			fee := model.Fee{
				// StudentID:    student.ID,
				Amount:       req.FeeAmount,
				Status:       "upcoming",
				PaidAmount:   0,
				GenerateDate: time.Date(sessionYearStart, currentMonth, 1, 0, 0, 0, 0, time.UTC),
				Payments:     []primitive.ObjectID{},
			}
			fees = append(fees, fee)
		}
		for m := time.January; m <= time.March; m++ {
			currentMonth := m
			fee := model.Fee{
				// StudentID:    student.ID,
				Amount:       req.FeeAmount,
				Status:       "upcoming",
				PaidAmount:   0,
				GenerateDate: time.Date(sessionYearEnd, currentMonth, 1, 0, 0, 0, 0, time.UTC),
				Payments:     []primitive.ObjectID{},
			}
			fees = append(fees, fee)
		}
	}
	// If the start month is before april, continue the loop for startMonth to March of the next year
	if startMonth <= time.March {
		for m := startMonth; m <= time.March; m++ {
			fee := model.Fee{
				// StudentID:    student.ID,
				Amount:       req.FeeAmount,
				Status:       "upcoming",
				PaidAmount:   0,
				GenerateDate: time.Date(sessionYearEnd, m, 1, 0, 0, 0, 0, time.UTC),
				Payments:     []primitive.ObjectID{},
			}
			fees = append(fees, fee)
		}
	}
	student.FeeStatuses = make(map[string][]model.Fee)
	student.FeeStatuses[session] = fees
	result, err := studentCollection.InsertOne(context.Background(), student)
	if err != nil {
		log.Printf("Could not create Student: %v", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}
	// for _, fee := range fees {
	// 	fee.StudentID = result.InsertedID.(primitive.ObjectID)
	// 	_, feeErr := feeCollection.InsertOne(context.Background(), fee)
	// 	if feeErr != nil {
	// 		log.Printf("Could not create Fee: %v", feeErr)
	// 	}
	// }
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.InsertedID)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var updatedStudent model.Student
	err := json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	studentID := vars["id"]
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	_, err = studentCollection.ReplaceOne(context.Background(), bson.M{"_id": objID}, updatedStudent)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	_, err = studentCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

// func CreateStudent(w http.ResponseWriter, r *http.Request) {
// 	var student model.Student
// 	err := json.NewDecoder(r.Body).Decode(&student)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	// Set default enrollment date if not provided
// if student.EnrollmentDate.IsZero() {
// 	student.EnrollmentDate = time.Now()
// }
// 	// Initialize fee statuses for the current year
// 	currentYear := time.Now().Year()
// 	student.FeeStatuses = make(map[int]model.FeeStatus)
// 	feeStatus := model.FeeStatus{
// 		FeeAmount: 0,                        // Set default fee amount
// 		Statuses:  make([]model.Status, 12), // Initialize statuses for 12 months
// 	}
// 	currentMonth := time.Now().Month()
// 	for i := 0; i < len(feeStatus.Statuses); i++ {
// 		if time.Month(i+1) < currentMonth {
// 			feeStatus.Statuses[i].Status = "none"
// 		} else if time.Month(i+1) == currentMonth {
// 			feeStatus.Statuses[i].Status = "pending"
// 		} else {
// 			feeStatus.Statuses[i].Status = "upcoming"
// 		}
// 	}
// 	student.FeeStatuses[currentYear] = feeStatus
// 	// Save student to database
// 	result, err := studentCollection.InsertOne(context.Background(), student)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//		// Respond with success
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusCreated)
//		json.NewEncoder(w).Encode(result.InsertedID)
//	}
