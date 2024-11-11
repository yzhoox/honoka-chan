package handler

import (
	"encoding/base64"
	"honoka-chan/config"
	"honoka-chan/encrypt"

	"xorm.io/xorm"
)

var (
	SifCdnServer string
	AsCdnServer  string
	ErrorMsg     = `{"code":20001,"message":""}`
	MainEng      *xorm.Engine
	UserEng      *xorm.Engine
)

func init() {
	SifCdnServer = config.Conf.Settings.SifCdnServer

	MainEng = config.MainEng
	UserEng = config.UserEng
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IsSigned(unitId int) bool {
	exists, err := MainEng.Table("unit_sign_asset_m").Where("unit_id = ?", unitId).Exist()
	CheckErr(err)

	return exists
}

func GenXMS(resp []byte) string {
	return base64.StdEncoding.EncodeToString(encrypt.RSASignSHA1(resp))
}
