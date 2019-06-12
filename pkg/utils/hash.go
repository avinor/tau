package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
