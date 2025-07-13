package middleware

import (
	"strings"

	"blog/internal/service"
	"blog/internal/utils"
	"blog/model/entity"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Auth(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "缺少认证token")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "认证token格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 验证token
		user, err := userService.ValidateToken(token)
		if err != nil {
			utils.Unauthorized(c, "无效的认证token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user", user)
		c.Set("user_id", user.ID)

		c.Next()
	}
}

// GetCurrentUser 从上下文获取当前用户
func GetCurrentUser(c *gin.Context) (*entity.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*entity.User), true
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}
