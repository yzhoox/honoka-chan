package handler

import (
	"honoka-chan/config"
	"honoka-chan/pkg/db"

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

	MainEng = db.MainEng
	UserEng = db.UserEng
}
