package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminAlreadyExists = errors.New("admin already exists")
)

type Admin struct {
	UserID   primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	Login    string             `json:"login" bson:"login"`
	Role     Role               `json:"role" bson:"role"`
	Password string             `json:"password" bson:"password"`
	Session  Session            `json:"session,omitempty" bson:"session,omitempty"`
}
