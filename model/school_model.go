package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type School struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Name       string               `bson:"name"`
	Location   string               `bson:"location"`
	Contact    string               `bson:"contact"`
	StudentIDs []primitive.ObjectID `bson:"student_ids"`
}
