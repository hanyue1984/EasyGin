package routes

import (
	"EasyGin/app/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupAdminRouter(name string, engine *gin.Engine) {
	router := engine.Group(fmt.Sprintf("/%s/admin", name))
	{
		var UserController controllers.UserController
		//后台登录
		router.POST("/v1/token", UserController.AdminToken)
		router.POST("/v1/register", UserController.Register)
	}
}
