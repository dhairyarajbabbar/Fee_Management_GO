// func CreatePartialPayment(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// // Parse request body to get school ID, student ID, and paid amount
	// var requestBody struct {
	// 	SchoolID   string  `json:"school_id"`
	// 	PaidAmount float64 `json:"paid_amount"`
	// }
	// decoder := json.NewDecoder(r.Body)
	// if err := decoder.Decode(&requestBody); err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	// // Convert school ID string to primitive.ObjectID
	// schoolID, err := primitive.ObjectIDFromHex(requestBody.SchoolID)
	// if err != nil {
	// 	http.Error(w, "Invalid school_id", http.StatusBadRequest)
	// 	return
	// }
	// // Get student ID from request parameters
	// params := mux.Vars(r)
	// studentID, err := primitive.ObjectIDFromHex(params["student_id"])
	// if err != nil {
	// 	http.Error(w, "Invalid student_id", http.StatusBadRequest)
	// 	return
	// }
	// // Retrieve student by ID
	// var student model.Student
	// filter := bson.M{"_id": studentID}
	// err = studentCollection.FindOne(context.Background(), filter).Decode(&student)
	// if err != nil {
	// 	http.Error(w, "Failed to retrieve student", http.StatusInternalServerError)
	// 	return
	// }
	// // Check if there are any previous partial payments
	// currentYear := time.Now().Year()
	// feeStatus := student.FeeStatuses[currentYear]
	// previousPartialAmount := 0.0
	// for _, status := range feeStatus.Statuses {
	// 	if strings.HasPrefix(status.Status, "partial_paid-") {
	// 		// Extract the paid amount from the status
	// 		amountStr := strings.TrimPrefix(status.Status, "partial_paid-")
	// 		amount, err := strconv.ParseFloat(amountStr, 64)
	// 		if err != nil {
	// 			http.Error(w, "Failed to parse partial paid amount", http.StatusInternalServerError)
	// 			return
	// 		}
	// 		previousPartialAmount += amount
	// 	}
	// }
	// // Calculate total paid amount including previous partial payments
	// totalPaidAmount := previousPartialAmount + requestBody.PaidAmount
	// // Mark the month as partially paid with the paid amount in the student's status
	// month := GetFirstUnpaidMonth(student) // Assuming partial payment goes to the first unpaid month
	// feeStatus.Statuses[month-1] = model.Status{
	// 	Status:    fmt.Sprintf("partial_paid-%.2f", totalPaidAmount),
	// 	PaymentID: primitive.NewObjectID(), // Generate new payment ID
	// }
	// // Save the updated student back to the database
	// update := bson.M{"$set": bson.M{"fee_statuses": student.FeeStatuses}}
	// _, err = studentCollection.UpdateOne(context.Background(), filter, update)
	// if err != nil {
	// 	http.Error(w, "Failed to update student's fee status", http.StatusInternalServerError)
	// 	return
	// }
	// // Mark the payment as PARTIAL PAYMENT
	// payment := model.Payment{
	// 	StudentID:   studentID,
	// 	SchoolID:    schoolID,
	// 	Amount:      requestBody.PaidAmount,
	// 	PaymentDate: time.Now(), // Assuming payment date is current date
	// 	Method:      "PARTIAL PAYMENT",
	// }
	// // Save the payment
	// _, err = paymentCollection.InsertOne(context.Background(), payment)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte(`{"message": "Partial payment recorded successfully"}`))
// }



//	func CreateCashPayment(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json")
//		// Parse request body to get school ID
//		var requestBody struct {
//			SchoolID string `json:"school_id"`
//		}
//		decoder := json.NewDecoder(r.Body)
//		if err := decoder.Decode(&requestBody); err != nil {
//			http.Error(w, "Invalid request body", http.StatusBadRequest)
//			return
//		}
//		// Convert school ID string to primitive.ObjectID
//		schoolID, err := primitive.ObjectIDFromHex(requestBody.SchoolID)
//		if err != nil {
//			http.Error(w, "Invalid school_id", http.StatusBadRequest)
//			return
//		}
//		// Get student ID from request parameters
//		params := mux.Vars(r)
//		studentID, err := primitive.ObjectIDFromHex(params["id"])
//		if err != nil {
//			http.Error(w, "Invalid student_id", http.StatusBadRequest)
//			return
//		}
//		// Retrieve student by ID
//		var student model.Student
//		filter := bson.M{"_id": studentID}
//		err = studentCollection.FindOne(context.Background(), filter).Decode(&student)
//		if err != nil {
//			http.Error(w, "Failed to retrieve student", http.StatusInternalServerError)
//			return
//		}
//		month, err := strconv.Atoi(r.FormValue("month"))
//		if err != nil {
//			month = GetFirstUnpaidMonth(student)
//		}
//		// Mark the month as paid in the student's status
//		currentYear := time.Now().Year()
//		feeStatus := student.FeeStatuses[currentYear]
//		feeStatus.Statuses[month-1] = model.Status{
//			Status:    "paid",
//			PaymentID: primitive.NewObjectID(), // Generate new payment ID
//		}
//		// Save the updated student back to the database
//		update := bson.M{"$set": bson.M{"fee_statuses": student.FeeStatuses}}
//		_, err = studentCollection.UpdateOne(context.Background(), filter, update)
//		if err != nil {
//			http.Error(w, "Failed to update student's fee status", http.StatusInternalServerError)
//			return
//		}
//		// fmt.Println(update)
//		payment := model.Payment{
//			StudentID:   studentID,
//			SchoolID:    schoolID,
//			Amount:      GetPaymentAmount(studentID, month),
//			PaymentDate: time.Now(),
//			Method:      "CASH",
//		}
//		// fmt.Println(payment)
//		_, err = paymentCollection.InsertOne(context.Background(), payment)
//		if err != nil {
//			log.Fatal(err)
//		}
//		w.WriteHeader(http.StatusCreated)
//		w.Write([]byte(`{"message": "Cash payment recorded successfully"}`))
//	}


func GetPaymentAmount(studentID primitive.ObjectID, month int) float64 {
	// // Assuming fee amount is stored in the first fee status for the current year
	// currentYear := time.Now().Year()
	// // Retrieve student from the database
	// var student model.Student
	// filter := bson.M{"_id": studentID}
	// err := studentCollection.FindOne(context.Background(), filter).Decode(&student)
	// if err != nil {
	// 	log.Println("Failed to retrieve student:", err)
	// 	return 0
	// }
	// feeStatus, ok := student.FeeStatuses[currentYear]
	// if !ok {
	// 	return 0
	// }
	// if month < 1 || month > 12 {
	// 	// Invalid month provided, return 0
	// }
	// // Get the fee amount from the first fee status for the current year
	// return feeStatus.FeeAmount
	return 0
}


