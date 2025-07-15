package hash_util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Hash256String hashes a string using SHA-256 and returns the hexadecimal representation of the hash.
func Hash256String(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("input string cannot be empty")
	}

	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:]), nil
}

// Hash256Bytes hashes a byte slice using SHA-256 and returns the hexadecimal representation of the hash.
func Hash256Bytes(b []byte) (string, error) {
	if len(b) == 0 {
		return "", fmt.Errorf("input byte slice cannot be empty")
	}

	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}
