package v012

import (
	"encoding/base64"
)

// encodeName returns base64 encoding of string. Using raw base64 to avoid = at the end
func encodeName(name []byte) string {
	return base64.RawStdEncoding.EncodeToString(name)
}

// decodeName decodes from base64 to string.
func decodeName(name string) string {
	decoded, err := base64.RawStdEncoding.DecodeString(name)
	if err != nil {
		return ""
	}

	return string(decoded)
}
