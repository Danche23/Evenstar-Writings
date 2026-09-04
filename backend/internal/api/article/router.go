package article

import (
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文章模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *ArticleHandler) {
	// 前台：公开
	group.GET("/articles", h.ListArticles)
	group.GET("/articles/hot", h.HotArticles)
	group.GET("/articles/:id", h.GetArticle)
	// 浏览统计：可选登录 + IP 限流（同 IP 1 分钟 30 次）
	group.POST("/articles/:id/view", middleware.OptionalAuth(), middleware.RateLimit(30, time.Minute), h.RecordView)

	// 后台：需登录 + 管理员
	admin := group.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	{
		admin.GET("/articles", h.AdminListArticles)
		admin.POST("/articles", h.AdminCreateArticle)
		admin.GET("/articles/:id", h.AdminGetArticle)
		admin.PUT("/articles/:id", h.AdminUpdateArticle)
		admin.DELETE("/articles/:id", h.AdminDeleteArticle)
	}
}
