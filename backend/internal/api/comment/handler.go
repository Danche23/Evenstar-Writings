package comment

import (
	"strconv"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/middleware"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// CommentHandler 评论模块处理器
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// ListComments 文章评论列表（两级树）
func (h *CommentHandler) ListComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := h.commentService.ListComments(uint(articleID), page, size)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// CreateComment 发表评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	var req dto.CommentWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	userID := middleware.GetUserID(c)
	resp, err := h.commentService.CreateComment(userID, uint(articleID), req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// DeleteComment 删除评论（用户删自己的，管理员删任意）
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	userID := middleware.GetUserID(c)
	isAdmin := middleware.GetUserRole(c) == 1
	if err := h.commentService.DeleteComment(userID, uint(commentID), isAdmin); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// AdminListComments 后台评论列表
func (h *CommentHandler) AdminListComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	articleID, _ := strconv.ParseUint(c.Query("article_id"), 10, 32)

	resp, err := h.commentService.AdminListComments(page, size, uint(articleID))
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// AdminDeleteComment 后台删除评论
func (h *CommentHandler) AdminDeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	if err := h.commentService.AdminDeleteComment(uint(commentID)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// AdminToggleTop 置顶/取消置顶评论
func (h *CommentHandler) AdminToggleTop(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	if err := h.commentService.AdminToggleTop(uint(commentID)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
