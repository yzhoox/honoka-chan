package main

import (
	"context"
	"honoka-chan/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := app.Start("."); err != nil {
		log.Fatalln(err)
	}

	// 处理系统信号，确保程序退出时关闭数据库
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("正在退出...")
	if err := app.Stop(context.Background()); err != nil {
		log.Println(err.Error())
	}
	os.Exit(0)
}
