package controllers

import (
	"EasyGin/app/services"
	"github.com/gin-gonic/gin"
)

type RBACController struct {
}

func (r *RBACController) RoleList(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.RolesList(ctx))
}
func (r *RBACController) RoleAllList(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.RolesAllList(ctx))
}
func (r *RBACController) EditRole(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.RolesEdit(ctx))
}
func (r *RBACController) DelRole(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.RolesDelete(ctx))
}
func (r *RBACController) CreateRole(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.RolesCreate(ctx))
}
func (r *RBACController) SetUserRole(ctx *gin.Context) {

}

func (r *RBACController) GetListMenu(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.MenuList())
}
func (r *RBACController) GetMyMenu(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.MenuMy(ctx))
}
func (r *RBACController) CreateMenu(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.MenuCreate(ctx))
}
func (r *RBACController) EditMenu(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.MenuEdit(ctx))
}
func (r *RBACController) DeleteMenu(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.MenuDelete(ctx))
}

// CreateDept 创建组
func (r *RBACController) CreateDept(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.DeptCreate(ctx))
}

// EditDept 编辑组
func (r *RBACController) EditDept(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.DeptEdit(ctx))
}

// GetDeptList 获取所有组
func (r *RBACController) GetDeptList(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.DeptList())
}
func (r *RBACController) DeleteDept(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.DeptDelete(ctx))
}

func (r *RBACController) EditUserRole(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.UsersEdit(ctx))
}

func (r *RBACController) UsersList(ctx *gin.Context) {
	var service services.RBAC
	ctx.Set("data", service.UsersList(ctx))
}
