package article

import (
	"strconv"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/middleware"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// ArticleHandler 文章模块处理器
type ArticleHandler struct {
	articleService *service.ArticleService
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

// ListArticles 前台文章列表
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	tagID, _ := strconv.ParseUint(c.Query("tag_id"), 10, 32)
	keyword := c.Query("keyword")

	resp, err := h.articleService.ListArticles(page, size, uint(categoryID), uint(tagID), keyword)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// GetArticle 前台文章详情
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	resp, err := h.articleService.GetArticle(uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// HotArticles 热门文章
func (h *ArticleHandler) HotArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	resp, err := h.articleService.HotArticles(limit)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// RecordView 记录文章浏览（可选登录）
func (h *ArticleHandler) RecordView(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	userID := middleware.GetUserID(c)
	views, err := h.articleService.RecordView(uint(id), userID, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, gin.H{"views": views})
}

// AdminListArticles 后台文章列表
func (h *ArticleHandler) AdminListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status, _ := strconv.Atoi(c.Query("status"))
	keyword := c.Query("keyword")

	resp, err := h.articleService.AdminListArticles(page, size, int8(status), keyword)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// AdminGetArticle 后台文章详情
func (h *ArticleHandler) AdminGetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	resp, err := h.articleService.AdminGetArticle(uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// AdminCreateArticle 后台创建文章
func (h *ArticleHandler) AdminCreateArticle(c *gin.Context) {
	var req dto.ArticleWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	authorID := middleware.GetUserID(c)
	id, err := h.articleService.AdminCreateArticle(authorID, req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, gin.H{"id": id})
}

// AdminUpdateArticle 后台更新文章
func (h *ArticleHandler) AdminUpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	var req dto.ArticleWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请求参数错误")
		return
	}
	if err := h.articleService.AdminUpdateArticle(uint(id), req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// AdminDeleteArticle 后台删除文章
func (h *ArticleHandler) AdminDeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	if err := h.articleService.AdminDeleteArticle(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
