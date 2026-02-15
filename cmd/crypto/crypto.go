// Package crypto - пекет в котором реализованы методы кодирования и
// декодиврования сообщений. Методы Encrypt и Decrypt принимают на вход
// местоположение открытого и закрытого ключей.
package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
)

// Encrypt - кодирует сообщение.
func Encrypt(path string, data any) ([]byte, error) {
	byteData, err := json.Marshal(data)
	if err != nil {
		return byteData, err
	}

	if path == "" {
		return byteData, nil
	}

	publicKey, err := getPublicKey(path)
	if err != nil {
		return nil, err
	}

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, byteData)
	if err != nil {
		return nil, err
	}

	return cipher, nil
}

// Decrypt - декодирует сообщение.
func Decrypt(path string, data []byte) ([]byte, error) {
	if path == "" {
		return data, nil
	}

	privateKey, err := getPrivateKey(path)
	if err != nil {
		return nil, err
	}

	res, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getPrivateKey(path string) (*rsa.PrivateKey, error) {
	pemBlock, err := getPemBlock(path)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func getPublicKey(path string) (*rsa.PublicKey, error) {
	pemBlock, err := getPemBlock(path)
	if err != nil {
		return nil, err
	}

	publicKey, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func getPemBlock(path string) (*pem.Block, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(keyBytes)

	return pemBlock, nil
}
