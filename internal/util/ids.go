package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func RandomHex(n int) (string, error) {
	if n < 0 {
		return "", errors.New("random hex length must be non-negative")
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
