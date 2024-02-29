package main

import (
	Config "EasyGin/app/config"
	"EasyGin/app/models"
	"EasyGin/app/routes"
	"EasyGin/common/middleware"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// 解析命令行参数以获取当前环境
	envPtr := flag.String("env", "dev", "Specify the environment to run the application")
	flag.Parse()

	currentEnvironment := *envPtr

	router := gin.Default()

	//注册公共中间件
	router.Use(middleware.HandleMiddleware())

	config := Config.LoadConfig(currentEnvironment) // 加载配置

	Config.AppConfig = config

	routes.LoadRoutes(fmt.Sprintf("%s", config.Name), router) // 自动加载路由文件

	models.ConnectDB(config.Database) //连接数据库

	router.Run(fmt.Sprintf(":%v", config.Port)) // 启动 Web 服务器
}
