package message

import (
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册留言板模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *MessageHandler) {
	group.GET("/messages", h.List)
	group.POST("/messages", middleware.Auth(), h.Create)
	group.DELETE("/messages/:id", middleware.Auth(), h.Delete)
}
