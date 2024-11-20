package hasher

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const SaltLen = 8

func HashPass(plainPassword string) ([]byte, error) {
	salt := make([]byte, SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return salt, fmt.Errorf("generate random bytes: %w", err)
	}
	return hash(salt, plainPassword)
}

func CheckPass(passHash []byte, plainPassword string) bool {
	salt := make([]byte, SaltLen)
	copy(salt, passHash[:SaltLen])

	userPassHash, err := hash(salt, plainPassword)
	if err != nil {
		return false
	}

	return bytes.Equal(userPassHash, passHash)
}

func hash(salt []byte, plainPassword string) ([]byte, error) {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)

	return append(salt, hashedPass...), nil
}
