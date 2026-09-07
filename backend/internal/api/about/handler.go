package about

import (
	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// AboutHandler 前台「关于」页内容处理器
type AboutHandler struct {
	aboutService *service.AboutService
}

// NewAboutHandler 创建关于页处理器
func NewAboutHandler(aboutService *service.AboutService) *AboutHandler {
	return &AboutHandler{aboutService: aboutService}
}

// Get 读取关于页内容（公开与后台共用；未保存过返回默认文案）
func (h *AboutHandler) Get(c *gin.Context) {
	doc, err := h.aboutService.Get()
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, doc)
}

// Update 保存关于页内容（后台）
func (h *AboutHandler) Update(c *gin.Context) {
	var req dto.AboutContent
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	if err := h.aboutService.Save(&req); err != nil {
		response.BizError(c, err)
		return
	}
	// 回读已保存内容，便于前端确认
	doc, err := h.aboutService.Get()
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, doc)
}
