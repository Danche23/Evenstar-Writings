package api

import (
	"github.com/Danche23/Evenstar-Writings/internal/api/article"
	"github.com/Danche23/Evenstar-Writings/internal/api/auth"
	"github.com/Danche23/Evenstar-Writings/internal/api/category"
	"github.com/Danche23/Evenstar-Writings/internal/api/comment"
	"github.com/Danche23/Evenstar-Writings/internal/api/tag"
	"github.com/Danche23/Evenstar-Writings/internal/api/upload"
	"github.com/Danche23/Evenstar-Writings/internal/api/user"
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router 路由（整个 API 的总路由入口）
type Router struct {
	// ========== 在这里添加各模块 Handler 字段（依赖注入） ==========
	authHandler *auth.AuthHandler
}

// NewRouter 创建路由
func NewRouter(authHandler *auth.AuthHandler) *Router {
	return &Router{authHandler: authHandler}
}

// Setup 设置路由（总入口）
func (r *Router) Setup(engine *gin.Engine) {
	// 全局中间件
	engine.Use(middleware.Recovery())
	engine.Use(middleware.RequestLogger())
	engine.Use(middleware.CORS())

	// 健康检查
	engine.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Evenstar API is running",
		})
	})

	// API 路由组
	apiGroup := engine.Group("/api")
	{
		auth.RegisterRoutes(apiGroup, r.authHandler)
		user.RegisterRoutes(apiGroup)
		article.RegisterRoutes(apiGroup)
		category.RegisterRoutes(apiGroup)
		tag.RegisterRoutes(apiGroup)
		comment.RegisterRoutes(apiGroup)
		upload.RegisterRoutes(apiGroup)
	}
}

// Close 关闭所有路由连接
func (r *Router) Close() error {
	return nil
}
