package refresh

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const refreshTokenMaxLength = 30

func GenerateRefreshToken() (string, error) {
	b := make([]byte, refreshTokenMaxLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("generate random string: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(b)

	return encoded, err
}
