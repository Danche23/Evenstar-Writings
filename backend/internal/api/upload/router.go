package upload

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册上传模块路由（由顶层 api/router.go 统一调用）。
// 在这里把本模块各个接口挂到 group 上，代码你自己写。
func RegisterRoutes(group *gin.RouterGroup) {
	// 例如：group.POST("/upload", ...)
}
