package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	StudentID   primitive.ObjectID `bson:"student_id"`
	SchoolID    primitive.ObjectID `bson:"school_id"`
	Amount      float64            `bson:"amount"`
	PaymentDate time.Time          `bson:"payment_date"`
	Method      string             `bson:"method"`
}
