package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func HandleMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 在这里执行中间件逻辑
		start := time.Now()
		fmt.Println("执行中间件")
		context.Next()
		if !context.IsAborted() {
			fmt.Println("Log: Request handled successfully")
		}
		elapsed := time.Since(start)
		fmt.Printf("Request processed in %v\n", elapsed)
		err := context.Errors.Last()
		if err != nil {
			// 处理错误，比如记录日志或返回特定的错误信息给客户端
		}
	}
}
