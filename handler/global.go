package handler

import (
	"honoka-chan/config"

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
