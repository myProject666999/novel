package controllers

import (
	"novel-backend/models"
	"novel-backend/pkg/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 会员管理
func GetMemberList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	level := c.Query("level")

	query := database.DB.Model(&models.Member{}).Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}

	var total int64
	query.Count(&total)

	var members []models.Member
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&members)

	Success(c, gin.H{
		"list":      members,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateMember(c *gin.Context) {
	memberID := c.Param("id")

	var req struct {
		Level    *int `json:"level"`
		Status   *int `json:"status"`
		EndDate  *string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var member models.Member
	if err := database.DB.First(&member, memberID).Error; err != nil {
		NotFound(c, "会员不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Level != nil {
		updates["level"] = *req.Level
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.EndDate != nil {
		// 解析日期
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			BadRequest(c, "日期格式错误")
			return
		}
		updates["end_date"] = endDate
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&member).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

// 会员反馈管理
func GetFeedbackList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := database.DB.Model(&models.Feedback{}).Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var feedbacks []models.Feedback
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&feedbacks)

	Success(c, gin.H{
		"list":      feedbacks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func HandleFeedback(c *gin.Context) {
	feedbackID := c.Param("id")

	var req struct {
		Reply string `json:"reply" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var feedback models.Feedback
	if err := database.DB.First(&feedback, feedbackID).Error; err != nil {
		NotFound(c, "反馈不存在")
		return
	}

	now := time.Now()
	feedback.Reply = req.Reply
	feedback.Status = 1
	feedback.ReplyTime = &now

	if err := database.DB.Save(&feedback).Error; err != nil {
		InternalServerError(c, "处理失败")
		return
	}

	Success(c, gin.H{
		"message": "处理成功",
	})
}

// 小说推荐管理
func GetRecommendList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	position := c.Query("position")
	status := c.Query("status")

	query := database.DB.Model(&models.NovelRecommend{}).Preload("Novel").Preload("Novel.Author").Preload("Novel.Category")

	if position != "" {
		query = query.Where("position = ?", position)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var recommends []models.NovelRecommend
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("sort ASC, created_at DESC").Find(&recommends)

	Success(c, gin.H{
		"list":      recommends,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateRecommend(c *gin.Context) {
	var req struct {
		NovelID  uint `json:"novel_id" binding:"required"`
		Position int  `json:"position"`
		Sort     int  `json:"sort"`
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

	// 检查是否已推荐
	var existing models.NovelRecommend
	if err := database.DB.Where("novel_id = ?", req.NovelID).First(&existing).Error; err == nil {
		BadRequest(c, "该小说已在推荐列表中")
		return
	}

	recommend := models.NovelRecommend{
		NovelID:  req.NovelID,
		Position: req.Position,
		Sort:     req.Sort,
		Status:   1,
	}

	if err := database.DB.Create(&recommend).Error; err != nil {
		InternalServerError(c, "创建推荐失败")
		return
	}

	Success(c, gin.H{
		"message":  "创建成功",
		"recommend": recommend,
	})
}

func UpdateRecommend(c *gin.Context) {
	recommendID := c.Param("id")

	var req struct {
		NovelID  *uint `json:"novel_id"`
		Position *int  `json:"position"`
		Sort     *int  `json:"sort"`
		Status   *int  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var recommend models.NovelRecommend
	if err := database.DB.First(&recommend, recommendID).Error; err != nil {
		NotFound(c, "推荐不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.NovelID != nil {
		// 检查小说是否存在
		var novel models.Novel
		if err := database.DB.First(&novel, *req.NovelID).Error; err != nil {
			NotFound(c, "小说不存在")
			return
		}
		updates["novel_id"] = *req.NovelID
	}
	if req.Position != nil {
		updates["position"] = *req.Position
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&recommend).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteRecommend(c *gin.Context) {
	recommendID := c.Param("id")

	var recommend models.NovelRecommend
	if err := database.DB.First(&recommend, recommendID).Error; err != nil {
		NotFound(c, "推荐不存在")
		return
	}

	// 软删除
	if err := database.DB.Delete(&recommend).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}

// 网站信息管理
func GetSiteInfo(c *gin.Context) {
	var siteInfo models.SiteInfo
	// 如果没有记录，创建默认记录
	if err := database.DB.First(&siteInfo).Error; err != nil {
		// 创建默认网站信息
		siteInfo = models.SiteInfo{
			SiteName:        "小说阅读网",
			SiteLogo:        "",
			SiteKeywords:    "小说,阅读,网络小说",
			SiteDescription: "专业的小说阅读平台",
			Copyright:       "©2024 小说阅读网 版权所有",
			Icp:             "",
			ContactEmail:    "contact@novel.com",
			ContactPhone:    "",
		}
	}

	Success(c, siteInfo)
}

func UpdateSiteInfo(c *gin.Context) {
	var req struct {
		SiteName        string `json:"site_name"`
		SiteLogo        string `json:"site_logo"`
		SiteKeywords    string `json:"site_keywords"`
		SiteDescription string `json:"site_description"`
		Copyright       string `json:"copyright"`
		Icp             string `json:"icp"`
		ContactEmail    string `json:"contact_email"`
		ContactPhone    string `json:"contact_phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var siteInfo models.SiteInfo
	// 尝试获取现有记录，不存在则创建
	if err := database.DB.First(&siteInfo).Error; err != nil {
		// 创建新记录
		siteInfo = models.SiteInfo{
			SiteName:        req.SiteName,
			SiteLogo:        req.SiteLogo,
			SiteKeywords:    req.SiteKeywords,
			SiteDescription: req.SiteDescription,
			Copyright:       req.Copyright,
			Icp:             req.Icp,
			ContactEmail:    req.ContactEmail,
			ContactPhone:    req.ContactPhone,
		}
		if err := database.DB.Create(&siteInfo).Error; err != nil {
			InternalServerError(c, "创建失败")
			return
		}
	} else {
		// 更新现有记录
		updates := make(map[string]interface{})
		if req.SiteName != "" {
			updates["site_name"] = req.SiteName
		}
		if req.SiteLogo != "" {
			updates["site_logo"] = req.SiteLogo
		}
		if req.SiteKeywords != "" {
			updates["site_keywords"] = req.SiteKeywords
		}
		if req.SiteDescription != "" {
			updates["site_description"] = req.SiteDescription
		}
		if req.Copyright != "" {
			updates["copyright"] = req.Copyright
		}
		if req.Icp != "" {
			updates["icp"] = req.Icp
		}
		if req.ContactEmail != "" {
			updates["contact_email"] = req.ContactEmail
		}
		if req.ContactPhone != "" {
			updates["contact_phone"] = req.ContactPhone
		}

		if len(updates) > 0 {
			if err := database.DB.Model(&siteInfo).Updates(updates).Error; err != nil {
				InternalServerError(c, "更新失败")
				return
			}
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

// 友情链接管理
func GetLinkList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	query := database.DB.Model(&models.FriendLink{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var links []models.FriendLink
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("sort ASC, created_at DESC").Find(&links)

	Success(c, gin.H{
		"list":      links,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateLink(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Url    string `json:"url" binding:"required"`
		Logo   string `json:"logo"`
		Sort   int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	link := models.FriendLink{
		Name:   req.Name,
		Url:    req.Url,
		Logo:   req.Logo,
		Sort:   req.Sort,
		Status: 1,
	}

	if err := database.DB.Create(&link).Error; err != nil {
		InternalServerError(c, "创建友情链接失败")
		return
	}

	Success(c, gin.H{
		"message": "创建成功",
		"link":    link,
	})
}

func UpdateLink(c *gin.Context) {
	linkID := c.Param("id")

	var req struct {
		Name   *string `json:"name"`
		Url    *string `json:"url"`
		Logo   *string `json:"logo"`
		Sort   *int    `json:"sort"`
		Status *int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var link models.FriendLink
	if err := database.DB.First(&link, linkID).Error; err != nil {
		NotFound(c, "友情链接不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Url != nil {
		updates["url"] = *req.Url
	}
	if req.Logo != nil {
		updates["logo"] = *req.Logo
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&link).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteLink(c *gin.Context) {
	linkID := c.Param("id")

	var link models.FriendLink
	if err := database.DB.First(&link, linkID).Error; err != nil {
		NotFound(c, "友情链接不存在")
		return
	}

	// 软删除
	if err := database.DB.Delete(&link).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}

// 新闻管理
func GetNewsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := database.DB.Model(&models.News{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var news []models.News
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&news)

	Success(c, gin.H{
		"list":      news,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateNews(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Cover   string `json:"cover"`
		Author  string `json:"author"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	if req.Status == 0 {
		req.Status = 1 // 默认发布
	}

	news := models.News{
		Title:   req.Title,
		Content: req.Content,
		Cover:   req.Cover,
		Author:  req.Author,
		Status:  req.Status,
	}

	if err := database.DB.Create(&news).Error; err != nil {
		InternalServerError(c, "创建新闻失败")
		return
	}

	Success(c, gin.H{
		"message": "创建成功",
		"news":    news,
	})
}

func UpdateNews(c *gin.Context) {
	newsID := c.Param("id")

	var req struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Cover   *string `json:"cover"`
		Author  *string `json:"author"`
		Status  *int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var news models.News
	if err := database.DB.First(&news, newsID).Error; err != nil {
		NotFound(c, "新闻不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Cover != nil {
		updates["cover"] = *req.Cover
	}
	if req.Author != nil {
		updates["author"] = *req.Author
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&news).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteNews(c *gin.Context) {
	newsID := c.Param("id")

	var news models.News
	if err := database.DB.First(&news, newsID).Error; err != nil {
		NotFound(c, "新闻不存在")
		return
	}

	// 软删除
	if err := database.DB.Delete(&news).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}

// 类别管理
func GetCategoryList(c *gin.Context) {
	var categories []models.Category
	database.DB.Order("sort ASC, id ASC").Find(&categories)

	Success(c, gin.H{
		"list": categories,
	})
}

func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	// 检查分类名是否已存在
	var existing models.Category
	if err := database.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		BadRequest(c, "分类名已存在")
		return
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      1,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		InternalServerError(c, "创建分类失败")
		return
	}

	Success(c, gin.H{
		"message":  "创建成功",
		"category": category,
	})
}

func UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Sort        *int    `json:"sort"`
		Status      *int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		NotFound(c, "分类不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		// 检查分类名是否已被其他分类使用
		var existing models.Category
		if err := database.DB.Where("name = ? AND id != ?", *req.Name, categoryID).First(&existing).Error; err == nil {
			BadRequest(c, "分类名已存在")
			return
		}
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&category).Updates(updates).Error; err != nil {
			InternalServerError(c, "更新失败")
			return
		}
	}

	Success(c, gin.H{
		"message": "更新成功",
	})
}

func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		NotFound(c, "分类不存在")
		return
	}

	// 检查是否有小说使用此分类
	var novelCount int64
	database.DB.Model(&models.Novel{}).Where("category_id = ?", categoryID).Count(&novelCount)
	if novelCount > 0 {
		BadRequest(c, "该分类下还有小说，无法删除")
		return
	}

	// 软删除
	if err := database.DB.Delete(&category).Error; err != nil {
		InternalServerError(c, "删除失败")
		return
	}

	Success(c, gin.H{
		"message": "删除成功",
	})
}
