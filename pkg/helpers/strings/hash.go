package strings

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/avinor/tau/pkg/helpers/ui"
)

// HashFromBytes generates a hash code from byte array
func HashFromBytes(str []byte) string {
	h := sha1.New()
	if _, err := h.Write(str); err != nil {
		ui.Fatal("Failed creating hash from string: %s", err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Hash generates a hash code from string
func Hash(str string) string {
	return HashFromBytes([]byte(str))
}
