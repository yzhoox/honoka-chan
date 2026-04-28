package startup

import (
	"honoka-chan/internal/handler/ghome/account"
	"log"
)

const (
	defaultPhone    = "1"
	defaultPassword = "klsbgames"
)

func CreateDefaultUser() {
	_, code, msg, created, err := account.AddUser(defaultPhone, defaultPassword, true)
	if err != nil {
		log.Fatalln("默认用户创建失败:", err.Error())
	}
	if code != 0 {
		log.Fatalf("默认用户创建失败: code=%d msg=%s", code, msg)
	}
	if created {
		log.Printf("默认用户创建成功, 账号: %s, 密码: %s\n", defaultPhone, defaultPassword)
		return
	}
}
