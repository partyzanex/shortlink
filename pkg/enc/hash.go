package enc

import (
	"crypto/rand"
	"encoding/base64"
)

func Hash(n int) (string, error) {
	r := make([]byte, n)

	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(r), nil
}
