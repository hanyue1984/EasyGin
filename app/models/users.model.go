package models

import (
	"EasyGin/common/tools"
	"encoding/json"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Users struct {
	ID            uint            `gorm:"primaryKey"`      //主键自增
	Hash          string          `gorm:"not null;unique"` //唯一
	Nickname      string          `gorm:"not null;default:平台游客"`
	Name          string          //姓名
	AppId         string          `gorm:"not null"`  //APPID
	Sex           uint8           `gorm:"default:0"` //性别0 未知 1男 2女
	Signature     string          //签名
	Email         string          //邮箱
	Address       string          //省市县(区)
	Avatar        string          //头像地址
	Birthday      time.Time       //生日
	Mobile        string          //手机号
	State         uint8           `gorm:"default:1"` //1正常用户
	NoTalk        bool            `gorm:"default:false"`
	LoginCount    uint            `gorm:"default:1"`
	LastLoginTime time.Time       //最后一次登录时间
	Admin         bool            `gorm:"default:false"` //是否是管理员
	Roles         pq.StringArray  `gorm:"type:text[]"`   //所属角色
	Dept          int             `gorm:"default:0"`     //所属部门
	DeviceId      string          //设备ID
	Extension     json.RawMessage `gorm:"type:json"` //扩展参数
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"` //软删除 并且创建一个索引提高软删除操作的性能
}

func (u *Users) Create(appid string, nickname string, email string) (*Users, error) {
	newUser := Users{
		Hash:          tools.Common.GenerateUniqueUid(),
		AppId:         appid,
		Nickname:      nickname,
		Email:         email,
		LastLoginTime: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	result := db.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newUser, nil
}

func (u *Users) FindOne(query interface{}, args ...interface{}) (*Users, error) {
	var user Users
	result := db.Where(query, args).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *Users) FindById(Uid string) (*Users, error) {
	var user Users
	result := db.Where("uid = ?", Uid).First(&user)
	if result.Error != nil {
		return &user, result.Error
	}
	return &user, nil
}
func (u *Users) FindAll(page int, limit int, search string, groupId string) ([]Users, int64) {
	var total int64
	// 获取总记录数
	db.Model(u).Count(&total)
	var users []Users
	// 查询特定页的数据
	if len(search) > 0 {
		db.Where("nickname LIKE ?", "%"+search+"%").Offset((page - 1) * limit).Limit(limit).Find(&users)
	} else if groupId != "" {
		db.Where("group = ?", groupId).Offset((page - 1) * limit).Limit(limit).Find(&users)
	} else {
		db.Offset((page - 1) * limit).Limit(limit).Find(&users)
	}
	if total > 0 {
		return users, total
	}
	return nil, 0
}
func (u *Users) Save() {
	db.Save(&u)
}
