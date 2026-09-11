package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const KeyPrefix = "dgr_"

func HashKey(rawKey string) string {
	sum := sha256.Sum256([]byte(rawKey))

	return fmt.Sprintf("%x", sum)
}

func Prefix(rawKey string) string {
	if len(rawKey) < 8 {
		return rawKey
	}
	return rawKey[:8]
}

func GenerateKey() (rawKey, hash, prefix string, err error) {
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return "", "", "", fmt.Errorf("generating key: %w", err)
	}

	rawKey = KeyPrefix + base64.RawURLEncoding.EncodeToString(buf)
	hash = HashKey(rawKey)
	prefix = Prefix(rawKey)

	return rawKey, hash, prefix, nil
}
