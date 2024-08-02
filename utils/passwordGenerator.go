package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomPassword generates a random password with the specified length
func GenerateRandomPassword() (string, error) {
	length := 12
	numBytes := length / 4 * 3
	if length%4 != 0 {
		numBytes += 3
	}

	randomBytes := make([]byte, numBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	password := base64.URLEncoding.EncodeToString(randomBytes)

	password = password[:length]

	return password, nil
}
