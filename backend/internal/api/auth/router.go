package auth

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册认证模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *AuthHandler) {
	group.POST("/auth/send-code", h.SendCode)
	group.POST("/auth/register", h.Register)
	group.POST("/auth/reset-password", h.ResetPassword)
	group.POST("/auth/login", h.Login)
}
