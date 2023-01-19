package nanoid

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	random = "ABCDEFGHYJKLMNOPQRCTUVWXYZ0123456789"
	size   = 6
)

// GenerateNanoID will generate 6 char sequence from english alphabet and numbers
func GenerateNanoID() (string, error) {
	return gonanoid.Generate(random, size)
}
