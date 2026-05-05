package routes

import (
	"novel-backend/controllers"
	"novel-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 基础路由
	r.GET("/health", controllers.HealthCheck)

	// API路由组
	api := r.Group("/api")
	{
		// 公共路由（无需认证）
		public := api.Group("/public")
		{
			// 用户认证
			public.POST("/login", controllers.Login)
			public.POST("/register", controllers.Register)

			// 小说相关（公开访问）
			public.GET("/novels", controllers.GetNovelList)
			public.GET("/novels/:id", controllers.GetNovelDetail)
			public.GET("/novels/rank", controllers.GetNovelRank)
			public.GET("/novels/recommend", controllers.GetRecommendNovels)
			public.GET("/novels/search", controllers.SearchNovels)

			// 章节相关
			public.GET("/chapters/:id", controllers.GetChapterDetail)

			// 分类相关
			public.GET("/categories", controllers.GetCategories)
		}

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			// 用户相关
			user := auth.Group("/user")
			{
				user.GET("/profile", controllers.GetUserProfile)
				user.PUT("/profile", controllers.UpdateUserProfile)
				user.PUT("/password", controllers.ChangePassword)
			}

			// 读者功能
			reader := auth.Group("/reader")
			{
				// 书架
				reader.GET("/bookshelf", controllers.GetBookshelf)
				reader.POST("/bookshelf", controllers.AddToBookshelf)
				reader.DELETE("/bookshelf/:id", controllers.RemoveFromBookshelf)

				// 评论
				reader.POST("/comments", controllers.CreateComment)
				reader.GET("/comments/:novel_id", controllers.GetNovelComments)

				// 充值订阅
				reader.POST("/recharge", controllers.Recharge)
				reader.POST("/subscribe", controllers.Subscribe)
				reader.GET("/orders", controllers.GetUserOrders)
			}

			// 作家功能
			writer := auth.Group("/writer")
			writer.Use(middleware.WriterAuth())
			{
				// 小说管理
				writer.GET("/novels", controllers.GetWriterNovels)
				writer.POST("/novels", controllers.CreateNovel)
				writer.PUT("/novels/:id", controllers.UpdateNovel)
				writer.DELETE("/novels/:id", controllers.DeleteNovel)

				// 章节管理
				writer.GET("/novels/:novel_id/chapters", controllers.GetNovelChapters)
				writer.POST("/novels/:novel_id/chapters", controllers.CreateChapter)
				writer.PUT("/chapters/:id", controllers.UpdateChapter)
				writer.DELETE("/chapters/:id", controllers.DeleteChapter)
			}

			// 管理员功能
			admin := auth.Group("/admin")
			admin.Use(middleware.AdminAuth())
			{
				// 用户管理
				admin.GET("/users", controllers.GetUserList)
				admin.PUT("/users/:id/status", controllers.UpdateUserStatus)
				admin.DELETE("/users/:id", controllers.DeleteUser)

				// 角色管理
				admin.GET("/roles", controllers.GetRoleList)
				admin.POST("/roles", controllers.CreateRole)
				admin.PUT("/roles/:id", controllers.UpdateRole)
				admin.DELETE("/roles/:id", controllers.DeleteRole)

				// 菜单管理
				admin.GET("/menus", controllers.GetMenuList)
				admin.POST("/menus", controllers.CreateMenu)
				admin.PUT("/menus/:id", controllers.UpdateMenu)
				admin.DELETE("/menus/:id", controllers.DeleteMenu)

				// 小说管理
				admin.GET("/novels", controllers.GetAdminNovelList)
				admin.PUT("/novels/:id/status", controllers.UpdateNovelStatus)
				admin.DELETE("/novels/:id", controllers.DeleteNovelAdmin)

				// 评论管理
				admin.GET("/comments", controllers.GetCommentList)
				admin.DELETE("/comments/:id", controllers.DeleteCommentAdmin)

				// 作家管理
				admin.GET("/writers", controllers.GetWriterList)
				admin.PUT("/writers/:id/status", controllers.UpdateWriterStatus)

				// 邀请码管理
				admin.GET("/invite-codes", controllers.GetInviteCodeList)
				admin.POST("/invite-codes", controllers.CreateInviteCode)
				admin.DELETE("/invite-codes/:id", controllers.DeleteInviteCode)

				// 会员管理
				admin.GET("/members", controllers.GetMemberList)
				admin.PUT("/members/:id", controllers.UpdateMember)

				// 会员反馈管理
				admin.GET("/feedbacks", controllers.GetFeedbackList)
				admin.PUT("/feedbacks/:id", controllers.HandleFeedback)

				// 小说推荐管理
				admin.GET("/recommends", controllers.GetRecommendList)
				admin.POST("/recommends", controllers.CreateRecommend)
				admin.PUT("/recommends/:id", controllers.UpdateRecommend)
				admin.DELETE("/recommends/:id", controllers.DeleteRecommend)

				// 网站信息管理
				admin.GET("/site-info", controllers.GetSiteInfo)
				admin.PUT("/site-info", controllers.UpdateSiteInfo)

				// 友情链接管理
				admin.GET("/links", controllers.GetLinkList)
				admin.POST("/links", controllers.CreateLink)
				admin.PUT("/links/:id", controllers.UpdateLink)
				admin.DELETE("/links/:id", controllers.DeleteLink)

				// 新闻管理
				admin.GET("/news", controllers.GetNewsList)
				admin.POST("/news", controllers.CreateNews)
				admin.PUT("/news/:id", controllers.UpdateNews)
				admin.DELETE("/news/:id", controllers.DeleteNews)

				// 类别管理
				admin.GET("/categories", controllers.GetCategoryList)
				admin.POST("/categories", controllers.CreateCategory)
				admin.PUT("/categories/:id", controllers.UpdateCategory)
				admin.DELETE("/categories/:id", controllers.DeleteCategory)

				// 订单管理
				admin.GET("/orders", controllers.GetOrderList)
				admin.PUT("/orders/:id/status", controllers.UpdateOrderStatus)

				// 系统日志管理
				admin.GET("/logs", controllers.GetLogList)
			}
		}
	}
}
