package upload

import (
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册上传模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *UploadHandler) {
	// 上传：需登录（场景权限在 service 层校验）
	group.POST("/upload", middleware.Auth(), h.Upload)

	// 后台：需登录 + 管理员
	admin := group.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	{
		admin.GET("/uploads", h.AdminList)
		admin.DELETE("/uploads/:id", h.AdminDelete)
	}
}
