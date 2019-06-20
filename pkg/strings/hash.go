package strings

import (
	"crypto/sha1"
	"encoding/hex"
)

// HashFromBytes generates a hash code from byte array
func HashFromBytes(str []byte) string {
	h := sha1.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// Hash generates a hash code from string
func Hash(str string) string {
	return HashFromBytes([]byte(str))
}