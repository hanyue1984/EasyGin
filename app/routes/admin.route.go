package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupAdminRouter(name string, engine *gin.Engine) {
	router := engine.Group(fmt.Sprintf("/%s/admin", name))
	{
		router.GET("/v1/index", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "Post request GIN Admin API",
			})
		})
	}
}
