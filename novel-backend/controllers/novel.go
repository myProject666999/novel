package controllers

import (
	"novel-backend/models"
	"novel-backend/pkg/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetNovelList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID := c.Query("category_id")
	status := c.Query("status")

	query := database.DB.Model(&models.Novel{}).Preload("Category").Preload("Author")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var novels []models.Novel
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&novels)

	Success(c, gin.H{
		"list":       novels,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

func GetNovelDetail(c *gin.Context) {
	id := c.Param("id")

	var novel models.Novel
	if err := database.DB.Preload("Category").Preload("Author").First(&novel, id).Error; err != nil {
		NotFound(c, "小说不存在")
		return
	}

	// 增加点击量
	database.DB.Model(&novel).Update("click_count", novel.ClickCount+1)

	Success(c, novel)
}

func GetNovelRank(c *gin.Context) {
	rankType := c.DefaultQuery("type", "click")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var query *gorm.DB
	switch rankType {
	case "click":
		query = database.DB.Model(&models.Novel{}).Order("click_count DESC")
	case "collect":
		query = database.DB.Model(&models.Novel{}).Order("collect_count DESC")
	case "recommend":
		query = database.DB.Model(&models.Novel{}).Where("recommend = 1").Order("updated_at DESC")
	default:
		query = database.DB.Model(&models.Novel{}).Order("click_count DESC")
	}

	var novels []models.Novel
	query.Preload("Category").Preload("Author").Limit(limit).Find(&novels)

	Success(c, gin.H{
		"type": rankType,
		"list": novels,
	})
}

func GetRecommendNovels(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var recommends []models.NovelRecommend
	database.DB.Where("status = 1").Order("sort ASC").Preload("Novel").Preload("Novel.Category").Preload("Novel.Author").Limit(limit).Find(&recommends)

	Success(c, gin.H{
		"list": recommends,
	})
}

func SearchNovels(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if keyword == "" {
		BadRequest(c, "请输入搜索关键词")
		return
	}

	query := database.DB.Model(&models.Novel{}).Preload("Category").Preload("Author")
	query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	var total int64
	query.Count(&total)

	var novels []models.Novel
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("click_count DESC").Find(&novels)

	Success(c, gin.H{
		"keyword":   keyword,
		"list":      novels,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetChapterDetail(c *gin.Context) {
	chapterID := c.Param("id")

	var chapter models.Chapter
	if err := database.DB.Preload("Novel").First(&chapter, chapterID).Error; err != nil {
		NotFound(c, "章节不存在")
		return
	}

	// 检查是否是VIP章节
	if chapter.VIP == 1 {
		// TODO: 检查用户是否已购买或是否是会员
		// userID, exists := c.Get("user_id")
		// if !exists {
		//     Unauthorized(c, "请先登录")
		//     return
		// }
		// 检查用户是否有权限阅读
	}

	Success(c, chapter)
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Where("status = 1").Order("sort ASC").Find(&categories)

	Success(c, gin.H{
		"list": categories,
	})
}

func GetNovelChapters(c *gin.Context) {
	novelID := c.Param("novel_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	var novel models.Novel
	if err := database.DB.First(&novel, novelID).Error; err != nil {
		NotFound(c, "小说不存在")
		return
	}

	var total int64
	database.DB.Model(&models.Chapter{}).Where("novel_id = ?", novelID).Count(&total)

	var chapters []models.Chapter
	offset := (page - 1) * pageSize
	database.DB.Where("novel_id = ?", novelID).Order("chapter_num ASC").Offset(offset).Limit(pageSize).Find(&chapters)

	Success(c, gin.H{
		"novel_id":  novelID,
		"list":      chapters,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetNovelComments(c *gin.Context) {
	novelID := c.Param("novel_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	database.DB.Model(&models.Comment{}).Where("novel_id = ? AND status = 1", novelID).Count(&total)

	var comments []models.Comment
	offset := (page - 1) * pageSize
	database.DB.Where("novel_id = ? AND status = 1", novelID).Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments)

	Success(c, gin.H{
		"novel_id":  novelID,
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
