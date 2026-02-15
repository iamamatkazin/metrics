// Package common предоставляет утилитарные функции для системы метрик.
// Включает функции криптографической подписи с использованием HMAC-SHA256.
package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// CalcSign вычисляет HMAC-SHA256 подпись для данных.
func CalcSign(key, src []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(src)

	return hex.EncodeToString(h.Sum(nil))
}

// PrintBuild выводит в консоль информацию о сборке.
func PrintBuild(version, date, commit string) {
	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}
