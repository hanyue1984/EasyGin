package main

import (
	Config "EasyGin/app/config"
	"EasyGin/app/models"
	"EasyGin/app/routes"
	"EasyGin/app/services"
	"EasyGin/common/middleware"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func setupLogger() {
	// 打开一个文件，用于写入日志
	logsDir := "./logs/"
	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		log.Fatalf("无法创建日志目录：%v", err)
	}
	file, err := os.OpenFile("./logs/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法打开日志文件: ", err)
	}
	// 设置日志输出到文件
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
}
func main() {
	// 设置日志输出到文件
	setupLogger()
	// 解析命令行参数以获取当前环境
	envPtr := flag.String("env", "dev", "Specify the environment to run the application")
	flag.Parse()

	currentEnvironment := *envPtr

	if currentEnvironment == "api" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()
	// 处理对 /favicon.ico 的请求
	router.GET("/favicon.ico", func(c *gin.Context) {
		// 返回一个空的响应或者合适的错误响应
		c.String(204, "")
	})
	router.SetTrustedProxies([]string{"127.0.0.1"})
	// 加载配置
	config := Config.LoadConfig(currentEnvironment)
	//注册公共中间件
	router.Use(middleware.HandleMiddleware(true, config.GatewayIdentification))

	Config.AppConfig = config

	routes.LoadRoutes(fmt.Sprintf("%s", config.Name), router) // 自动加载路由文件

	models.ConnectDB(config.Database) //连接数据库

	services.ListenService{}.Listen()

	router.Run(fmt.Sprintf(":%v", config.Port)) // 启动 Web 服务器
}
