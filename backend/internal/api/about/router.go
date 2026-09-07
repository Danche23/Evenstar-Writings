package about

import (
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册「关于」页内容路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *AboutHandler) {
	// 前台：公开读取
	group.GET("/about", h.Get)

	// 后台：需登录 + 管理员
	admin := group.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	{
		admin.GET("/about", h.Get)
		admin.PUT("/about", h.Update)
	}
}
