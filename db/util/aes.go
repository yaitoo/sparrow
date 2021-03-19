package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

var (
	ErrSecretTooShort = errors.New("secret is too short, length must be 16/24/32")
	ErrCipherTooShort = errors.New("ciphertext is too short")
)

func AesEncrypt(plaintext, secret string) (string, error) {
	secretLen := len(secret)
	if secretLen == 16 || secretLen == 24 || secretLen == 32 {
		block, err := aes.NewCipher([]byte(secret))
		if err != nil {
			return "", nil
		}
		cipherText := make([]byte, aes.BlockSize+len(plaintext))
		iv := cipherText[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", nil
		}
		cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:],
			[]byte(plaintext))
		return hex.EncodeToString(cipherText), nil
	}
	return "", ErrSecretTooShort
}

func AesDecrypt(cipherText, secret string) (string, error) {
	secretLen := len(secret)
	if secretLen == 16 || secretLen == 24 || secretLen == 32 {
		cipherBytes, err := hex.DecodeString(cipherText)
		if err != nil {
			return "", err
		}
		block, err := aes.NewCipher([]byte(secret))
		if err != nil {
			return "", err
		}
		if len(cipherBytes) < aes.BlockSize {
			return "", ErrCipherTooShort
		}
		iv := cipherBytes[:aes.BlockSize]
		cipherBytes = cipherBytes[aes.BlockSize:]
		cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherBytes, cipherBytes)
		return string(cipherBytes), nil
	}
	return "", ErrSecretTooShort
}
