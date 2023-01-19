package domain

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	PayOnPickup = Pay{"on pickup"}
	PayOnline   = Pay{"online"}
)

type Pay struct {
	P string
}

func (p Pay) String() string {
	return p.P
}
func (p Pay) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Pay string `json:"pay"`
	}{
		Pay: p.String(),
	})
}

func (p *Pay) UnmarshalJSON(raw []byte) error {
	decoded := new(struct {
		Pay string `json:"pay"`
	})
	if err := json.Unmarshal(raw, decoded); err != nil {
		return err
	}
	p.P = decoded.Pay
	return nil
}

func (p Pay) MarshalBSON() ([]byte, error) {
	return bson.Marshal(struct {
		Pay string `bson:"pay"`
	}{
		Pay: p.String(),
	})
}

func (p *Pay) UnmarshalBSON(raw []byte) error {
	decoded := new(struct {
		Pay string `bson:"pay"`
	})
	if err := bson.Unmarshal(raw, decoded); err != nil {
		return err
	}
	p.P = decoded.Pay
	return nil
}
