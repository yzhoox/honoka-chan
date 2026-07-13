package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"honoka-chan/config"
	"os"
)

func RSAGen(bits int) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	privateKeyFile, err := os.Create(config.PrivateKeyPath)
	if err != nil {
		panic(err)
	}
	defer privateKeyFile.Close()
	if err := pem.Encode(privateKeyFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	}); err != nil {
		panic(err)
	}

	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	publicKeyFile, err := os.Create(config.PublicKeyPath)
	if err != nil {
		panic(err)
	}
	defer publicKeyFile.Close()
	if err := pem.Encode(publicKeyFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	}); err != nil {
		panic(err)
	}
}

func RSAEncrypt(plainText []byte, publickeypath string) ([]byte, error) {
	data, err := os.ReadFile(publickeypath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid RSA public key PEM")
	}
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("PEM does not contain an RSA public key")
	}

	//nolint:staticcheck // Legacy clients require PKCS#1 v1.5 encryption.
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
}

func RSADecrypt(cipherText []byte) ([]byte, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return nil, err
	}

	//nolint:staticcheck // Legacy clients send PKCS#1 v1.5 encrypted payloads.
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
}

func RSASignSHA1(cipherText []byte) ([]byte, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return nil, err
	}

	msgHash := sha1.New()
	_, _ = msgHash.Write(cipherText)
	msgHashSum := msgHash.Sum(nil)

	return rsa.SignPSS(rand.Reader, privateKey, crypto.SHA1, msgHashSum, nil)
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid RSA private key PEM")
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("PEM does not contain an RSA private key")
	}
	return privateKey, nil
}
