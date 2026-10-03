package password

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const (
	algorithm  = "pbkdf2_sha256"
	iterations = 210000
	saltSize   = 16
	keySize    = 32
)

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(password string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("can't generate salt: %w", err)
	}

	key, err := pbkdf2.Key(sha256.New, password, salt, iterations, keySize)
	if err != nil {
		return "", fmt.Errorf("can't derive key: %w", err)
	}

	return strings.Join([]string{
		algorithm,
		strconv.Itoa(iterations),
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	}, "$"), nil
}

func (h *Hasher) Compare(hash string, password string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 4 || parts[0] != algorithm {
		return false
	}

	hashIterations, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	expectedKey, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}

	actualKey, err := pbkdf2.Key(sha256.New, password, salt, hashIterations, len(expectedKey))
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(actualKey, expectedKey) == 1
}
