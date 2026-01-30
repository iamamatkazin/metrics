// Package common предоставляет утилитарные функции для системы метрик.
// Включает функции криптографической подписи с использованием HMAC-SHA256.
package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// CalcSign вычисляет HMAC-SHA256 подпись для данных.
func CalcSign(key, src []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(src)

	return hex.EncodeToString(h.Sum(nil))
}
