package auth

import (
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册认证模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *AuthHandler) {
	group.POST("/auth/send-code", h.SendCode)
	group.POST("/auth/register", h.Register)
	group.POST("/auth/reset-password", h.ResetPassword)
	// 登录：同 IP 10 分钟最多 20 次（复用现有 RateLimit 中间件，IP 维度独立限制，与邮箱维度互不替代）
	group.POST("/auth/login", middleware.RateLimit(20, 10*time.Minute), h.Login)
}
