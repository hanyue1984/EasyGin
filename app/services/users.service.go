package services

import (
	Config "EasyGin/app/config"
	"EasyGin/app/models"
	"EasyGin/common/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type UsersService struct {
}

func (u UsersService) GetUser(ctx *gin.Context, key string) *models.Users {
	var UsersRedis tools.RedisClient
	var user models.Users
	client := UsersRedis.Connect("User", Config.AppConfig.RedisCommon)
	result := client.Get(ctx, key, models.Users{})
	// 检查是否成功转换
	if user, ok := result.(models.Users); ok {
		// 成功转换为models.User类型
		fmt.Println("User:", user)
	} else {
		fmt.Println("无法将object转换为models.User类型")
		return nil
	}
	return &user
}

func (u UsersService) SetUser(ctx *gin.Context, key string) bool {
	var UsersRedis *tools.RedisClient
	UsersRedis.Connect("User", Config.AppConfig.RedisCommon)
	UsersRedis.Set(ctx, key, models.Users{
		Hash:      "12345",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}, 0)
	return true
}
