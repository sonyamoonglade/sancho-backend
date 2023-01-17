package domain

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	StatusWaitingForVerification = OrderStatus{"waiting for verification"}
	StatusVerified               = OrderStatus{"verified"}
	StatusCompleted              = OrderStatus{"completed"}
	StatusCancelled              = OrderStatus{"cancelled"}
)

type OrderStatus struct {
	V string
}

func (o OrderStatus) String() string {
	return o.V
}

func (o OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status string `json:"status"`
	}{
		Status: o.String(),
	})
}

func (o *OrderStatus) UnmarshalJSON(raw []byte) error {
	decoded := new(struct {
		Status string `json:"status"`
	})
	if err := json.Unmarshal(raw, decoded); err != nil {
		return err
	}
	o.V = decoded.Status
	return nil
}

func (o OrderStatus) MarshalBSON() ([]byte, error) {
	return bson.Marshal(struct {
		Status string `bson:"status"`
	}{
		Status: o.String(),
	})
}

func (o *OrderStatus) UnmarshalBSON(raw []byte) error {
	decoded := new(struct {
		Status string `bson:"status"`
	})
	if err := bson.Unmarshal(raw, decoded); err != nil {
		return err
	}
	o.V = decoded.Status
	return nil
}
