package utils

import (
	"encoding/base64"
	"honoka-chan/pkg/db"
	"honoka-chan/pkg/encrypt"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IsSigned(unitId int) bool {
	exists, err := db.MainEng.Table("unit_sign_asset_m").Where("unit_id = ?", unitId).Exist()
	CheckErr(err)

	return exists
}

func GenXMS(resp []byte) string {
	return base64.StdEncoding.EncodeToString(encrypt.RSASignSHA1(resp))
}
