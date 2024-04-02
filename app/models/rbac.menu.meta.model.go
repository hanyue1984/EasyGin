package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type RbacMenuMeta struct {
	ID               uint `gorm:"primaryKey;autoIncrement"`
	RbacMenuID       uint
	Title            string          `gorm:"not null"`                //显示名称
	Hidden           bool            `gorm:"default:false"`           //是否隐藏菜单
	HiddenBreadcrumb bool            `gorm:"default:false"`           //是否隐藏面包屑
	FullPage         bool            `gorm:"default:false"`           //整页路由
	Icon             string          `gorm:"default:'el-icon-apple'"` //图标
	Type             string          `gorm:"not null"`                //菜单、IFrame、外链、按钮
	Color            string          //颜色
	Tag              string          //标签
	Affix            bool            //是否固定
	Extension        json.RawMessage `gorm:"type:json"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"` //软删除 并且创建一个索引提高软删除操作的性能
}

func (r *RbacMenuMeta) Save() {
	db.Save(&r)
}
