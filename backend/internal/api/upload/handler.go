package upload

import (
	"io"
	"net/http"
	"strconv"

	"github.com/Danche23/Evenstar-Writings/internal/middleware"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// UploadHandler 上传模块处理器
type UploadHandler struct {
	uploadService *service.UploadService
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// Upload 上传图片
func (h *UploadHandler) Upload(c *gin.Context) {
	userID := middleware.GetUserID(c)
	isAdmin := middleware.GetUserRole(c) == 1
	scene := c.PostForm("scene")

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "请上传文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		response.Error(c, apperrors.CodeInternalError, "打开文件失败")
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		response.Error(c, apperrors.CodeInternalError, "读取文件失败")
		return
	}

	// 用实际内容检测 MIME（不信任前端传的 Content-Type）
	mime := http.DetectContentType(data)

	id, url, err := h.uploadService.Upload(userID, isAdmin, scene, file.Filename, data, mime)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, gin.H{"id": id, "url": url})
}

// AdminList 后台上传文件列表
func (h *UploadHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	scene := c.Query("scene")

	resp, err := h.uploadService.AdminList(page, size, scene)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}

// AdminDelete 后台删除上传文件
func (h *UploadHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperrors.CodeInvalidParam, "参数错误")
		return
	}
	force := c.Query("force") == "true"

	refs, err := h.uploadService.AdminDelete(uint(id), force)
	if err != nil {
		if refs != nil {
			response.ErrorWithData(c, apperrors.CodeConflict, "文件被文章引用", gin.H{"articles": refs})
			return
		}
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
