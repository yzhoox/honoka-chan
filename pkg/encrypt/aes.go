package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

var (
	iv = []byte("12345678abcdefgh")
)

func padding(plainText []byte, blockSize int) []byte {
	n := blockSize - len(plainText)%blockSize
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

func unPadding(cipherText []byte) ([]byte, error) {
	if len(cipherText) == 0 {
		return nil, fmt.Errorf("AES plaintext is empty")
	}

	end := cipherText[len(cipherText)-1]
	if end == 0 || int(end) > len(cipherText) {
		return nil, fmt.Errorf("invalid AES padding")
	}
	for _, value := range cipherText[len(cipherText)-int(end):] {
		if value != end {
			return nil, fmt.Errorf("invalid AES padding")
		}
	}
	return cipherText[:len(cipherText)-int(end)], nil
}

func AESCBCEncrypt(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plainText = padding(plainText, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	return cipherText, nil
}

func AESCBCDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("AES ciphertext must be a non-empty multiple of the block size")
	}
	plainText := make([]byte, len(cipherText))
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(plainText, cipherText)
	return unPadding(plainText)
}
