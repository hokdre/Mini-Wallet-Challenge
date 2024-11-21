package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Encryption interface {
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
}

func NewAesEncryption(secret string) (*aesEncryption, error) {
	keyLength := len(secret)
	if keyLength != 16 && keyLength != 24 && keyLength != 32 {
		return nil, errors.New("invlid length of secret")
	}
	return &aesEncryption{secret: secret}, nil
}

type aesEncryption struct {
	secret string
}

func (a *aesEncryption) Encrypt(text string) (string, error) {
	keyBytes := []byte(a.secret)
	plainTextBytes := []byte(text)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	plainTextBytes = a.pkcs7Padding(plainTextBytes, aes.BlockSize)
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainTextBytes)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (a *aesEncryption) pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func (a *aesEncryption) Decrypt(text string) (string, error) {
	cipherTextBytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(a.secret))
	if err != nil {
		return "", err
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherTextBytes, cipherTextBytes)

	plainTextBytes, err := a.pkcs7UnPadding(cipherTextBytes)
	if err != nil {
		return "", err
	}

	return string(plainTextBytes), nil
}

func (a *aesEncryption) pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("invalid padding")
	}

	return data[:length-padding], nil
}
