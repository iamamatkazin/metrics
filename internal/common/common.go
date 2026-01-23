package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// CalcSign - вычисляем подпись.
func CalcSign(key, src []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(src)

	return hex.EncodeToString(h.Sum(nil))
}
