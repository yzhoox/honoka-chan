package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	Conf = &AppConfigs{}

	PackageVersion = "97.4.6"

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
	UnlockAllSpecialRotation bool   `json:"unlock_all_special_rotation"`
}

func InitConfig() error {
	conf, err := Load("./config.json")
	if err != nil {
		return err
	}
	Conf = conf
	return nil
}

func DefaultConfigs() *AppConfigs {
	return &AppConfigs{
		AppName: "honoka-chan",
		Settings: Settings{
			ListenPort:               "8080",
			CdnServer:                "http://127.0.0.1:8080/static",
			UnlockAllSpecialRotation: false,
		},
	}
}

func Load(p string) (*AppConfigs, error) {
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		conf := DefaultConfigs()
		if err := conf.Save(p); err != nil {
			return nil, fmt.Errorf("create default config: %w", err)
		}
		return conf, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	c := AppConfigs{}
	if err := json.Unmarshal(data, &c); err == nil {
		return &c, nil
	}

	backup := p + ".backup" + strconv.FormatInt(time.Now().Unix(), 10)
	if err := os.Rename(p, backup); err != nil {
		return nil, fmt.Errorf("backup invalid config: %w", err)
	}
	conf := DefaultConfigs()
	if err := conf.Save(p); err != nil {
		return nil, fmt.Errorf("write default config: %w", err)
	}
	return conf, nil
}

func (c *AppConfigs) Save(p string) error {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, append(data, '\n'), 0644)
}
