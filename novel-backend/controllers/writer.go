package controllers

import (
	"novel-backend/models"
	"novel-backend/pkg/database"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetWriterNovels(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := database.DB.Model(&models.Novel{}).Where("author_id = ?", userID.(uint)).Preload("Category")

	if status != "" {
		query = query.Where("status = ?", status)
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

func CreateNovel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var req struct {
		Title       string `json:"title" binding:"required"`
		CategoryID  uint   `json:"category_id" binding:"required"`
		Cover       string `json:"cover"`
		Description string `json:"description"`
		VIP         int    `json:"vip"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 检查分类是否存在
	var category models.Category
	if err := database.DB.First(&category, req.CategoryID).Error; err != nil {
		NotFound(c, "分类不存在")
		return
	}

	novel := models.Novel{
		Title:       req.Title,
		AuthorID:    userID.(uint),
		CategoryID:  req.CategoryID,
		Cover:       req.Cover,
		Description: req.Description,
		Status:      1, // 默认连载中
		VIP:         req.VIP,
	}

	if err := database.DB.Create(&novel).Error; err != nil {
		InternalServerError(c, "创建小说失败")
		return
	}

	Success(c, gin.H{
		"message": "创建成功",
		"novel":   novel,
	})
}

func UpdateNovel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	novelID := c.Param("id")

	var novel models.Novel
	if err := database.DB.Where("id = ? AND author_id = ?", novelID, userID.(uint)).First(&novel).Error; err != nil {
		NotFound(c, "小说不存在或无权修改")
		return
	}

	var req struct {
		Title       string `json:"title"`
		CategoryID  uint   `json:"category_id"`
		Cover       string `json:"cover"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
		VIP         *int   `json:"vip"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.CategoryID != 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.VIP != nil {
		updates["vip"] = *req.VIP
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&novel).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteNovel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	novelID := c.Param("id")

	var novel models.Novel
	if err := database.DB.Where("id = ? AND author_id = ?", novelID, userID.(uint)).First(&novel).Error; err != nil {
		NotFound(c, "小说不存在或无权删除")
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

func CreateChapter(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	novelID := c.Param("novel_id")

	// 检查小说是否属于当前作家
	var novel models.Novel
	if err := database.DB.Where("id = ? AND author_id = ?", novelID, userID.(uint)).First(&novel).Error; err != nil {
		NotFound(c, "小说不存在或无权操作")
		return
	}

	var req struct {
		Title   string  `json:"title" binding:"required"`
		Content string  `json:"content" binding:"required"`
		VIP     int     `json:"vip"`
		Price   float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 获取当前最大章节号
	var maxChapterNum int
	database.DB.Model(&models.Chapter{}).Where("novel_id = ?", novelID).Select("COALESCE(MAX(chapter_num), 0)").Scan(&maxChapterNum)

	chapter := models.Chapter{
		NovelID:     novel.ID,
		Title:       req.Title,
		Content:     req.Content,
		WordCount:   len(strings.ReplaceAll(req.Content, "\n", "")),
		ChapterNum:  maxChapterNum + 1,
		VIP:         req.VIP,
		Price:       req.Price,
		Status:      1,
	}

	if err := database.DB.Create(&chapter).Error; err != nil {
		InternalServerError(c, "创建章节失败")
		return
	}

	// 更新小说的字数
	database.DB.Model(&novel).Update("word_count", novel.WordCount+chapter.WordCount)

	Success(c, gin.H{
		"message":  "创建成功",
		"chapter":  chapter,
	})
}

func UpdateChapter(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	chapterID := c.Param("id")

	// 检查章节是否属于当前作家
	var chapter models.Chapter
	if err := database.DB.Preload("Novel").First(&chapter, chapterID).Error; err != nil {
		NotFound(c, "章节不存在")
		return
	}

	if chapter.Novel.AuthorID != userID.(uint) {
		Forbidden(c, "无权修改此章节")
		return
	}

	var req struct {
		Title   string  `json:"title"`
		Content string  `json:"content"`
		VIP     *int    `json:"vip"`
		Price   *float64 `json:"price"`
		Status  *int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	oldWordCount := chapter.WordCount

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
		newWordCount := len(strings.ReplaceAll(req.Content, "\n", ""))
		updates["word_count"] = newWordCount
		// 更新小说字数
		database.DB.Model(&models.Novel{}).Where("id = ?", chapter.NovelID).Update("word_count", gorm.Expr("word_count - ? + ?", oldWordCount, newWordCount))
	}
	if req.VIP != nil {
		updates["vip"] = *req.VIP
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&chapter).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteChapter(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	chapterID := c.Param("id")

	var chapter models.Chapter
	if err := database.DB.Preload("Novel").First(&chapter, chapterID).Error; err != nil {
		NotFound(c, "章节不存在")
		return
	}

	if chapter.Novel.AuthorID != userID.(uint) {
		Forbidden(c, "无权删除此章节")
		return
	}

	// 更新小说字数
	database.DB.Model(&models.Novel{}).Where("id = ?", chapter.NovelID).Update("word_count", gorm.Expr("word_count - ?", chapter.WordCount))

	// 软删除
	if err := database.DB.Delete(&chapter).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}
