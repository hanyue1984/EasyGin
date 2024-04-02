package routes

import "github.com/gin-gonic/gin"

func LoadRoutes(name string, route *gin.Engine) {
	SetupAdminRouter(name, route)
	SetupApiRouter(name, route)
	SetupAdminRBACRouter(name, route)
}
