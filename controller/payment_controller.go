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

var paymentCollection *mongo.Collection

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
	paymentCollection = client.Database(dbName).Collection("payment")
}
func GetAllPayments(w http.ResponseWriter, r *http.Request) {
	var payments []model.Payment
	cursor, err := paymentCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var payment model.Payment
		if err := cursor.Decode(&payment); err != nil {
			log.Fatal(err)
		}
		payments = append(payments, payment)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}
func GetPaymentsForDate(w http.ResponseWriter, r *http.Request) {
}
func GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["id"]
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}
	var payment model.Payment
	err = paymentCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Payment not found", http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}
func GetFirstUnpaidMonth(student model.Student) int {
	// currentYear := time.Now().Year()
	// feeStatus, ok := student.FeeStatuses[currentYear]
	// if !ok { // If no fee status available for the current year, consider all months unpaid
	// 	return 1
	// }
	// for monthIndex, status := range feeStatus.Statuses {
	// 	// Check if the status is unpaid for this month (pending or upcoming)
	// 	if status.Status == "pending" || status.Status == "upcoming" {
	// 		// Return the first unpaid month index (monthIndex starts from 0, so add 1)
	// 		return monthIndex + 1
	// 	}
	// }
	return 0
}

func CreateCashPaymentByAmount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		StudentID string  `json:"student_id"`
		SchoolID  string  `json:"school_id"`
		Amount    float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	studentID, err := primitive.ObjectIDFromHex(req.StudentID)
	if err != nil {
		http.Error(w, "Invalid student_id", http.StatusBadRequest)
		return
	}

	schoolID, err := primitive.ObjectIDFromHex(req.SchoolID)
	if err != nil {
		http.Error(w, "Invalid school_id", http.StatusBadRequest)
		return
	}

	var student model.Student
	if err := studentCollection.FindOne(context.Background(), bson.M{"_id": studentID}).Decode(&student); err != nil {
		http.Error(w, "Failed to find student", http.StatusInternalServerError)
		return
	}

	currentYear := fmt.Sprintf("%d", time.Now().Year())
	fees, exists := student.FeeStatuses[currentYear]
	if !exists {
		http.Error(w, "No fees found for the current year", http.StatusNotFound)
		return
	}

	var firstUnpaidOrPartial *model.Fee
	firstUnpaidOrPartialIndex := -1
	for i, fee := range fees {
		if fee.Status == "unpaid" || fee.Status == "partially_paid" {
			firstUnpaidOrPartial = &fees[i]
			firstUnpaidOrPartialIndex = i
			break
		}
	}

	if firstUnpaidOrPartial == nil {
		http.Error(w, "All fees are already paid for the current year", http.StatusNotFound)
		return
	}

	remainingAmount := req.Amount
	for i := firstUnpaidOrPartialIndex; i < len(fees) && remainingAmount > 0; i++ {
		initialPaidAmount := fees[i].PaidAmount
		fees[i].PaidAmount += remainingAmount

		if fees[i].PaidAmount >= fees[i].Amount {
			fees[i].Status = "paid"
			remainingAmount -= fees[i].Amount - initialPaidAmount
		} else {
			fees[i].Status = "partially_paid"
			break
		}
	}

	// Create Payment object
	payment := model.Payment{
		StudentID:   studentID,
		SchoolID:    schoolID,
		Amount:      req.Amount,
		PaymentDate: time.Now(),
		Method:      "cash",
	}

	// Save the payment to the database
	if _, err := paymentCollection.InsertOne(context.Background(), payment); err != nil {
		http.Error(w, "Failed to save payment", http.StatusInternalServerError)
		return
	}

	// Update student's fee statuses
	_, err = studentCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": studentID},
		bson.M{"$set": bson.M{"fee_statuses": student.FeeStatuses}},
	)
	if err != nil {
		http.Error(w, "Failed to update student's fee statuses", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	var updatedPayment model.Payment
	err := json.NewDecoder(r.Body).Decode(&updatedPayment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	paymentID := vars["id"]
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}
	_, err = paymentCollection.ReplaceOne(context.Background(), bson.M{"_id": objID}, updatedPayment)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

func DeletePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["id"]
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}
	_, err = paymentCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

// determineSessionYear calculates the session year based on the current date
func determineSessionYear(currentTime time.Time) string {
	year := currentTime.Year()
	if currentTime.Month() < time.April {
		year-- // If before April, it's the previous session year
	}
	return fmt.Sprintf("%d-%d", year, year+1)
}

// initializeFeesForSession creates an array of fees for the current session starting from the current month up to March of the next year
func initializeFeesForSession(startTime time.Time) []model.Fee {
	fees := make([]model.Fee, 0)
	// Determine the start month and year for the fee generation
	startMonth := startTime.Month()
	startYear := startTime.Year()
	if startMonth < time.April {
		startMonth = time.April
	} else if startMonth > time.March {
		startYear++
		startMonth = time.April
	}
	// Generate fees for each month from the start month to March of the next year
	for month := startMonth; month <= time.March || (startYear == startTime.Year() && month <= startTime.Month()); month++ {
		fee := model.Fee{
			GenerateDate: startTime,
			Status:       "upcoming",
			PaidAmount:   0,
		}
		fees = append(fees, fee)
		// Increment the month for the next fee
		nextMonthTime := startTime.AddDate(0, 1, 0)
		startTime = nextMonthTime
		if startTime.Month() == time.April {
			break
		}
	}
	return fees
}
