package domain

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	UserID string `bson:"_id,omitempty" json:"userId"`
	Role   Role   `bson:"role" json:"role"`
}

var (
	RoleUnknown  = Role{"unknown"}
	RoleCustomer = Role{"customer"}
	RoleWorker   = Role{"worker"}
	RoleAdmin    = Role{"admin"}

	permissions = map[Role]int{
		RoleUnknown:  0,
		RoleCustomer: 1,
		RoleWorker:   2,
		RoleAdmin:    3,
	}
)

type Role struct {
	R string
}

func FromString(s string) Role {
	switch s {
	case RoleCustomer.R:
		return RoleCustomer
	case RoleWorker.R:
		return RoleWorker
	case RoleAdmin.R:
		return RoleAdmin
	default:
		return RoleUnknown
	}
}

func (r Role) String() string {
	return r.R
}

// CheckPermissions checks if compareWith Role has higher or equal permissions as current Role
func (r Role) CheckPermissions(required Role) bool {
	return permissions[r] >= permissions[required]
}

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Role string `json:"role"`
	}{
		Role: r.String(),
	})
}

func (r *Role) UnmarshalJSON(raw []byte) error {
	decoded := new(struct {
		Role string `json:"role"`
	})
	if err := json.Unmarshal(raw, decoded); err != nil {
		return err
	}
	r.R = decoded.Role
	return nil
}

func (r Role) MarshalBSON() ([]byte, error) {
	return bson.Marshal(struct {
		Role string `bson:"role" json:"role"`
	}{
		Role: r.String(),
	})
}

func (r *Role) UnmarshalBSON(raw []byte) error {
	decoded := new(struct {
		Role string `bson:"role" json:"role"`
	})
	if err := bson.Unmarshal(raw, decoded); err != nil {
		return err
	}
	r.R = decoded.Role
	return nil
}
