package middleware

import (
	"EasyGin/common/tools"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(client tools.RedisClient) gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}
