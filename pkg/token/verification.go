package token

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateVerificationToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}