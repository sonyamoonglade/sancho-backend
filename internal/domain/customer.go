package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerExists   = errors.New("customer with such phone number exists")
)

type Customer struct {
	UserID          primitive.ObjectID   `json:"userId" bson:"_id,omitempty"`
	Role            Role                 `json:"role" bson:"role"`
	PhoneNumber     string               `json:"phoneNumber" bson:"phoneNumber"`
	Name            *string              `json:"name,omitempty" bson:"name,omitempty"`
	DeliveryAddress *UserDeliveryAddress `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	Session         *Session             `json:"session,omitempty" bson:"session,omitempty"`
}

type UserDeliveryAddress struct {
	Address   string `json:"address" bson:"address"`
	Entrance  int64  `json:"entrance" bson:"entrance"`
	Floor     int64  `json:"floor" bson:"floor"`
	Apartment int64  `json:"apartment" bson:"apartment"`
}
