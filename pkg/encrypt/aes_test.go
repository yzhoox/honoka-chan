package encrypt

import (
	"bytes"
	"testing"
)

func TestAESCBCDecryptRoundTrip(t *testing.T) {
	key := []byte("1234567890abcdef")
	plainText := []byte("login payload")

	cipherText, err := AESCBCEncrypt(plainText, key)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	decrypted, err := AESCBCDecrypt(cipherText, key)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(decrypted, plainText) {
		t.Fatalf("decrypted plaintext = %q, want %q", decrypted, plainText)
	}
}

func TestAESCBCDecryptRejectsMalformedCiphertext(t *testing.T) {
	key := []byte("1234567890abcdef")
	for _, cipherText := range [][]byte{nil, make([]byte, 15)} {
		if _, err := AESCBCDecrypt(cipherText, key); err == nil {
			t.Fatalf("AESCBCDecrypt(%d bytes) succeeded", len(cipherText))
		}
	}

	if _, err := unPadding([]byte{1, 2, 3, 2}); err == nil {
		t.Fatal("unPadding accepted invalid padding")
	}
}
