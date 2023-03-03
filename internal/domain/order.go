package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	invalidEntrance     = "invalid entrance"
	invalidFloor        = "invalid floor"
	invalidApartment    = "invalid apartment"
	invalidDeliveryTime = "invalid delivery time"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrHavePendingOrder = errors.New("have pending order")
)

type Order struct {
	OrderID           primitive.ObjectID    `json:"orderId" bson:"_id,omitempty"`
	NanoID            string                `json:"nanoId" bson:"nanoId"`
	CustomerID        string                `json:"customerId" bson:"customerId"`
	Cart              []CartProduct         `json:"cart" bson:"cart"`
	Pay               Pay                   `json:"pay" bson:"pay"`
	Amount            int64                 `json:"amount" bson:"amount"`
	Discount          float64               `json:"discount" bson:"discount"`
	DiscountedAmount  int64                 `json:"discountedAmount" bson:"discountedAmount"`
	Status            OrderStatus           `json:"status" bson:"status"`
	IsDelivered       bool                  `json:"isDelivered" bson:"isDelivered"`
	DeliveryAddress   *OrderDeliveryAddress `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	CreatedAt         time.Time             `json:"createdAt" bson:"createdAt"`
	VerifiedAt        *time.Time            `json:"verifiedAt,omitempty" bson:"verifiedAt,omitempty"`
	CompletedAt       *time.Time            `json:"completedAt,omitempty" bson:"completedAt,omitempty"`
	CancelledAt       *time.Time            `json:"cancelledAt,omitempty" bson:"cancelledAt,omitempty"`
	CancelExplanation *string               `json:"cancelExplanation,omitempty" bson:"cancelExplanation,omitempty"`
}

type OrderDeliveryAddress struct {
	IsAsap      bool      `json:"isAsap" bson:"isAsap"`
	Address     string    `json:"address" bson:"address"`
	Entrance    int64     `json:"entrance" bson:"entrance"`
	Floor       int64     `json:"floor" bson:"floor"`
	Apartment   int64     `json:"apartment" bson:"apartment"`
	DeliveredAt time.Time `json:"deliveredAt" bson:"deliveredAt"`
}

func (o *OrderDeliveryAddress) ToUserDeliveryAddress() *UserDeliveryAddress {
	return &UserDeliveryAddress{
		Address:   o.Address,
		Entrance:  o.Entrance,
		Floor:     o.Floor,
		Apartment: o.Apartment,
	}
}

func (o *OrderDeliveryAddress) IsValid() (bool, string) {
	if o.Entrance <= 0 {
		return false, invalidEntrance
	}
	if o.Apartment <= 0 {
		return false, invalidApartment
	}
	if o.Floor <= 0 {
		return false, invalidFloor
	}
	if o.DeliveredAt.Before(time.Now().UTC()) {
		return false, invalidDeliveryTime
	}
	return true, ""
}
