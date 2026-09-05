package category

import (
	"strconv"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类模块处理器
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// List 前台分类列表
func (h *CategoryHandler) List(c *gin.Context) {
	list, err := h.categoryService.List()
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, list)
}

// AdminCreate 后台新建分类
func (h *CategoryHandler) AdminCreate(c *gin.Context) {
	var req dto.CategoryWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	id, err := h.categoryService.Create(req.Name)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, gin.H{"id": id})
}

// AdminUpdate 后台编辑分类
func (h *CategoryHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	var req dto.CategoryWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	if err := h.categoryService.Update(uint(id), req.Name); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// AdminDelete 后台删除分类
func (h *CategoryHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	if err := h.categoryService.Delete(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
