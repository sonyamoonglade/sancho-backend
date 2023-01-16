package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminAlreadyExists = errors.New("admin already exists")
)

type Admin struct {
	UserID   string  `json:"userId" bson:"userId,omitempty"`
	Login    string  `json:"login" bson:"login"`
	Role     Role    `json:"role" bson:"role"`
	Password string  `json:"password" bson:"password"`
	Session  Session `json:"session,omitempty" bson:"session,omitempty"`
}
