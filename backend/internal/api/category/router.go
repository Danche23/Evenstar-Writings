package category

import (
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册分类模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *CategoryHandler) {
	// 前台：公开
	group.GET("/categories", h.List)

	// 后台：需登录 + 管理员
	admin := group.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	{
		admin.POST("/categories", h.AdminCreate)
		admin.PUT("/categories/:id", h.AdminUpdate)
		admin.DELETE("/categories/:id", h.AdminDelete)
	}
}
