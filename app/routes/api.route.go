package routes

import (
	"EasyGin/app/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupApiRouter(name string, engine *gin.Engine) {
	router := engine.Group(fmt.Sprintf("/%s/api", name))
	{
		var UserController controllers.UserController
		router.GET("/v1/get", UserController.GetUser)
		router.POST("/v1/login", UserController.Login)
		router.POST("/v1/register", UserController.Register)
	}
}
