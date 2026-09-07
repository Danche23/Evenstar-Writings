package message

import (
	"strconv"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/middleware"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// MessageHandler 留言板处理器
type MessageHandler struct {
	messageService *service.MessageService
}

// NewMessageHandler 创建留言处理器
func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// List 留言列表（公开）
func (h *MessageHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	resp, err := h.messageService.List(page, size)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// Create 写留言（需登录）
func (h *MessageHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.MessageWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	resp, err := h.messageService.Create(userID, req.Content)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// Delete 删除留言（自己或管理员）
func (h *MessageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	userID := middleware.GetUserID(c)
	isAdmin := middleware.GetUserRole(c) == 1
	if err := h.messageService.Delete(userID, uint(id), isAdmin); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
