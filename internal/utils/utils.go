package utils

import (
	"encoding/base64"
	"honoka-chan/pkg/encrypt"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GenXMS(resp []byte) string {
	return base64.StdEncoding.EncodeToString(encrypt.RSASignSHA1(resp))
}
