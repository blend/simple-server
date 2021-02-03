package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
)

const nonceSize = 12

func loadSecret(secretPath string) ([]byte, error) {
	base64Secret, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return nil, err
	}
	base64Secret = bytes.TrimSpace(base64Secret)
	secret, err := base64.RawURLEncoding.DecodeString(string(base64Secret))
	if err != nil {
		return nil, err
	}

	// Validate length
	switch len(secret) {
	case 16, 24, 32:
		return secret, nil
	default:
		return nil, fmt.Errorf("Expect secret key length of 16, 24, or 32 bytes, got %d", len(secret))
	}
}

func encrypt(input string, key []byte) ([]byte, error) {
	// Based on https://golang.org/src/crypto/cipher/example_test.go
	plaintext := []byte(input)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

func decrypt(input string, key []byte) (string, error) {
	// Based on https://golang.org/src/crypto/cipher/example_test.go
	inputBytes, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	nonce, ciphertext := inputBytes[:nonceSize], inputBytes[nonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
