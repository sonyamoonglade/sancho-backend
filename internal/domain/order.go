package domain

import (
	"time"
)

type Order struct {
	OrderID           string                `json:"orderId" bson:"_id,omitempty"`
	NanoID            string                `json:"nanoId" bson:"nanoId"`
	CustomerID        string                `json:"customerId" bson:"customerId"`
	Cart              []Product             `json:"cart" bson:"cart"`
	Pay               Pay                   `json:"pay" bson:"pay"`
	Amount            int64                 `json:"amount" bson:"amount"`
	Discount          int64                 `json:"discount" bson:"discount"`
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
