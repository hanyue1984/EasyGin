package models

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	id        uint   `gorm:"primaryKey"` //主键自增
	Hash      string `gorm:"unique"`     //唯一
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` //软删除 并且创建一个索引提高软删除操作的性能
}
