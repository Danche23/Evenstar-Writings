package auth

import (
	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// 认证模块
type AuthHandler struct {
	authService *service.AuthService
}

// 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	// 1. 绑定并校验请求参数
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, resp)
}
