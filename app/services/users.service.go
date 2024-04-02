package services

import (
	Config "EasyGin/app/config"
	"EasyGin/app/models"
	"EasyGin/common/lib"
	"EasyGin/common/tools"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type UsersService struct {
	Hash     string
	Token    string
	Platform string
	AppId    string
}

var UsersModel models.Users
var AccountsModel models.Accounts

func (u *UsersService) UsernameLogin(ctx *gin.Context, platform string, username string, password string) gin.H {
	appid := ctx.GetHeader("appid")
	if username == "" || password == "" {
		lib.CustomError(100, "缺少必要参数 username password")
	}
	account, err := AccountsModel.FindOne(platform, username, appid)
	if err != nil || account == nil {
		lib.CustomError(101, "用户名密码错误请重新尝试")
	}
	users := models.Users{}
	user, err := users.FindOne("hash = ?", account.Hash)
	if err != nil {
		lib.CustomError(101, "用户名密码错误请重新尝试")
	}
	user.LoginCount += 1
	user.LastLoginTime = time.Now()
	account.CreateToken()
	err = account.Save()
	if err != nil {
		return nil
	}
	user.Save()
	u.Hash = user.Hash
	u.Token = account.Token
	u.Platform = platform
	u.AppId = appid
	err = u.tLogin(ctx, user)
	if err != nil {
		return nil
	}
	return gin.H{
		"token":    account.Token,
		"userInfo": user,
	}
}

func (u *UsersService) Register(ctx *gin.Context, platform string, username string, password string, nickname string, email string, userType string) gin.H {
	switch platform {
	case "password":
		appid := ctx.GetHeader("appid")
		if username == "" || password == "" {
			lib.CustomError(100, "缺少必要参数 username password")
		}
		account, _ := AccountsModel.FindOne(platform, username, appid)
		if account == nil {
			user, _ := UsersModel.Create(appid, nickname, email)
			if user != nil {
				account, _ = AccountsModel.Create(appid, user.Hash, platform, "", username, password)
				if account != nil {
					u.Hash = user.Hash
					u.Token = account.Token
					u.Platform = platform
					u.AppId = appid
					err := u.tLogin(ctx, user)
					if err != nil {
						return nil
					}
					return gin.H{
						"token":    account.Token,
						"hash":     user.Hash,
						"nickname": user.Nickname,
					}
				}
			}
		} else {
			lib.CustomError(102, "该用户名已被注册")
		}
	case "mobile":
		lib.CustomError(100, "暂未开通手机注册")
	case "wechat":
		lib.CustomError(100, "暂未开通注册")
	default:
		lib.CustomError(100, "platform错误参数或者缺少")
	}
	return nil
}

func (u *UsersService) LoadUserInfo(ctx *gin.Context, hash string, userInfo interface{}) error {
	UserRedis := tools.RedisClient{}.Connect("User", Config.AppConfig.RedisCommon)
	HashExpiration := time.Hour * 24 * 15
	err := UserRedis.Get(ctx, fmt.Sprintf("user_body:hash_%s", hash), userInfo)
	if errors.Is(err, redis.Nil) { //如果redis中没有键值
		user, _ := UsersModel.FindOne("hash = ?", hash)
		if user != nil {
			//重置用户数据到redis
			UserErr := UserRedis.Set(ctx, fmt.Sprintf("user_body:hash_%s", user.Hash), user, HashExpiration)
			if UserErr != nil {
				return UserErr
			}
			if err != nil {
				return fmt.Errorf("error decoding JSON: %v", err)
			}
			userInfo = user
			return nil
		}
	} else if err != nil {
		return err
	}
	return nil
}
func (u *UsersService) tLogin(ctx *gin.Context, user *models.Users) error {
	TokenExpiration := time.Hour * 24 * 7 // 15 天的持续时间
	UsersRedis := tools.RedisClient{}.Connect("User", Config.AppConfig.RedisCommon)
	HashExpiration := time.Hour * 24 * 15
	TokenErr := UsersRedis.Set(ctx, fmt.Sprintf("token_body:token_%s", u.Token), u, TokenExpiration)
	if TokenErr != nil {
		return TokenErr
	} //3600秒后过期
	UserErr := UsersRedis.Set(ctx, fmt.Sprintf("user_body:hash_%s", user.Hash), user, HashExpiration)
	if UserErr != nil {
		return UserErr
	}
	return nil
}

func (u *UsersService) upExp(ctx *gin.Context, user *models.Users) error {
	TokenExpiration := time.Hour * 24 * 7 // 15 天的持续时间
	HashExpiration := time.Hour * 24 * 15
	UsersRedis := tools.RedisClient{}.Connect("User", Config.AppConfig.RedisCommon)
	TokenErr := UsersRedis.Expire(ctx, fmt.Sprintf("token_body:token_%s", u.Token), TokenExpiration)
	if TokenErr != nil {
		return TokenErr
	}
	UserErr := UsersRedis.Expire(ctx, fmt.Sprintf("user_body:hash_%s", user.Hash), HashExpiration)
	if UserErr != nil {
		return UserErr
	}
	return nil
}

func (u *UsersService) uCache(ctx *gin.Context, user *models.Users, force bool) error {
	HashExpiration := time.Hour * 24 * 15
	UsersRedis := tools.RedisClient{}.Connect("User", Config.AppConfig.RedisCommon)
	if user.State <= 0 { //如果账号被封
		UsersRedis.Destroy(ctx, fmt.Sprintf("token_body:token_%s", u.Token))
		UsersRedis.Destroy(ctx, fmt.Sprintf("user_body:hash_%s", user.Hash))
	} else {
		var userCache models.Users
		err := UsersRedis.Get(ctx, fmt.Sprintf("user_body:hash_%s", user.Hash), userCache)
		if err != nil {
			return err
		}
		if force {
			UsersRedis.Set(ctx, fmt.Sprintf("user_body:hash_%s", user.Hash), user, HashExpiration) //强制刷新
		}
	}
	message := tools.ListenData{
		Key:  "users",
		Cmd:  "user_info_update",
		Data: user,
	}
	UsersRedis.Publish("users_modules", message)
	return nil
}
