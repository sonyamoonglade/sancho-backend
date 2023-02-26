package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Worker struct {
	UserID   primitive.ObjectID `json:"userId" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Login    string             `json:"login" bson:"login"`
	Password string             `json:"password" bson:"password"`
	Role     Role               `json:"role" bson:"role"`
	Session  Session            `json:"session,omitempty" bson:"session,omitempty"`
}
