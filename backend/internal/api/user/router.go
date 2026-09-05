package user

import (
	"github.com/Danche23/Evenstar-Writings/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户模块路由（由顶层 api/router.go 统一调用）
func RegisterRoutes(group *gin.RouterGroup, h *UserHandler) {
	// 前台：需登录
	authed := group.Group("")
	authed.Use(middleware.Auth())
	{
		authed.GET("/user/profile", h.GetProfile)
		authed.PUT("/user/profile", h.UpdateProfile)
		authed.PUT("/user/password", h.UpdatePassword)
	}

	// 后台：需登录 + 管理员
	admin := group.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	{
		admin.GET("/users", h.AdminListUsers)
		admin.DELETE("/users/:id", h.AdminDeleteUser)
		admin.PUT("/users/:id/status", h.AdminUpdateUserStatus)
	}
}
