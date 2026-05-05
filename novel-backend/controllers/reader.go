package controllers

import (
	"novel-backend/models"
	"novel-backend/pkg/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBookshelf(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var bookshelf []models.Bookshelf
	database.DB.Where("user_id = ?", userID.(uint)).Preload("Novel").Preload("Novel.Category").Preload("Novel.Author").Order("updated_at DESC").Find(&bookshelf)

	Success(c, gin.H{
		"list": bookshelf,
	})
}

func AddToBookshelf(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var req struct {
		NovelID uint `json:"novel_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 检查小说是否存在
	var novel models.Novel
	if err := database.DB.First(&novel, req.NovelID).Error; err != nil {
		NotFound(c, "小说不存在")
		return
	}

	// 检查是否已加入书架
	var existingBookshelf models.Bookshelf
	if err := database.DB.Where("user_id = ? AND novel_id = ?", userID.(uint), req.NovelID).First(&existingBookshelf).Error; err == nil {
		BadRequest(c, "已加入书架")
		return
	}

	bookshelf := models.Bookshelf{
		UserID:  userID.(uint),
		NovelID: req.NovelID,
	}

	if err := database.DB.Create(&bookshelf).Error; err != nil {
		InternalServerError(c, "加入书架失败")
		return
	}

	// 增加收藏数
	database.DB.Model(&novel).Update("collect_count", novel.CollectCount+1)

	Success(c, gin.H{
		"message": "加入书架成功",
	})
}

func RemoveFromBookshelf(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	bookshelfID := c.Param("id")

	var bookshelf models.Bookshelf
	if err := database.DB.Where("id = ? AND user_id = ?", bookshelfID, userID.(uint)).First(&bookshelf).Error; err != nil {
		NotFound(c, "记录不存在")
		return
	}

	// 减少收藏数
	database.DB.Model(&models.Novel{}).Where("id = ?", bookshelf.NovelID).Update("collect_count", gorm.Expr("collect_count - 1"))

	if err := database.DB.Delete(&bookshelf).Error; err != nil {
		InternalServerError(c, "移除失败")
		return
	}

	Success(c, gin.H{
		"message": "移除成功",
	})
}

func CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var req struct {
		NovelID  uint   `json:"novel_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 检查小说是否存在
	var novel models.Novel
	if err := database.DB.First(&novel, req.NovelID).Error; err != nil {
		NotFound(c, "小说不存在")
		return
	}

	comment := models.Comment{
		UserID:   userID.(uint),
		NovelID:  req.NovelID,
		Content:  req.Content,
		ParentID: req.ParentID,
		Status:   1,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		InternalServerError(c, "评论失败")
		return
	}

	// 增加评论数
	database.DB.Model(&novel).Update("comment_count", novel.CommentCount+1)

	Success(c, gin.H{
		"message": "评论成功",
		"comment": comment,
	})
}

func Recharge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	if req.Amount <= 0 {
		BadRequest(c, "充值金额必须大于0")
		return
	}

	// 创建订单
	orderNo := generateOrderNo()
	order := models.Order{
		OrderNo:   orderNo,
		UserID:    userID.(uint),
		OrderType: 1, // 充值
		Amount:    req.Amount,
		Status:    0, // 待支付
	}

	if err := database.DB.Create(&order).Error; err != nil {
		InternalServerError(c, "创建订单失败")
		return
	}

	// TODO: 这里应该调用支付接口，这里简化处理，直接模拟支付成功
	// 实际应用中应该等待支付回调

	Success(c, gin.H{
		"order_no": orderNo,
		"amount":   req.Amount,
		"message":  "订单创建成功，请完成支付",
	})
}

func Subscribe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var req struct {
		NovelID uint `json:"novel_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 检查小说是否存在
	var novel models.Novel
	if err := database.DB.First(&novel, req.NovelID).Error; err != nil {
		NotFound(c, "小说不存在")
		return
	}

	// 检查是否已订阅
	var existingSub models.Subscription
	if err := database.DB.Where("user_id = ? AND novel_id = ?", userID.(uint), req.NovelID).First(&existingSub).Error; err == nil {
		BadRequest(c, "已订阅该小说")
		return
	}

	// TODO: 检查用户余额是否足够，这里简化处理
	// 实际应用中应该扣除用户余额

	subscription := models.Subscription{
		UserID:  userID.(uint),
		NovelID: req.NovelID,
		Status:  1,
	}

	if err := database.DB.Create(&subscription).Error; err != nil {
		InternalServerError(c, "订阅失败")
		return
	}

	Success(c, gin.H{
		"message": "订阅成功",
	})
}

func GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := database.DB.Model(&models.Order{}).Where("user_id = ?", userID.(uint))
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	Success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func generateOrderNo() string {
	return "ORD" + time.Now().Format("20060102150405") + strconv.Itoa(int(time.Now().UnixNano()%10000))
}
