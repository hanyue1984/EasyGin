package middleware

import (
	"EasyGin/common/lib"
	"EasyGin/common/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func HandleMiddleware(authIgnore bool, gatewayId string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 在这里执行中间件逻辑
		startTime := time.Now()
		result := gin.H{
			"cache":    false,
			"success":  false,
			"duration": -1,
			"hash":     tools.Common.GenerateUniqueHash(32),
			"error":    nil,
			"data":     nil,
		}
		ctx.Set("Models", nil)
		appid := ctx.GetHeader("appid")
		authHeader := ctx.GetHeader("Authorization")
		// 使用strings.Split函数按空格分割字符串
		parts := strings.Split(authHeader, " ")
		var token string
		// 检查是否得到正确的Bearer结构（至少包含两个部分）
		if len(authHeader) > 1 {
			if len(parts) > 1 && parts[0] == "Bearer" {
				// 提取出Bearer后面的令牌
				token = parts[1]
				ctx.Set("token", token)
			} else {
				result["error"] = gin.H{
					"code":    101,
					"message": "Token为非法格式",
				}
				ctx.JSON(http.StatusOK, result)
				ctx.Abort()
				return
			}
		}
		platform := ctx.GetHeader("platform")
		//deviceID := context.GetHeader("x-device-id")
		cookies := ctx.GetHeader("x-cookies")
		gatewayIdentification := ctx.GetHeader("gateway-identification")
		if authIgnore {
			if gatewayIdentification != gatewayId {
				result["error"] = gin.H{
					"code":    4002,
					"message": "您缺少访问权限",
				}
				ctx.JSON(http.StatusOK, result)
				ctx.Abort()
				return
			}
		}
		// Handle missing or undefined values
		if token == "" || token == "undefined" || token == "null" {
			token = tools.Common.GenerateUniqueHash(16)
		}

		if cookies == "" || cookies == "undefined" || cookies == "null" {
			cookies = fmt.Sprintf("cookie_%d_%s", time.Now().Unix(), tools.Common.GenerateUniqueHash(16))
		}

		if platform == "" || platform == "undefined" {
			platform = "default"
		}

		if appid == "" || appid == "undefined" || appid == "null" {
			result["error"] = gin.H{
				"code":    4001,
				"message": "您缺少必要参数",
			}
			ctx.JSON(http.StatusOK, result)
			ctx.Abort()
			return
		}
		defer func() {
			if r := recover(); r != nil {
				// 处理错误信息
				result["success"] = false
				if err, ok := r.(lib.HTTPError); ok {
					result["error"] = gin.H{
						"code":    err.Code,
						"message": err.Message,
					}
				} else {
					// 无法识别的错误类型
					result["error"] = gin.H{
						"code":    http.StatusInternalServerError,
						"message": "内部服务器错误",
					}
				}
				logError(r, ctx)
				ctx.JSON(http.StatusOK, result)
				ctx.Abort()
			}
		}()
		ctx.Next()
		data, _ := ctx.Get("data")
		result["data"] = data
		if data != nil {
			result["success"] = true
		}
		//清除数据
		ctx.Set("data", nil)
		endTime := time.Now()
		duration := endTime.Sub(startTime).Milliseconds()
		result["duration"] = duration
		if exp, ok := ctx.Get("exp"); ok {
			result["exp"] = exp
		}

		if ctx.Writer.Status() != http.StatusOK {
			if ctx.Writer.Size() == -1 {
				// Set default error response
				result["success"] = false
				result["error"] = gin.H{
					"code":    ctx.Writer.Status(),
					"message": "接口地址不正确",
				}
			}
		}

		ctx.Header("access-token", token)
		ctx.Header("x-cookies", cookies)
		ctx.Header("platform", platform)
		ctx.Header("response-success", "true")
		// Logging logic goes here
		ctx.JSON(http.StatusOK, result)
	}
}

// logError 日志记录错误信息的函数
func logError(err interface{}, context *gin.Context) {
	// 假设这里实现了错误日志的记录逻辑
	fmt.Println("Error:", err, "Request ID:", context.Request.Header.Get("X-Request-Id"))
}
