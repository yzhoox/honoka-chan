package config

import (
	"encoding/json"
	"honoka-chan/pkg/utils"
	"os"
	"strconv"
	"time"
)

var (
	Conf = &AppConfigs{}

	PackageVersion = "97.4.6"
	ConfigPath     = "./config.json"

	PrivateKeyPath = "assets/certs/privatekey.pem"
	PublicKeyPath  = "assets/certs/publickey.pem"
)

type AppConfigs struct {
	AppName  string   `json:"app_name"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	ListenPort               string `json:"listen_port"`
	CdnServer                string `json:"cdn_server"`
	ReloadToken              string `json:"reload_token"`
	UnlockAllSpecialRotation bool   `json:"unlock_all_special_rotation"`
}

func InitConfig() {
	ReloadConfig()
}

func ReloadConfig() *AppConfigs {
	Conf = Load(ConfigPath)
	return Conf
}

func DefaultConfigs() *AppConfigs {
	return &AppConfigs{
		AppName: "honoka-chan",
		Settings: Settings{
			ListenPort:               "8080",
			CdnServer:                "http://127.0.0.1:8080/static",
			ReloadToken:              "",
			UnlockAllSpecialRotation: false,
		},
	}
}

func Load(p string) *AppConfigs {
	if !utils.PathExists(p) {
		_ = DefaultConfigs().Save(p)
	}
	c := AppConfigs{}
	err := json.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	if err != nil {
		_ = os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		_ = DefaultConfigs().Save(p)
	}
	c = AppConfigs{}
	_ = json.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	return &c
}

func (c *AppConfigs) Save(p string) error {
	data, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return err
	}
	utils.WriteAllText(p, string(data)+"\n")
	return nil
}
