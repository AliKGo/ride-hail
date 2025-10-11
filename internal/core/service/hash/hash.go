package hash

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	saltSize  = 16
	iterCount = 100000
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password empty")
	}

	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	h := sha512.New()
	h.Write([]byte(password))
	h.Write(salt)
	sum := h.Sum(nil)

	for i := 0; i < iterCount; i++ {
		h = sha512.New()
		h.Write(sum)
		h.Write(salt)
		sum = h.Sum(nil)
	}

	saltB64 := base64.StdEncoding.EncodeToString(salt)
	hashB64 := base64.StdEncoding.EncodeToString(sum)

	return fmt.Sprintf("%s$%s", saltB64, hashB64), nil
}

func VerifyPassword(stored, password string) (bool, error) {
	parts := strings.SplitN(stored, "$", 2)
	if len(parts) != 2 {
		return false, errors.New("invalid stored hash format")
	}
	saltB64 := parts[0]
	hashB64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, err
	}
	expected, err := base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, err
	}

	h := sha512.New()
	h.Write([]byte(password))
	h.Write(salt)
	sum := h.Sum(nil)

	for i := 0; i < iterCount; i++ {
		h = sha512.New()
		h.Write(sum)
		h.Write(salt)
		sum = h.Sum(nil)
	}

	if len(sum) != len(expected) {
		return false, nil
	}
	if subtle.ConstantTimeCompare(sum, expected) == 1 {
		return true, nil
	}
	return false, nil
}
