package startup

import (
	"honoka-chan/internal/handler/ghome/account"
	usermodel "honoka-chan/internal/model/user"
	"log"
)

const (
	defaultPassword = "klsbgames"
)

func CreateDefaultUser() {
	_, code, msg, created, err := account.AddUser(usermodel.DefaultSystemPhone, defaultPassword, true)
	if err != nil {
		log.Fatalln("默认用户创建失败:", err.Error())
	}
	if code != 0 {
		log.Fatalf("默认用户创建失败: code=%d msg=%s", code, msg)
	}
	if created {
		log.Printf("默认用户创建成功, 账号: %s, 密码: %s\n", usermodel.DefaultSystemPhone, defaultPassword)
		return
	}
}
