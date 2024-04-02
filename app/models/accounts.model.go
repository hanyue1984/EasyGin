package models

import (
	"EasyGin/common/tools"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Accounts struct {
	Id             uint            `gorm:"primaryKey"`      //主键自增
	Hash           string          `gorm:"not null;unique"` //对应User的hash唯一
	AppId          string          `gorm:"not null"`        //APPID
	Token          string          `gorm:"not null"`        //TOKEN
	Platform       string          `gorm:"not null"`        //平台类型
	PlatformUserID string          //对应平台ID
	Username       string          //用户名 当platform为手机号时 username为手机号
	Password       string          //密码或者凭证
	Extension      json.RawMessage `gorm:"type:json"` //扩展参数
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"` //软删除 并且创建一个索引提高软删除操作的性能
}

func (a *Accounts) Create(appid string, hash string, platform string, platformId string, username string, password string) (*Accounts, error) {
	newAccount := Accounts{
		Hash:           hash,
		AppId:          appid,
		Token:          a.CreateToken(),
		Platform:       platform,
		PlatformUserID: platformId,
		Username:       username,
		Password:       tools.Common.GenerateMD5Hash(password),
	}
	result := db.Create(&newAccount)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newAccount, nil
}

func (a *Accounts) FindOne(platform string, username string, appid string) (*Accounts, error) {
	var account Accounts
	result := db.Where("platform = ? AND username = ? AND app_id = ?", platform, username, appid).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	} else {
		return &account, nil
	}
}

func (a *Accounts) FindOneByPlatform(platform string) (*Accounts, error) {
	var account Accounts
	result := db.Where("platform = ?", platform).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	} else {
		return &account, nil
	}
}
func (a *Accounts) CreateToken() string {
	a.Token = tools.Common.GenerateMD5Hash(tools.Common.GenerateUniqueUid())
	return a.Token
}

func (a *Accounts) Save() error {
	return db.Save(a).Error
}
