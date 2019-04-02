package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// RememberTokenBytes is the set number of bytes used to generate random strings
const RememberTokenBytes = 32

// Bytes is a wrapper for crytpo/rand.read
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String creates a random string
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken creates a remember token
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
