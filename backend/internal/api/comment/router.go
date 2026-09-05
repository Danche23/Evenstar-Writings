package comment

import (
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册评论模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *CommentHandler) {
	// 前台
	group.GET("/articles/:id/comments", h.ListComments)
	// 发表评论：需登录（限流在 service 层）
	group.POST("/articles/:id/comments", middleware.Auth(), h.CreateComment)
	// 删除评论：需登录（删自己/管理员）
	group.DELETE("/comments/:id", middleware.Auth(), h.DeleteComment)

	// 后台：需登录 + 管理员
	admin := group.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	{
		admin.GET("/comments", h.AdminListComments)
		admin.DELETE("/comments/:id", h.AdminDeleteComment)
		admin.PUT("/comments/:id/top", h.AdminToggleTop)
	}
}
