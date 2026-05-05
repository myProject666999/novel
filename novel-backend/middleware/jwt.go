package middleware

import (
	"novel-backend/controllers"
	"novel-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			controllers.Unauthorized(c, "未登录或登录已过期")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			controllers.Unauthorized(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_id", claims.RoleID)

		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			controllers.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		// 假设角色ID 1是管理员
		// 实际应用中应该从数据库查询角色权限
		if roleID.(uint) != 1 {
			controllers.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

func WriterAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			controllers.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		// 假设角色ID 2是作家，1是管理员
		// 管理员也可以访问作家功能
		if roleID.(uint) != 1 && roleID.(uint) != 2 {
			controllers.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}
