package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fee struct {
	Amount       float64              `bson:"amount" json:"amount"`
	GenerateDate time.Time            `bson:"generate_date" json:"generate_date"`
	Status       string               `bson:"status" json:"status"`           // paid/unpaid/partially_paid/upcoming
	PaidAmount   float64              `bson:"paid_amount" json:"paid_amount"` //0 default
	Payments     []primitive.ObjectID `bson:"payments" json:"payments"`
}

type Student struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	RollNumber     string             `bson:"roll_number"`
	Password       string             `bson:"password"`
	SchoolID       primitive.ObjectID `bson:"school_id"`
	Contact        string             `bson:"contact"`
	EnrollmentDate time.Time          `bson:"enrollment_date"`
	FeeStatuses    map[string][]Fee   `bson:"fee_statuses"`
}
