package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Worker struct {
	UserID   primitive.ObjectID `json:"userId" bson:"_id"`
	Login    string             `json:"login" bson:"login"`
	Role     Role               `json:"role" bson:"role"`
	Password string             `json:"password" bson:"password"`
	Session  Session            `json:"session,omitempty" bson:"session,omitempty"`
	Name     string             `json:"name" bson:"name"`
}
