package startup

import (
	"fmt"
	"honoka-chan/internal/handler/ghome/account"
	usermodel "honoka-chan/internal/model/user"
	"log"
)

const (
	defaultPassword = "klsbgames"
)

func CreateDefaultUser() error {
	_, code, msg, created, err := account.AddUser(usermodel.DefaultSystemPhone, defaultPassword, true)
	if err != nil {
		return fmt.Errorf("默认用户创建失败: %w", err)
	}
	if code != 0 {
		return fmt.Errorf("默认用户创建失败: code=%d msg=%s", code, msg)
	}
	if created {
		log.Printf("默认用户创建成功, 账号: %s, 密码: %s\n", usermodel.DefaultSystemPhone, defaultPassword)
		return nil
	}
	return nil
}
