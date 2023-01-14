package domain

type Customer struct {
	UserID          string           `json:"userId" bson:"userId,omitempty"`
	PhoneNumber     string           `json:"phoneNumber" bson:"phoneNumber"`
	DeliveryAddress *DeliveryAddress `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	Name            string           `json:"name" bson:"name"`
	Session         Session          `json:"session,omitempty" bson:"session,omitempty"`
}

type DeliveryAddress struct {
	Address   string `json:"address" bson:"address"`
	Entrance  int64  `json:"entrance" bson:"entrance"`
	Floor     int64  `json:"floor" bson:"floor"`
	Apartment int64  `json:"apartment" bson:"apartment"`
}
