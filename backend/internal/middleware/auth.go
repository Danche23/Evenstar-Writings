package middleware

import (
	"strings"

	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/pkg/database"
	"github.com/Danche23/Evenstar-Writings/pkg/jwt"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	// ContextUserID 用户 ID 上下文键
	ContextUserID = "user_id"
	// ContextUsername 用户名 上下文键
	ContextUsername = "username"
	// ContextUserRole 用户角色 上下文键
	ContextUserRole = "user_role"
)

// Auth JWT 认证中间件：验签 + 查库校验用户存在 / status=1 / token_version 匹配
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请提供认证令牌")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "令牌格式错误")
			c.Abort()
			return
		}

		// 验签解析 Token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 查库校验：用户存在 + status=1 + token_version 匹配（JWT 失效机制 M1）
		var user model.User
		if err := database.GetMySQL().First(&user, claims.UserID).Error; err != nil {
			response.Unauthorized(c, "用户不存在")
			c.Abort()
			return
		}
		if user.Status != 1 {
			response.Unauthorized(c, "用户已被禁用")
			c.Abort()
			return
		}
		if user.TokenVersion != claims.TV {
			response.Unauthorized(c, "令牌已失效，请重新登录")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextUserID, claims.GetUserID())
		c.Set(ContextUsername, claims.GetUsername())
		c.Set(ContextUserRole, user.Role)

		c.Next()
	}
}

// AdminOnly 管理员权限中间件（需在 Auth 之后使用）
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserRole(c) != 1 {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get(ContextUserID); exists {
		return userID.(uint)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(ContextUsername); exists {
		return username.(string)
	}
	return ""
}

// GetUserRole 从上下文获取用户角色（1=管理员 2=普通用户）
func GetUserRole(c *gin.Context) int8 {
	if role, exists := c.Get(ContextUserRole); exists {
		return role.(int8)
	}
	return 0
}

// OptionalAuth 可选的 JWT 认证中间件（有 token 就解析，无 token 或失效则放行，不强制）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err == nil {
			c.Set(ContextUserID, claims.GetUserID())
			c.Set(ContextUsername, claims.GetUsername())
		}

		c.Next()
	}
}
