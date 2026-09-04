package auth

import "github.com/gin-gonic/gin"

// 注册认证模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *AuthHandler) {
	group.POST("/auth/login", h.Login)
}
