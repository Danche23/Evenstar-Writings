package tag

import (
	"strconv"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// TagHandler 标签模块处理器
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler 创建标签处理器
func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// List 前台标签列表
func (h *TagHandler) List(c *gin.Context) {
	list, err := h.tagService.List()
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, list)
}

// AdminCreate 后台新建标签
func (h *TagHandler) AdminCreate(c *gin.Context) {
	var req dto.TagWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	id, err := h.tagService.Create(req.Name)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, gin.H{"id": id})
}

// AdminUpdate 后台编辑标签
func (h *TagHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	var req dto.TagWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	if err := h.tagService.Update(uint(id), req.Name); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// AdminDelete 后台删除标签
func (h *TagHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	if err := h.tagService.Delete(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
