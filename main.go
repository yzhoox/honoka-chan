package main

import (
	"honoka-chan/config"
	"honoka-chan/db"
	"honoka-chan/router"
	_ "honoka-chan/tools"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// 处理系统信号，确保程序退出时关闭数据库
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Shutting down...")
		db.DB.Close()
		os.Exit(0)
	}()

	// Gin
	gin.SetMode(gin.ReleaseMode)

	// Router
	r := gin.Default()

	// SIF
	router.SifRouter(r)

	r.Run(":" + config.Conf.Settings.ListenPort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
