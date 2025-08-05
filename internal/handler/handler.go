package handler

import (
	"honoka-chan/config"
)

var (
	SifCdnServer string
	ErrorMsg     = `{"code":20001,"message":""}`
)

func init() {
	SifCdnServer = config.Conf.Settings.SifCdnServer
}
