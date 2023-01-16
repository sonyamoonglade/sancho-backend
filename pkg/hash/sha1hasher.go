package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

type SHA1Hasher struct{}

func NewSHA1Hasher() *SHA1Hasher {
	return &SHA1Hasher{}
}

func (s SHA1Hasher) Hash(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
