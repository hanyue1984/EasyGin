package routes

import (
	"EasyGin/app/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupApiRouter(name string, engine *gin.Engine) {
	router := engine.Group(fmt.Sprintf("/%s/api", name))
	{
		router.GET("/v1/get", controllers.UserController{}.GetUser)
	}
}
