package v012

import (
	"encoding/base64"
)

func encodeName(name []byte) string {
	return base64.RawStdEncoding.EncodeToString(name)
}

func decodeName(name string) string {
	decoded, err := base64.RawStdEncoding.DecodeString(name)
	if err != nil {
		return ""
	}

	return string(decoded)
}
