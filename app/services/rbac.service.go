package services

import (
	"EasyGin/app/models"
	"EasyGin/common/lib"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type RBAC struct {
}

type groupRequest struct {
	Id       int    `json:"id"`
	Label    string `json:"label"`
	ParentId int    `json:"parentId"`
	Remark   string `json:"remark"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}

type MetaData struct {
	Title            string `json:"title"`
	Type             string `json:"type"`
	Hidden           bool   `json:"hidden"`
	HiddenBreadcrumb bool   `json:"hiddenBreadcrumb"`
	FullPage         bool   `json:"fullPage"`
	Icon             string `json:"icon"`
	Color            string `json:"color"`
	Tag              string `json:"tag"`
	Affix            bool   `json:"affix"`
}
type MenuJSONData struct {
	ID        int      `json:"id"`
	Sort      int      `json:"sort"`
	Status    string   `json:"status"`
	Remark    string   `json:"remark"`
	Alias     string   `json:"alias"`
	Label     string   `json:"label"`
	ParentID  int      `json:"parentId"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Redirect  string   `json:"redirect"`
	Component string   `json:"component"`
	Active    string   `json:"active"`
	Meta      MetaData `json:"meta"`
}
type MenuDeleteJsonData struct {
	Ids []int `json:"ids"`
}

type userEditJson struct {
	ID       int      `json:"id"`
	UserName string   `json:"userName"`
	Avatar   string   `json:"avatar"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	Dept     int      `json:"dept"`
}

// ======================Users==

func (r *RBAC) UsersEdit(ctx *gin.Context) uint {
	var requestData userEditJson
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		lib.CustomError(100, fmt.Sprintf("缺少必要参数Error:%v", err))
	}
	var users models.Users
	one, err := users.FindOne("id = ?", requestData.ID)
	if err != nil {
		return 0
	}
	one.Nickname = requestData.UserName
	one.Avatar = requestData.Avatar
	one.Name = requestData.Name
	one.Roles = requestData.Roles
	one.Dept = requestData.Dept
	if len(one.Roles) < 0 {
		one.Admin = false
	} else {
		one.Admin = true
	}
	one.Save()
	return one.ID
}

func (r *RBAC) UsersList(ctx *gin.Context) gin.H {
	params := ctx.Request.URL.Query()
	page, err := strconv.Atoi(params.Get("page"))
	pageSize, err2 := strconv.Atoi(params.Get("pageSize"))
	search := params.Get("search")
	groupId := params.Get("groupId")
	if err != nil || err2 != nil {
		lib.CustomError(100, "page 或者PageSize 为非数字")
	}
	if page < 1 || pageSize < 1 {
		lib.CustomError(100, "page 或者PageSize 不能小于1")
	}
	m := models.Users{}
	users, total := m.FindAll(page, pageSize, search, groupId)
	if total == 0 {
		return gin.H{}
	}
	var rows []gin.H
	for _, user := range users {
		var roles models.RbacRoles
		var rolesList []string
		for _, RoleID := range user.Roles {
			if RoleID != "normal" && RoleID != "root" && RoleID != "0" {
				Id, err := strconv.Atoi(RoleID)
				if err != nil {
					lib.CustomError(100, "roleID 转换错误")
				}
				role, _ := roles.FindOne(Id)
				rolesList = append(rolesList, role.Label)
			}
		}
		var account models.Accounts
		acc, _ := account.FindOneByPlatform("password")
		if acc != nil {
			rows = append(rows, gin.H{
				"id":        user.ID,
				"userName":  acc.Username,
				"avatar":    user.Avatar,
				"mail":      user.Email,
				"name":      user.Name,
				"roles":     user.Roles,
				"dept":      user.Dept,
				"groupName": strings.Join(rolesList, ","),
				"date":      user.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}
	return gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"summary": gin.H{
			"id":   "20",
			"name": "999",
		},
		"rows": rows,
	}
}

// ==================Menu==================

func (r *RBAC) MenuCreate(ctx *gin.Context) uint {
	var requestData MenuJSONData
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		lib.CustomError(100, fmt.Sprintf("缺少必要参数Error:%v", err))
	}
	var Menu models.RbacMenu
	Menu.Parent = uint(requestData.ParentID)
	Menu.Name = requestData.Name
	Menu.Path = requestData.Path
	Menu.Component = requestData.Component
	Menu.CreatedAt = time.Now()
	Menu.Meta.Title = requestData.Meta.Title
	Menu.Meta.Type = requestData.Meta.Type
	err := Menu.Create()
	if err != nil {
		lib.CustomError(100, "创建失败")
	}
	return Menu.ID
}

func (r *RBAC) MenuEdit(ctx *gin.Context) uint {
	var requestData MenuJSONData
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}
	var Menu models.RbacMenu
	data := models.RbacMenu{
		ID:   uint(requestData.ID),
		Name: requestData.Name,
		Path: requestData.Path,
		Meta: models.RbacMenuMeta{
			Title:            requestData.Meta.Title,
			Type:             requestData.Meta.Type,
			Hidden:           requestData.Meta.Hidden,
			HiddenBreadcrumb: requestData.Meta.HiddenBreadcrumb,
			FullPage:         requestData.Meta.FullPage,
			Icon:             requestData.Meta.Icon,
			Color:            requestData.Meta.Color,
			Tag:              requestData.Meta.Tag,
		},
		Component: requestData.Component,
		Parent:    uint(requestData.ParentID),
		Sort:      requestData.Sort,
		Version:   "",
		Active:    requestData.Active,
	}
	err := Menu.Edit(data)
	if err != nil {
		lib.CustomError(101, "修改数据失败")
	}
	return data.ID
}

// 去重
func (r *RBAC) removeDuplicates(s []string) []string {
	uniqueMap := make(map[string]bool)
	var result []string

	for _, str := range s {
		if _, ok := uniqueMap[str]; !ok {
			uniqueMap[str] = true
			result = append(result, str)
		}
	}

	return result
}
func (r *RBAC) MenuMy(ctx *gin.Context) gin.H {
	UserInfoInterface, _ := ctx.Get("UserInfo")
	if UserInfo, ok := UserInfoInterface.(models.Users); ok {
		if UserInfo.Roles[0] == "root" {
			return gin.H{
				"menu":        r.MenuList(),
				"permissions": []string{},
				"dashboardGrid": []string{
					"welcome",
					"ver",
					"time",
					"progress",
					"echarts",
					"about",
				},
			}
		}
		fmt.Printf("UserInfoInterface:%v", UserInfo)
		rolesIdArray := UserInfo.Roles
		var menus []string
		for _, sRoleId := range rolesIdArray {
			if sRoleId != "normal" && sRoleId != "0" {
				iRoleId, err := strconv.Atoi(sRoleId)
				if err != nil {
					return nil
				}
				var role models.RbacRoles
				role, err = role.FindOne(iRoleId)
				if err != nil {
					continue
				}
				for _, menu := range role.Menus {
					menus = append(menus, menu)
				}
			}
		}
		output := r.removeDuplicates(menus)
		var menuModel models.RbacMenu
		var data []gin.H
		// 获取所有父级菜单
		MenusModel, _ := menuModel.FindMenusByNameAndParentID(output, 0)
		for _, menu := range MenusModel {
			menuData := r.buildMenu(menu, output)
			if menuData != nil {
				data = append(data, menuData)
			}
		}
		if len(data) == 0 {
			data = []gin.H{}
		}
		return gin.H{
			"menu":        data,
			"permissions": []string{},
			"dashboardGrid": []string{
				"welcome",
				"ver",
				"time",
				"progress",
				"echarts",
				"about",
			},
		}
	} else {
		lib.CustomError(100, "获取用户信息失败")
	}
	return gin.H{}
}
func (r *RBAC) MenuList() []gin.H {
	var data []gin.H
	var Routes models.RbacMenu
	Menus, _ := Routes.FindParentAll(0)
	for _, menu := range Menus {
		menuData := r.buildMenu(menu, []string{})
		if menuData != nil {
			data = append(data, menuData)
		}
	}
	return data
}

func (r *RBAC) buildMenu(menu models.RbacMenu, names []string) gin.H {
	meta := gin.H{
		"title":            menu.Meta.Title,
		"type":             menu.Meta.Type,
		"hidden":           menu.Meta.Hidden,
		"hiddenBreadcrumb": menu.Meta.HiddenBreadcrumb,
		"fullPage":         menu.Meta.FullPage,
		"icon":             menu.Meta.Icon,
		"color":            menu.Meta.Color,
		"tag":              menu.Meta.Tag,
		"affix":            menu.Meta.Affix,
	}
	if len(names) > 0 && !menu.IsStringInSortedSlice(names, menu.Name) {
		return nil
	}
	menuData := gin.H{
		"id":        menu.ID,
		"name":      menu.Name,
		"path":      menu.Path,
		"meta":      meta,
		"component": menu.Component,
		"parent":    menu.Parent,
		"sort":      menu.Sort,
		"version":   menu.Version,
		"active":    menu.Active,
		"children":  []gin.H{},
	}
	childrens, _ := menu.FindParentAll(menu.ID)
	for _, child := range childrens {
		childData := r.buildMenu(child, names)
		menuData["children"] = append(menuData["children"].([]gin.H), childData)
	}
	return menuData
}

func (r *RBAC) MenuDelete(ctx *gin.Context) bool {
	var data MenuDeleteJsonData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}
	var Menu models.RbacMenu
	for _, id := range data.Ids {
		err := Menu.Delete(uint(id))
		if err != nil {
			lib.CustomError(101, "删除失败")
		}
	}
	return true
}

// Dept Begin ==============================

// DeptCreate 组创建
func (r *RBAC) DeptCreate(ctx *gin.Context) gin.H {
	var requestData groupRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}
	var group models.RbacDept
	create, err := group.Create(requestData.Label, requestData.ParentId, requestData.Remark, requestData.Status, requestData.Sort)
	if err != nil {
		lib.CustomError(101, "创建失败")
	}
	return gin.H{
		"id":       create.ID,
		"label":    create.Label,
		"parentId": create.ParentId,
		"Remark":   create.Remarks,
		"Sort":     create.Sort,
		"Status":   create.Status,
		"data":     create.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (r *RBAC) DeptEdit(ctx *gin.Context) gin.H {
	var requestData groupRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}
	var group models.RbacDept
	group.Edit(requestData.Id, requestData.Label, requestData.ParentId, requestData.Remark, requestData.Status, requestData.Sort)
	return gin.H{
		"id":       group.ID,
		"label":    group.Label,
		"parentId": group.ParentId,
		"remark":   group.Remarks,
		"sort":     group.Sort,
		"status":   group.Status,
		"data":     group.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
func (r *RBAC) buildDept(sonDept models.RbacDept) gin.H {
	responseChild := gin.H{
		"id":       fmt.Sprintf("%d", sonDept.ID),
		"parentId": sonDept.ParentId,
		"label":    sonDept.Label,
		"date":     sonDept.CreatedAt.Format("2006-01-02 15:04:05"),
		"remark":   sonDept.Remarks,
		"status":   sonDept.Status,
		"sort":     sonDept.Sort,
		"children": []gin.H{},
	}

	childrens := sonDept.FindAll(sonDept.ID)
	for _, child := range childrens {
		childData := r.buildDept(child)
		responseChild["children"] = append(responseChild["children"].([]gin.H), childData)
	}
	return responseChild
}
func (r *RBAC) DeptList() []gin.H {
	var data []gin.H
	dept := models.RbacDept{}
	depths := dept.FindAll(0)
	for _, Dept := range depths {
		deptData := r.buildDept(Dept)
		data = append(data, deptData)
	}
	if len(data) == 0 || data == nil {
		return []gin.H{}
	}
	return data
}
func (r *RBAC) DeptDelete(ctx *gin.Context) gin.H {
	var request deleteJson
	if err := ctx.ShouldBindJSON(&request); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}
	var dept models.RbacDept
	if request.ID > 0 {
		err := dept.Delete(uint(request.ID))
		if err != nil {
			lib.CustomError(102, fmt.Sprintf("%s", err))
		}
	}
	if request.Ids != nil && len(request.Ids) > 0 {
		for _, id := range request.Ids {
			err := dept.Delete(uint(id))
			if err != nil {
				lib.CustomError(101, "删除失败")
			}
		}
	}
	return gin.H{
		"id":  request.ID,
		"ids": request.Ids,
	}
}

// Group End ==============================

// roles Begin ============================

type rbacRoles struct {
	Id     int      `json:"id"`
	Label  string   `json:"label"`
	Alias  string   `json:"alias"`
	Remark string   `json:"remark"`
	Sort   int      `json:"sort"`
	Status string   `json:"status"`
	Menus  []string `json:"menus"`
}

func (r *RBAC) RolesCreate(ctx *gin.Context) gin.H {
	var request rbacRoles
	if err := ctx.ShouldBindJSON(&request); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}
	var roles models.RbacRoles
	Status, err := strconv.Atoi(request.Status)
	if err != nil && request.Id != 0 {
		lib.CustomError(100, "Status不是一个字符串数字")
	}
	create, err := roles.Create(request.Label, request.Alias, request.Remark, Status, request.Sort)
	if err != nil {
		return nil
	}
	return gin.H{
		"id":     create.ID,
		"label":  create.Label,
		"alias":  create.Alias,
		"remark": create.Remark,
		"status": create.Status,
		"sort":   create.Sort,
		"data":   create.CreatedAt,
	}
}

func (r *RBAC) RolesEdit(ctx *gin.Context) gin.H {
	var request rbacRoles
	if err := ctx.ShouldBindJSON(&request); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}

	Status, err := strconv.Atoi(request.Status)
	if err != nil && request.Id != 0 {
		lib.CustomError(100, "Status不是一个字符串数字")
	}
	var roles models.RbacRoles
	if roles.Edit(request.Id, request.Label, request.Alias, request.Remark, Status, request.Sort, request.Menus) == false {
		lib.CustomError(101, "修改失败")
	}
	return gin.H{
		"id":     roles.ID,
		"label":  roles.Label,
		"alias":  roles.Alias,
		"remark": roles.Remark,
		"status": roles.Status,
		"sort":   roles.Sort,
		"data":   roles.CreatedAt,
	}
}

func (r *RBAC) RolesAllList(ctx *gin.Context) []gin.H {
	var roles models.RbacRoles
	rows, _ := roles.FindAll(1, 100)
	return rows
}

func (r *RBAC) RolesList(ctx *gin.Context) gin.H {
	params := ctx.Request.URL.Query()
	page, err := strconv.Atoi(params.Get("page"))
	pageSize, err2 := strconv.Atoi(params.Get("pageSize"))
	if err != nil || err2 != nil {
		lib.CustomError(100, "page 或者PageSize 为非数字")
	}

	if page < 1 || pageSize < 1 {
		lib.CustomError(100, "page 或者PageSize 不能小于1")
	}
	var roles models.RbacRoles
	rows, total := roles.FindAll(page, pageSize)
	return models.ConversionListData(total, page, pageSize, rows)
}

type deleteJson struct {
	ID  int   `json:"id"`
	Ids []int `json:"ids"`
}

func (r *RBAC) RolesDelete(ctx *gin.Context) gin.H {
	var request deleteJson
	if err := ctx.ShouldBindJSON(&request); err != nil {
		lib.CustomError(100, "缺少必要参数")
	}
	var roles models.RbacRoles
	if request.ID > 0 {
		err := roles.Delete(uint(request.ID))
		if err != nil {
			lib.CustomError(102, fmt.Sprintf("%s", err))
		}
	}
	if request.Ids != nil && len(request.Ids) > 0 {
		for _, id := range request.Ids {
			err := roles.Delete(uint(id))
			if err != nil {
				lib.CustomError(101, "删除失败")
			}
		}
	}
	return gin.H{
		"id":  request.ID,
		"ids": request.Ids,
	}
}

// roles End ============================
