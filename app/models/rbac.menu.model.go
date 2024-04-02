package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"sort"
	"time"
)

type RbacMenuApi struct {
	Code string `json:"code"` //标识
	Url  string `json:"url"`  //接口
}
type RbacMenu struct {
	ID           uint            `gorm:"primaryKey;autoIncrement"`
	Name         string          `gorm:"not null"` //名称 路由标识 要跟组件name一致否则 <keep-alive>失效
	Path         string          `gorm:"not null"` //路由地址
	Meta         RbacMenuMeta    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;references:ID"`
	RbacMenuApis []RbacMenuApi   `gorm:"type:jsonb" json:"rbac_menu_apis"` // 使用jsonb类型存储数组
	Component    string          //视图名称 /views/下的
	Parent       uint            `gorm:"default:0"` //上级菜单 0为顶级菜单 顶级菜单不一定有组件有可能直接指向一个页面
	Sort         int             `gorm:"default:99"`
	Version      string          `gorm:"not null;default:'1.0'"`
	Active       string          //需要高亮的上级菜单路由地址
	Redirect     string          //重定向地址
	Extension    json.RawMessage `gorm:"type:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"` //软删除 并且创建一个索引提高软删除操作的性能
}

func (r *RbacMenu) Create() error {
	route := db.Create(r)
	if route.Error != nil {
		return route.Error
	}
	return nil
}
func (r *RbacMenu) Edit(menu RbacMenu) error {
	var Menu RbacMenu
	result := db.Preload("Meta").Where("id = ?", menu.ID).First(&Menu)
	if result.Error != nil {
		return result.Error
	}
	menu.Meta.ID = Menu.Meta.ID
	menu.Meta.RbacMenuID = Menu.Meta.ID
	db.Save(&menu)
	db.Save(&menu.Meta)
	return nil
}

// FindMenusByNameAndParentID 根据名称和父ID获取菜单
func (r *RbacMenu) FindMenusByNameAndParentID(Names []string, ParentId uint) ([]RbacMenu, error) {
	//namesStr := strings.Join(Names, ",")
	var menus []RbacMenu
	if err := db.Preload("Meta").Where("name IN (?) AND parent = ?", Names, ParentId).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// IsStringInSortedSlice 判断字符串是否在已排序的切片中
func (r *RbacMenu) IsStringInSortedSlice(menus []string, target string) bool {
	sort.Strings(menus) // 如果切片尚未排序，则先排序
	index := sort.SearchStrings(menus, target)
	return index < len(menus) && menus[index] == target
}

// FindParentAll 查找所有Parent为0的数据
func (r *RbacMenu) FindParentAll(ParentId uint) ([]RbacMenu, error) {
	// 查找所有Parent为0的数据
	var topLevelRoutes []RbacMenu
	db.Preload("Meta").Where("parent = ?", ParentId).Order("sort DESC").Find(&topLevelRoutes)
	return topLevelRoutes, nil
}
func (r *RbacMenu) FindAll(id int) ([]RbacMenu, error) {
	// 查找所有Parent为0的数据
	var topLevelRoutes []RbacMenu
	db.Where("id = ?", id).Order("sort DESC").Find(&topLevelRoutes)
	return topLevelRoutes, nil
}
func (r *RbacMenu) Delete(id uint) error {
	var menu RbacMenu
	if err := db.Preload("Meta").Where("id = ?", id).Unscoped().Delete(&menu).Error; err != nil {
		return err
	}
	return nil
}

func (r *RbacMenu) Save() {
	db.Save(r)
}
