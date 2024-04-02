package routes

import (
	Config "EasyGin/app/config"
	"EasyGin/app/controllers"
	"EasyGin/common/middleware"
	"EasyGin/common/tools"
	"github.com/gin-gonic/gin"
)

func SetupAdminRBACRouter(name string, engine *gin.Engine) {
	UserRedis := tools.RedisClient{}.Connect("User", Config.AppConfig.RedisCommon)
	AdminPowerAuth := middleware.AdminPowerAuthMiddleware(UserRedis)
	RBACController := controllers.RBACController{}
	router := engine.Group("/rbac/admin")
	{
		//获取用户列表
		router.GET("/v1/system/user/list", AdminPowerAuth, RBACController.UsersList)
		//设置用户权限
		router.PATCH("/v1/system/user", AdminPowerAuth, RBACController.EditUserRole)
		//获取角色列表
		router.GET("/v1/system/role", AdminPowerAuth, RBACController.RoleList)

		router.GET("/v1/system/role/list", AdminPowerAuth, RBACController.RoleAllList)
		//设置角色
		router.PATCH("/v1/system/role", AdminPowerAuth, RBACController.EditRole)
		//创建一个角色
		router.POST("/v1/system/role", AdminPowerAuth, RBACController.CreateRole)
		//删除
		router.DELETE("/v1/system/role", AdminPowerAuth, RBACController.DelRole)

		//创建菜单
		router.POST("/v1/system/menu", AdminPowerAuth, RBACController.CreateMenu)
		//设置菜单
		router.PATCH("/v1/system/menu", AdminPowerAuth, RBACController.EditMenu)
		//获取后台菜单列表
		router.GET("/v1/system/menu", AdminPowerAuth, RBACController.GetListMenu)
		//获取后台菜单列表
		router.GET("/v1/system/menu/my", AdminPowerAuth, RBACController.GetMyMenu)
		//删除后台菜单
		router.DELETE("/v1/system/menu", AdminPowerAuth, RBACController.DeleteMenu)

		//创建部门
		router.POST("/v1/system/dept", AdminPowerAuth, RBACController.CreateDept)
		//编辑部门
		router.PATCH("/v1/system/dept", AdminPowerAuth, RBACController.EditDept)
		//获取部门列表
		router.GET("/v1/system/dept/list", AdminPowerAuth, RBACController.GetDeptList)
		//删除部门
		router.DELETE("/v1/system/dept", AdminPowerAuth, RBACController.DeleteDept)
	}
}
