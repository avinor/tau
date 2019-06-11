package utils

import (
	"encoding/hex"
	"crypto/sha1"
)

func Hash(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
