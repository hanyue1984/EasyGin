package middleware

import (
	"EasyGin/app/models"
	"EasyGin/common/lib"
	"EasyGin/common/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

func AdminPowerAuthMiddleware(client *tools.RedisClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Get("token")
		if token == "" || !err {
			lib.CustomError(401, "token is empty")
			return
		}
		token = token.(string)
		var tokenBody struct {
			Hash     string
			Token    string
			Platform string
			AppId    string
		}
		tokenBodyErr := client.Get(ctx, fmt.Sprintf("token_body:token_%s", token), &tokenBody)
		if tokenBodyErr != nil {
			lib.CustomError(401, "token is invalid")
		}
		var user models.Users
		userErr := client.Get(ctx, fmt.Sprintf("user_body:hash_%s", tokenBody.Hash), &user)
		if userErr != nil {
			lib.CustomError(401, "user is invalid")
		}
		if user.Admin {
			ctx.Set("UserInfo", user)
			ctx.Next()
		} else {
			lib.CustomError(401, "user is not admin")
		}
	}
}
