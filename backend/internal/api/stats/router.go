package stats

import (
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册统计模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *StatsHandler) {
	admin := group.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	{
		admin.GET("/stats", h.Stats)
	}
}
