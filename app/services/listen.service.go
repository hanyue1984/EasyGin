package services

import (
	Config "EasyGin/app/config"
	"EasyGin/common/tools"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type ListenService struct {
}

func (listen ListenService) Listen() {
	UsersRedis := tools.RedisClient{}.Connect("User", Config.AppConfig.RedisCommon)
	UsersRedis.Subscribe("users_modules", func(message *redis.Message) {
		fmt.Sprintln(message)
	})
}
