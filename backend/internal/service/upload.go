package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/storage"

	"github.com/redis/go-redis/v9"
)

// UploadService 上传业务逻辑
type UploadService struct {
	uploadRepo  *repository.UploadRepository
	articleRepo *repository.ArticleRepository
	storage     storage.Storage
	redis       *redis.Client
}

// NewUploadService 创建上传服务
func NewUploadService(uploadRepo *repository.UploadRepository, articleRepo *repository.ArticleRepository, st storage.Storage, rdb *redis.Client) *UploadService {
	return &UploadService{uploadRepo: uploadRepo, articleRepo: articleRepo, storage: st, redis: rdb}
}

// Upload 上传图片（scene 权限 + 大小/类型校验 + avatar 月次数限制 + 存储 + 记表）
func (s *UploadService) Upload(userID uint, isAdmin bool, scene, filename string, data []byte, mime string) (uint, string, error) {
	// 1. scene 校验
	if scene != "article" && scene != "avatar" {
		return 0, "", apperrors.New(apperrors.CodeInvalidParam, "scene 仅支持 article / avatar")
	}
	// 2. 场景权限：article 仅管理员
	if scene == "article" && !isAdmin {
		return 0, "", apperrors.ErrForbidden
	}
	// 3. 大小校验（≤ 5MB）
	if len(data) > 5*1024*1024 {
		return 0, "", apperrors.New(apperrors.CodeInvalidParam, "文件大小不能超过 5MB")
	}
	// 4. 类型校验
	ext, ok := extFromMime(mime)
	if !ok {
		return 0, "", apperrors.New(apperrors.CodeInvalidParam, "仅支持 jpg/png/gif/webp 图片")
	}
	// 5. avatar 月次数限制（每用户每月 5 次）
	if scene == "avatar" {
		if err := s.checkAvatarLimit(context.Background(), userID); err != nil {
			return 0, "", err
		}
	}
	// 6. 生成 objectKey
	objectKey := s.buildObjectKey(scene, userID, ext)
	// 7. 存储
	url, err := s.storage.Save(objectKey, data)
	if err != nil {
		return 0, "", apperrors.ErrInternalError
	}
	// 8. 记 uploads 表
	upload := &model.Upload{
		UserID:   userID,
		Scene:    scene,
		Filename: filename,
		URL:      url,
		Size:     int64(len(data)),
		Mime:     mime,
	}
	if err := s.uploadRepo.Create(upload); err != nil {
		return 0, "", apperrors.ErrInternalError
	}
	return upload.ID, url, nil
}

// AdminList 后台上传文件列表
func (s *UploadService) AdminList(page, size int, scene string) (*dto.PageData[dto.Upload], error) {
	uploads, total, err := s.uploadRepo.List(page, size, scene)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.Upload, 0, len(uploads))
	for i := range uploads {
		u := &uploads[i]
		list = append(list, dto.Upload{
			ID:        u.ID,
			UserID:    u.UserID,
			Scene:     u.Scene,
			Filename:  u.Filename,
			URL:       u.URL,
			Size:      u.Size,
			Mime:      u.Mime,
			CreatedAt: u.CreatedAt,
		})
	}
	return dto.NewPageData(list, total, page, size), nil
}

// AdminDelete 删除上传文件（引用检查 → 先删存储 → 再删 DB）
// 返回引用文章列表（非空表示有引用且未强制删除，需前端确认）
func (s *UploadService) AdminDelete(id uint, force bool) ([]dto.RefArticle, error) {
	upload, err := s.uploadRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.ErrResourceNotFound
	}
	// 引用检查
	refs := s.findReferences(upload.URL)
	if len(refs) > 0 && !force {
		return refs, apperrors.New(apperrors.CodeConflict, "文件被文章引用，需管理员确认")
	}
	// 先删存储，再删 DB
	if err := s.storage.Delete(upload.URL); err != nil {
		return nil, apperrors.ErrInternalError
	}
	if err := s.uploadRepo.Delete(id); err != nil {
		return nil, apperrors.ErrInternalError
	}
	return nil, nil
}

// checkAvatarLimit avatar 每用户每月 5 次限制
func (s *UploadService) checkAvatarLimit(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("upload:avatar:%d:%s", userID, time.Now().Format("2006-01"))
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return apperrors.ErrInternalError
	}
	if count == 1 {
		_ = s.redis.Expire(ctx, key, secondsUntilMonthEnd(time.Now()))
	}
	if count > 5 {
		return apperrors.New(apperrors.CodeTooManyRequests, "头像每月最多上传 5 次")
	}
	return nil
}

// findReferences 查引用该 url 的文章
func (s *UploadService) findReferences(url string) []dto.RefArticle {
	articles, err := s.articleRepo.FindByContentURL(url)
	if err != nil {
		return nil
	}
	refs := make([]dto.RefArticle, 0, len(articles))
	for _, a := range articles {
		refs = append(refs, dto.RefArticle{ID: a.ID, Title: a.Title, Status: a.Status})
	}
	return refs
}

// buildObjectKey 生成存储路径：article/{y}/{m}/{random}.{ext} 或 avatar/{uid}/{random}.{ext}
func (s *UploadService) buildObjectKey(scene string, userID uint, ext string) string {
	now := time.Now()
	if scene == "article" {
		return fmt.Sprintf("article/%d/%d/%s%s", now.Year(), int(now.Month()), randomHex(16), ext)
	}
	return fmt.Sprintf("avatar/%d/%s%s", userID, randomHex(16), ext)
}

// extFromMime 从 MIME 推导扩展名
func extFromMime(mime string) (string, bool) {
	switch mime {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "image/gif":
		return ".gif", true
	case "image/webp":
		return ".webp", true
	}
	return "", false
}

// randomHex 生成 n 字节随机 hex 字符串
func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// secondsUntilMonthEnd 到本月末的剩余时间
func secondsUntilMonthEnd(now time.Time) time.Duration {
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	return nextMonth.Sub(now)
}
