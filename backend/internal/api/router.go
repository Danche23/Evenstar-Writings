package api

import (
	"github.com/Danche23/Evenstar-Writings/internal/api/article"
	"github.com/Danche23/Evenstar-Writings/internal/api/auth"
	"github.com/Danche23/Evenstar-Writings/internal/api/category"
	"github.com/Danche23/Evenstar-Writings/internal/api/comment"
	"github.com/Danche23/Evenstar-Writings/internal/api/stats"
	"github.com/Danche23/Evenstar-Writings/internal/api/tag"
	"github.com/Danche23/Evenstar-Writings/internal/api/upload"
	"github.com/Danche23/Evenstar-Writings/internal/api/user"
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router 路由（整个 API 的总路由入口）
type Router struct {
	authHandler     *auth.AuthHandler
	userHandler     *user.UserHandler
	articleHandler  *article.ArticleHandler
	categoryHandler *category.CategoryHandler
	tagHandler      *tag.TagHandler
	commentHandler  *comment.CommentHandler
	uploadHandler   *upload.UploadHandler
	statsHandler    *stats.StatsHandler
}

// NewRouter 创建路由
func NewRouter(authHandler *auth.AuthHandler, userHandler *user.UserHandler, articleHandler *article.ArticleHandler, categoryHandler *category.CategoryHandler, tagHandler *tag.TagHandler, commentHandler *comment.CommentHandler, uploadHandler *upload.UploadHandler, statsHandler *stats.StatsHandler) *Router {
	return &Router{authHandler: authHandler, userHandler: userHandler, articleHandler: articleHandler, categoryHandler: categoryHandler, tagHandler: tagHandler, commentHandler: commentHandler, uploadHandler: uploadHandler, statsHandler: statsHandler}
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

	// 本地 mock 存储的静态文件访问（开发期，后续换 OSS 后移除）
	engine.Static("/uploads", "./storage")

	// API 路由组
	apiGroup := engine.Group("/api")
	{
		auth.RegisterRoutes(apiGroup, r.authHandler)
		user.RegisterRoutes(apiGroup, r.userHandler)
		article.RegisterRoutes(apiGroup, r.articleHandler)
		category.RegisterRoutes(apiGroup, r.categoryHandler)
		tag.RegisterRoutes(apiGroup, r.tagHandler)
		comment.RegisterRoutes(apiGroup, r.commentHandler)
		upload.RegisterRoutes(apiGroup, r.uploadHandler)
		stats.RegisterRoutes(apiGroup, r.statsHandler)
	}
}

// Close 关闭所有路由连接
func (r *Router) Close() error {
	return nil
}
