package controllers

import (
	"novel-backend/models"
	"novel-backend/pkg/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户管理
func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")
	status := c.Query("status")
	roleID := c.Query("role_id")

	query := database.DB.Model(&models.User{}).Preload("Role")

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users)

	Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		NotFound(c, "用户不存在")
		return
	}

	user.Status = req.Status
	if err := database.DB.Save(&user).Error; err != nil {
		InternalServerError(c, "更新失败")
		return
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		NotFound(c, "用户不存在")
		return
	}

	// 软删除
	if err := database.DB.Delete(&user).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}

// 角色管理
func GetRoleList(c *gin.Context) {
	var roles []models.Role
	database.DB.Preload("Menus").Order("id ASC").Find(&roles)

	Success(c, gin.H{
		"list": roles,
	})
}

func CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
		Description string `json:"description"`
		MenuIDs     []uint `json:"menu_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 检查角色名是否已存在
	var existingRole models.Role
	if err := database.DB.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		BadRequest(c, "角色名已存在")
		return
	}

	role := models.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Status:      1,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		InternalServerError(c, "创建角色失败")
		return
	}

	// 关联菜单
	if len(req.MenuIDs) > 0 {
		var menus []models.Menu
		database.DB.Find(&menus, req.MenuIDs)
		database.DB.Model(&role).Association("Menus").Replace(menus)
	}

	Success(c, gin.H{
		"message": "创建成功",
		"role":    role,
	})
}

func UpdateRole(c *gin.Context) {
	roleID := c.Param("id")

	var role models.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		NotFound(c, "角色不存在")
		return
	}

	var req struct {
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
		MenuIDs     []uint `json:"menu_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&role).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	// 更新菜单关联
	if req.MenuIDs != nil {
		var menus []models.Menu
		database.DB.Find(&menus, req.MenuIDs)
		database.DB.Model(&role).Association("Menus").Replace(menus)
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteRole(c *gin.Context) {
	roleID := c.Param("id")

	var role models.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		NotFound(c, "角色不存在")
		return
	}

	// 检查是否有用户使用此角色
	var userCount int64
	database.DB.Model(&models.User{}).Where("role_id = ?", roleID).Count(&userCount)
	if userCount > 0 {
		BadRequest(c, "该角色下还有用户，无法删除")
		return
	}

	// 软删除
	if err := database.DB.Delete(&role).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}

// 菜单管理
func GetMenuList(c *gin.Context) {
	var menus []models.Menu
	database.DB.Order("sort ASC, id ASC").Find(&menus)

	Success(c, gin.H{
		"list": menus,
	})
}

func CreateMenu(c *gin.Context) {
	var req struct {
		ParentID   *uint  `json:"parent_id"`
		Name       string `json:"name" binding:"required"`
		Path       string `json:"path"`
		Icon       string `json:"icon"`
		Sort       int    `json:"sort"`
		MenuType   int    `json:"menu_type"`
		Permission string `json:"permission"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	menu := models.Menu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Icon:       req.Icon,
		Sort:       req.Sort,
		MenuType:   req.MenuType,
		Permission: req.Permission,
		Status:     1,
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		InternalServerError(c, "创建菜单失败")
		return
	}

	Success(c, gin.H{
		"message": "创建成功",
		"menu":    menu,
	})
}

func UpdateMenu(c *gin.Context) {
	menuID := c.Param("id")

	var menu models.Menu
	if err := database.DB.First(&menu, menuID).Error; err != nil {
		NotFound(c, "菜单不存在")
		return
	}

	var req struct {
		ParentID   *uint  `json:"parent_id"`
		Name       string `json:"name"`
		Path       string `json:"path"`
		Icon       string `json:"icon"`
		Sort       *int   `json:"sort"`
		MenuType   *int   `json:"menu_type"`
		Permission string `json:"permission"`
		Status     *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Path != "" {
		updates["path"] = req.Path
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.MenuType != nil {
		updates["menu_type"] = *req.MenuType
	}
	if req.Permission != "" {
		updates["permission"] = req.Permission
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&menu).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteMenu(c *gin.Context) {
	menuID := c.Param("id")

	var menu models.Menu
	if err := database.DB.First(&menu, menuID).Error; err != nil {
		NotFound(c, "菜单不存在")
		return
	}

	// 检查是否有子菜单
	var childCount int64
	database.DB.Model(&models.Menu{}).Where("parent_id = ?", menuID).Count(&childCount)
	if childCount > 0 {
		BadRequest(c, "该菜单下还有子菜单，无法删除")
		return
	}

	// 软删除
	if err := database.DB.Delete(&menu).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}
