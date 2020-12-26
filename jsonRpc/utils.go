package jsonRpc

import (
	"encoding/base64"
	"encoding/hex"
)

func base64ToHex(str string) string {
	b, _ := base64.RawURLEncoding.DecodeString(str)
	return hex.EncodeToString(b)
}

func hexToBase64(str string) string {
	b, _ := hex.DecodeString(str)
	return base64.RawURLEncoding.EncodeToString(b)
}
