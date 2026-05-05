package controllers

import (
	"novel-backend/models"
	"novel-backend/pkg/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 小说管理（管理员）
func GetAdminNovelList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	title := c.Query("title")
	status := c.Query("status")
	categoryID := c.Query("category_id")

	query := database.DB.Model(&models.Novel{}).Preload("Category").Preload("Author")

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	query.Count(&total)

	var novels []models.Novel
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&novels)

	Success(c, gin.H{
		"list":      novels,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateNovelStatus(c *gin.Context) {
	novelID := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var novel models.Novel
	if err := database.DB.First(&novel, novelID).Error; err != nil {
		NotFound(c, "小说不存在")
		return
	}

	novel.Status = req.Status
	if err := database.DB.Save(&novel).Error; err != nil {
		InternalServerError(c, "更新失败")
		return
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteNovelAdmin(c *gin.Context) {
	novelID := c.Param("id")

	var novel models.Novel
	if err := database.DB.First(&novel, novelID).Error; err != nil {
		NotFound(c, "小说不存在")
		return
	}

	// 软删除
	if err := database.DB.Delete(&novel).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}

// 评论管理
func GetCommentList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := database.DB.Model(&models.Comment{}).Preload("User").Preload("Novel")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var comments []models.Comment
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments)

	Success(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func DeleteCommentAdmin(c *gin.Context) {
	commentID := c.Param("id")

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		NotFound(c, "评论不存在")
		return
	}

	// 减少小说评论数
	database.DB.Model(&models.Novel{}).Where("id = ?", comment.NovelID).Update("comment_count", gorm.Expr("comment_count - 1"))

	// 软删除
	if err := database.DB.Delete(&comment).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}

// 作家管理
func GetWriterList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")
	status := c.Query("status")

	// 假设角色ID 2是作家
	query := database.DB.Model(&models.User{}).Where("role_id = 2").Preload("Role")

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var writers []models.User
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&writers)

	Success(c, gin.H{
		"list":      writers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateWriterStatus(c *gin.Context) {
	writerID := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, writerID).Error; err != nil {
		NotFound(c, "作家不存在")
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

// 邀请码管理
func GetInviteCodeList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	code := c.Query("code")
	status := c.Query("status")

	query := database.DB.Model(&models.InviteCode{})

	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var inviteCodes []models.InviteCode
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&inviteCodes)

	Success(c, gin.H{
		"list":      inviteCodes,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateInviteCode(c *gin.Context) {
	var req struct {
		Code     string `json:"code"`
		MaxCount int    `json:"max_count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 如果没有提供code，自动生成
	if req.Code == "" {
		req.Code = generateInviteCode()
	}

	// 检查邀请码是否已存在
	var existingCode models.InviteCode
	if err := database.DB.Where("code = ?", req.Code).First(&existingCode).Error; err == nil {
		BadRequest(c, "邀请码已存在")
		return
	}

	if req.MaxCount <= 0 {
		req.MaxCount = 1
	}

	inviteCode := models.InviteCode{
		Code:      req.Code,
		UsedCount: 0,
		MaxCount:  req.MaxCount,
		Status:    1,
	}

	if err := database.DB.Create(&inviteCode).Error; err != nil {
		InternalServerError(c, "创建邀请码失败")
		return
	}

	Success(c, gin.H{
		"message": "创建成功",
		"code":    inviteCode,
	})
}

func DeleteInviteCode(c *gin.Context) {
	codeID := c.Param("id")

	var inviteCode models.InviteCode
	if err := database.DB.First(&inviteCode, codeID).Error; err != nil {
		NotFound(c, "邀请码不存在")
		return
	}

	// 软删除
	if err := database.DB.Delete(&inviteCode).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}

// 订单管理
func GetOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	orderNo := c.Query("order_no")
	status := c.Query("status")
	orderType := c.Query("order_type")

	query := database.DB.Model(&models.Order{}).Preload("User")

	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders)

	Success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		NotFound(c, "订单不存在")
		return
	}

	order.Status = req.Status
	if req.Status == 1 { // 已支付
		now := time.Now()
		order.PayTime = &now
	}

	if err := database.DB.Save(&order).Error; err != nil {
		InternalServerError(c, "更新失败")
		return
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

// 系统日志管理
func GetLogList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	module := c.Query("module")
	username := c.Query("username")

	query := database.DB.Model(&models.SystemLog{})

	if module != "" {
		query = query.Where("module = ?", module)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	var total int64
	query.Count(&total)

	var logs []models.SystemLog
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs)

	Success(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// 辅助函数
func generateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	now := time.Now().UnixNano()
	for i := range code {
		code[i] = charset[now%int64(len(charset))]
		now = now / int64(len(charset))
	}
	return string(code)
}
