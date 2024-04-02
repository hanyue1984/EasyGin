package models

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type RbacDept struct {
	ID        uint            `gorm:"primaryKey;autoIncrement"`
	Label     string          `gorm:"unique;not null"` //名称
	ParentId  int             `gorm:"not null"`        //当为0时为1级目录
	Remarks   string          //备注
	Status    int             `gorm:"default:0"`  //状态0 为停用 1为启用
	Sort      int             `gorm:"default:99"` //排序
	Extension json.RawMessage `gorm:"type:json"`  //扩展
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` //软删除 并且创建一个索引提高软删除操作的性能
}

func (g *RbacDept) Save() {
	db.Save(&g)
}

func (g *RbacDept) Create(label string, parentId int, remarks string, status int, sort int) (*RbacDept, error) {
	group := RbacDept{
		Label:    label,
		ParentId: parentId,
		Remarks:  remarks,
		Status:   status,
		Sort:     sort,
	}
	result := db.Create(&group)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}
func (g *RbacDept) Edit(id int, label string, parentId int, remarks string, status int, sort int) bool {
	if g.FindOne(id) == nil {
		//查找到了这个表
		g.Label = label
		g.ParentId = parentId
		g.Remarks = remarks
		g.Status = status
		g.Sort = sort
		g.Save()
		return true
	}
	return false
}

func (g *RbacDept) FindAll(parentId uint) []RbacDept {
	var groups []RbacDept
	db.Where("parent_id = ?", parentId).Order("sort DESC").Find(&groups)
	return groups
}
func (g *RbacDept) FindOne(id int) error {
	result := db.Where("id = ?", id).First(g)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return result.Error
		}
	} else {
		return nil
	}
}

func (g *RbacDept) Delete(id uint) error {
	var Dept RbacDept
	if err := db.Where("id = ?", id).Unscoped().Delete(&Dept).Error; err != nil {
		return err
	}
	return nil
}
