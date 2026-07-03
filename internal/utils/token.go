package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

const refreshTokenSize = 32

// GenerateRefreshToken creates a cryptographically secure random token.
func GenerateRefreshToken() (string, error) {
	token := make([]byte, refreshTokenSize)

	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}

// HashToken hashes a token before storing it in the database.
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
