package user

import (
	"strconv"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/middleware"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户模块处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile 获取当前登录用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	resp, err := h.userService.GetProfile(userID)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// UpdateProfile 修改个人资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	if err := h.userService.UpdateProfile(userID, req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// UpdatePassword 修改密码
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	if err := h.userService.UpdatePassword(userID, req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// AdminListUsers 后台用户列表（分页）
func (h *UserHandler) AdminListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	resp, err := h.userService.List(page, size, keyword)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// AdminDeleteUser 后台删除用户（软删）
func (h *UserHandler) AdminDeleteUser(c *gin.Context) {
	operatorID := middleware.GetUserID(c)
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	if err := h.userService.Delete(operatorID, uint(targetID)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// AdminUpdateUserStatus 后台禁用/解禁用户
func (h *UserHandler) AdminUpdateUserStatus(c *gin.Context) {
	operatorID := middleware.GetUserID(c)
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	if err := h.userService.UpdateStatus(operatorID, uint(targetID), req.Status); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
