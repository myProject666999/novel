package controllers

import (
	"novel-backend/models"
	"novel-backend/pkg/database"
	"novel-backend/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	InviteCode string `json:"invite_code"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.Preload("Role").Where("username = ?", req.Username).First(&user).Error; err != nil {
		Unauthorized(c, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		Forbidden(c, "账号已被禁用")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		InternalServerError(c, "生成令牌失败")
		return
	}

	Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"nickname":  user.Nickname,
			"avatar":    user.Avatar,
			"role_id":   user.RoleID,
			"role_name": user.Role.Name,
			"balance":   user.Balance,
			"vip_level": user.VipLevel,
		},
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		BadRequest(c, "用户名已存在")
		return
	}

	// 检查邀请码（如果是作家注册可能需要）
	var inviteCode *models.InviteCode
	if req.InviteCode != "" {
		inviteCode = &models.InviteCode{}
		if err := database.DB.Where("code = ? AND status = 1 AND used_count < max_count", req.InviteCode).First(inviteCode).Error; err != nil {
			BadRequest(c, "邀请码无效或已使用")
			return
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		InternalServerError(c, "密码加密失败")
		return
	}

	// 默认角色是读者（假设角色ID 3是读者）
	roleID := uint(3)

	// 如果有邀请码，可能是作家注册
	if inviteCode != nil {
		roleID = uint(2) // 作家角色
		inviteCode.UsedCount++
		database.DB.Save(inviteCode)
	}

	user := models.User{
		Username:     req.Username,
		Password:     hashedPassword,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       1,
		RoleID:       roleID,
		InviteCodeID: func() *uint { if inviteCode != nil { return &inviteCode.ID }; return nil }(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		InternalServerError(c, "注册失败")
		return
	}

	Success(c, gin.H{
		"message": "注册成功",
		"user_id": user.ID,
	})
}

func GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var user models.User
	if err := database.DB.Preload("Role").First(&user, userID.(uint)).Error; err != nil {
		NotFound(c, "用户不存在")
		return
	}

	Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"phone":      user.Phone,
		"avatar":     user.Avatar,
		"balance":    user.Balance,
		"vip_level":  user.VipLevel,
		"role_id":    user.RoleID,
		"role_name":  user.Role.Name,
		"created_at": user.CreatedAt,
	})
}

func UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID.(uint)).Error; err != nil {
		NotFound(c, "用户不存在")
		return
	}

	// 只更新允许修改的字段
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID.(uint)).Error; err != nil {
		NotFound(c, "用户不存在")
		return
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		BadRequest(c, "原密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		InternalServerError(c, "密码加密失败")
		return
	}

	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		InternalServerError(c, "修改密码失败")
		return
	}

	Success(c, gin.H{
		"message": "密码修改成功",
	})
}
