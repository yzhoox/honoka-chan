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
	"path/filepath"
	"runtime"
)

func RSAGen(bits int) {
	//get current path
	_, currentpath, _, _ := runtime.Caller(0)
	currentpath = filepath.Dir(currentpath)

	//----------------------------------------------private key

	// GenerateKey generates an RSA keypair of the given bit size using the
	// random source random (for example, crypto/rand.Reader).
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}

	//serialize privatekey to ASN.1 der by x509.MarshalPKCS8PrivateKey
	x509privatekey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	//encode x509 to pem and save to file
	//1. create privatefile
	privatekeyfile, err := os.Create(config.PrivateKeyPath)
	if err != nil {
		panic(err)
	}
	defer privatekeyfile.Close()
	//2. new a pem block struct object
	privatekeyblock := pem.Block{
		Type:    "PRIVATE KEY",
		Headers: nil,
		Bytes:   x509privatekey,
	}
	//3. save to file
	pem.Encode(privatekeyfile, &privatekeyblock)

	//----------------------------------------------public key

	//get public key
	publickey := privateKey.PublicKey
	//serialize publickey to ASN.1 der by x509.MarshalPKCS8PublicKey
	x509publickey, _ := x509.MarshalPKIXPublicKey(&publickey)

	//encode x509 to pem and save to file
	//1. create publickeyfile
	publickeyfile, err := os.Create(config.PublicKeyPath)
	if err != nil {
		panic(err)
	}
	defer publickeyfile.Close()

	//2. new a pem block struct object
	publickeyblock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   x509publickey,
	}

	//3. save to file
	pem.Encode(publickeyfile, &publickeyblock)
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

	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
}

func RSADecrypt(cipherText []byte) ([]byte, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return nil, err
	}

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
