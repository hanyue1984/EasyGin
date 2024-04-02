package controllers

import (
	"EasyGin/app/models"
	"EasyGin/app/services"
	"EasyGin/common/lib"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type UserController struct {
}

func (u UserController) GetUser(ctx *gin.Context) {
	//user := services.UsersService{}.GetUser(ctx, "12345")
	user := &models.Users{
		Hash: "xxxxx",
		//AppId:     "12",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
	if user == nil {
		lib.CustomError(500, "没有该用户")
	}
	ctx.Set("data", gin.H{
		"hash": user.Hash,
		//"Appid": user.AppId,
	})
}

type registerRequest struct {
	Platform string `json:"platform"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	UserType string `json:"userType"`
}

func (u UserController) Login(ctx *gin.Context) {
	var requestData registerRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		lib.CustomError(100, "参数为空缺少必要参数platform")
	}
	if requestData.Platform == "" {
		lib.CustomError(100, "参数为空缺少必要参数platform")
	}
	switch requestData.Platform {
	case "password":
		service := services.UsersService{}
		ctx.Set("data", service.UsernameLogin(ctx, requestData.Platform, requestData.Username, requestData.Password))
	case "wechat":
		lib.CustomError(100, "暂未开通微信登录")
	case "qq":
		lib.CustomError(100, "暂未开通QQ登录")
	case "mobile":
		lib.CustomError(100, "暂未开通手机登录")
	default:
		lib.CustomError(102, "传入错误的平台类型")
	}
}

func (u UserController) Register(ctx *gin.Context) {
	var requestData registerRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		lib.CustomError(100, "参数为空缺少必要参数platform")
	}
	platform := requestData.Platform
	if platform == "" {
		lib.CustomError(100, "参数为空缺少必要参数platform")
	}
	switch platform {
	case "password":
		service := services.UsersService{}
		ctx.Set("data", service.Register(ctx, platform, requestData.Username, requestData.Password, requestData.Nickname, requestData.Email, requestData.UserType))
	case "wechat":
		lib.CustomError(100, "暂未开通微信登录")
	case "qq":
		lib.CustomError(100, "暂未开通QQ登录")
	case "mobile":
		lib.CustomError(100, "暂未开通手机登录")
	default:
		lib.CustomError(102, "传入错误的平台类型")
	}
}

func (u UserController) AdminToken(ctx *gin.Context) {
	username := ""
	password := ""
	var data map[string]interface{}
	if err := ctx.BindJSON(&data); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}
	if val, ok := data["username"]; ok {
		username = val.(string)
	}

	if val, ok := data["password"]; ok {
		password = val.(string)
	}
	if username != "" && password != "" {
		service := services.UsersService{}
		responseData := service.UsernameLogin(ctx, "password", username, password)
		token := responseData["token"].(string)
		userInfo := responseData["userInfo"].(*models.Users)
		ctx.Set("data", gin.H{
			"token": token,
			"userInfo": map[string]interface{}{
				"userId":    userInfo.Hash,
				"userName":  userInfo.Nickname,
				"dashboard": "0",
				"dept":      userInfo.Dept,
				"admin":     userInfo.Admin,
			},
		})
	} else {
		lib.CustomError(100, "缺少必要参数 username password")
	}
}
