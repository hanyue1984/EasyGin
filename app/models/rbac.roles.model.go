package models

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type RbacRoles struct {
	ID        uint            `gorm:"primaryKey;autoIncrement"`
	Label     string          `gorm:"unique;not null"` //角色名称
	Alias     string          `gorm:"not null"`        //角色别名
	Remark    string          //备注 可以为空
	Status    int             `gorm:"default:0"`   // 0为关闭角色 1为开启
	Sort      int             `gorm:"default:99"`  //排序 数字越大越在前面
	Menus     pq.StringArray  `gorm:"type:text[]"` //菜单列表
	Extension json.RawMessage `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` //软删除 并且创建一个索引提高软删除操作的性能
}

func (r *RbacRoles) Save() {
	db.Save(&r)
}
func (r *RbacRoles) Create(label string, alias string, remark string, status int, sort int) (*RbacRoles, error) {
	group := RbacRoles{
		Label:  label,
		Alias:  alias,
		Remark: remark,
		Status: status,
		Sort:   sort,
		Menus:  []string{},
	}
	result := db.Create(&group)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}
func (r *RbacRoles) Edit(id int, label string, alias string, remark string, status int, sort int, menus []string) bool {
	role, err := r.FindOne(id)
	if err != nil {
		return false
	}
	role.Label = label
	role.Alias = alias
	role.Remark = remark
	role.Status = status
	role.Sort = sort
	role.Menus = menus
	role.Save()
	return true
}

func (r *RbacRoles) FindAll(page int, limit int) ([]gin.H, int64) {
	var total int64
	// 获取总记录数
	db.Model(r).Count(&total)
	var groups []RbacRoles
	// 查询特定页的数据
	db.Offset((page - 1) * limit).Limit(limit).Find(&groups)
	// 将查询结果转换为 []interface{} 类型的切片
	var interfaces []gin.H
	for _, group := range groups {
		interfaces = append(interfaces, gin.H{
			"id":     strconv.Itoa(int(group.ID)),
			"label":  group.Label,
			"alias":  group.Alias,
			"remark": group.Remark,
			"status": strconv.Itoa(group.Status),
			"sort":   group.Sort,
			"data":   group.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return interfaces, total
}

func (r *RbacRoles) FindOne(id int) (RbacRoles, error) {
	var role RbacRoles
	result := db.Where("id = ?", id).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return role, nil
		} else {
			return role, result.Error
		}
	} else {
		return role, nil
	}
}

func (r *RbacRoles) Delete(id uint) error {
	var Role RbacRoles
	if err := db.Where("id = ?", id).Unscoped().Delete(&Role).Error; err != nil {
		return err
	}
	return nil
}

func (r *RbacRoles) EditMenus(menus []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range menus {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
