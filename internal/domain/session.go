package domain

import "time"

type Session struct {
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt" bson:"expiresAt"`
}

func NewExpiresAt(ttl time.Duration) time.Time {
	return time.Now().UTC().Add(ttl)
}
