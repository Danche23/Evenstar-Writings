package tag

import (
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册标签模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *TagHandler) {
	// 前台：公开
	group.GET("/tags", h.List)

	// 后台：需登录 + 管理员
	admin := group.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	{
		admin.POST("/tags", h.AdminCreate)
		admin.PUT("/tags/:id", h.AdminUpdate)
		admin.DELETE("/tags/:id", h.AdminDelete)
	}
}
