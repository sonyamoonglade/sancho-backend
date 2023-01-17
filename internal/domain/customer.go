package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	UserID          primitive.ObjectID   `json:"userId" bson:"userId,omitempty"`
	PhoneNumber     string               `json:"phoneNumber" bson:"phoneNumber"`
	DeliveryAddress *UserDeliveryAddress `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	Role            Role                 `json:"role" bson:"role"`
	Name            string               `json:"name" bson:"name"`
	Session         Session              `json:"session,omitempty" bson:"session,omitempty"`
}

type UserDeliveryAddress struct {
	Address   string `json:"address" bson:"address"`
	Entrance  int64  `json:"entrance" bson:"entrance"`
	Floor     int64  `json:"floor" bson:"floor"`
	Apartment int64  `json:"apartment" bson:"apartment"`
}
