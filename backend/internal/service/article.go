package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"

	"github.com/redis/go-redis/v9"
)

// ArticleService 文章业务逻辑
type ArticleService struct {
	articleRepo  *repository.ArticleRepository
	userRepo     *repository.UserRepository
	categoryRepo *repository.CategoryRepository
	tagRepo      *repository.TagRepository
	redis        *redis.Client
}

// NewArticleService 创建文章服务
func NewArticleService(articleRepo *repository.ArticleRepository, userRepo *repository.UserRepository, categoryRepo *repository.CategoryRepository, tagRepo *repository.TagRepository, rdb *redis.Client) *ArticleService {
	return &ArticleService{articleRepo: articleRepo, userRepo: userRepo, categoryRepo: categoryRepo, tagRepo: tagRepo, redis: rdb}
}

// ListArticles 前台文章列表（仅已发布）
func (s *ArticleService) ListArticles(page, size int, categoryID, tagID uint, keyword string) (*dto.PageData[dto.ArticleListItem], error) {
	articles, total, err := s.articleRepo.ListPublished(page, size, categoryID, tagID, keyword)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.ArticleListItem, 0, len(articles))
	for i := range articles {
		list = append(list, s.toListItem(&articles[i]))
	}
	return dto.NewPageData(list, total, page, size), nil
}

// GetArticle 前台文章详情（仅已发布）
func (s *ArticleService) GetArticle(id uint) (*dto.ArticleDetail, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.ErrResourceNotFound
	}
	if article.Status != 2 {
		return nil, apperrors.ErrResourceNotFound
	}
	detail := s.toDetail(article)
	return &detail, nil
}

// HotArticles 热门文章（Redis ZSET 取 TOP N，过滤未发布/已删除）
func (s *ArticleService) HotArticles(limit int) ([]dto.ArticleListItem, error) {
	ctx := context.Background()
	ids, err := s.redis.ZRevRange(ctx, "hot:articles", 0, int64(limit-1)).Result()
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	if len(ids) == 0 {
		return []dto.ArticleListItem{}, nil
	}

	var articleIDs []uint
	for _, id := range ids {
		n, _ := strconv.ParseUint(id, 10, 32)
		articleIDs = append(articleIDs, uint(n))
	}

	articles, err := s.articleRepo.FindByIDs(articleIDs)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	m := make(map[uint]model.Article, len(articles))
	for _, a := range articles {
		m[a.ID] = a
	}

	// 按 ZSET 分数顺序组装，过滤未发布/已删除
	list := make([]dto.ArticleListItem, 0, len(articleIDs))
	for _, id := range articleIDs {
		a, ok := m[id]
		if !ok || a.Status != 2 {
			continue
		}
		list = append(list, s.toListItem(&a))
	}
	return list, nil
}

// RecordView 记录文章浏览（防刷 + 计数 + 热门，返回最新浏览量）
func (s *ArticleService) RecordView(articleID uint, userID uint, ip, userAgent string) (uint, error) {
	ctx := context.Background()

	article, err := s.articleRepo.FindByID(articleID)
	if err != nil || article.Status != 2 {
		return 0, apperrors.ErrResourceNotFound
	}

	// 防刷标识：登录用户按 user_id，游客按 IP + md5(UA)
	var identity string
	if userID > 0 {
		identity = fmt.Sprintf("article:view:user:%d:article:%d", userID, articleID)
	} else {
		identity = fmt.Sprintf("article:view:guest:%d:%s:%s", articleID, ip, md5Hash(userAgent))
	}

	// SETNX 15 分钟窗口：同一标识 15 分钟内同一文章只计一次
	ok, err := s.redis.SetNX(ctx, identity, 1, 15*time.Minute).Result()
	if err != nil {
		return 0, apperrors.ErrInternalError
	}
	if !ok {
		return s.currentViews(ctx, articleID, article.Views), nil
	}

	// 计数：INCR 浏览量 + ZINCRBY 热门排行
	viewKey := fmt.Sprintf("article:view:%d", articleID)
	if err := s.redis.Incr(ctx, viewKey).Err(); err != nil {
		return 0, apperrors.ErrInternalError
	}
	if err := s.redis.ZIncrBy(ctx, "hot:articles", 1, strconv.FormatUint(uint64(articleID), 10)).Err(); err != nil {
		return 0, apperrors.ErrInternalError
	}

	return s.currentViews(ctx, articleID, article.Views), nil
}

// AdminListArticles 后台文章列表（含草稿）
func (s *ArticleService) AdminListArticles(page, size int, status int8, keyword string) (*dto.PageData[dto.AdminArticleListItem], error) {
	articles, total, err := s.articleRepo.List(page, size, status, keyword)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.AdminArticleListItem, 0, len(articles))
	for i := range articles {
		a := &articles[i]
		list = append(list, dto.AdminArticleListItem{
			ID:          a.ID,
			AuthorID:    a.AuthorID,
			Title:       a.Title,
			Summary:     a.Summary,
			Cover:       a.Cover,
			Status:      a.Status,
			Views:       a.Views,
			PublishedAt: a.PublishedAt,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		})
	}
	return dto.NewPageData(list, total, page, size), nil
}

// AdminGetArticle 后台文章详情（含草稿 + category_ids/tag_ids）
func (s *ArticleService) AdminGetArticle(id uint) (*dto.AdminArticleDetail, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.ErrResourceNotFound
	}
	categoryIDs, _ := s.articleRepo.GetCategoryIDs(id)
	tagIDs, _ := s.articleRepo.GetTagIDs(id)
	return &dto.AdminArticleDetail{
		ArticleDetail: s.toDetail(article),
		Status:        article.Status,
		CategoryIDs:   categoryIDs,
		TagIDs:        tagIDs,
	}, nil
}

// AdminCreateArticle 后台创建文章
func (s *ArticleService) AdminCreateArticle(authorID uint, req dto.ArticleWriteRequest) (uint, error) {
	status := req.Status
	if status == 0 {
		status = 1
	}
	var publishedAt *time.Time
	if status == 2 {
		now := time.Now()
		publishedAt = &now
	}

	article := &model.Article{
		AuthorID:    authorID,
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		Cover:       req.Cover,
		Status:      status,
		PublishedAt: publishedAt,
	}
	if err := s.articleRepo.Create(article); err != nil {
		return 0, apperrors.ErrInternalError
	}
	if err := s.articleRepo.SetCategories(article.ID, req.CategoryIDs); err != nil {
		return 0, apperrors.ErrInternalError
	}
	if err := s.articleRepo.SetTags(article.ID, req.TagIDs); err != nil {
		return 0, apperrors.ErrInternalError
	}
	return article.ID, nil
}

// AdminUpdateArticle 后台更新文章（全量 + 关联全量替换）
func (s *ArticleService) AdminUpdateArticle(id uint, req dto.ArticleWriteRequest) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return apperrors.ErrResourceNotFound
	}
	status := req.Status
	if status == 0 {
		status = article.Status
	}

	updates := map[string]interface{}{
		"title":   req.Title,
		"summary": req.Summary,
		"content": req.Content,
		"cover":   req.Cover,
		"status":  status,
	}
	// 首次变 2 时写 published_at
	if status == 2 && article.PublishedAt == nil {
		now := time.Now()
		updates["published_at"] = &now
	}
	if err := s.articleRepo.Update(id, updates); err != nil {
		return apperrors.ErrInternalError
	}
	if err := s.articleRepo.SetCategories(id, req.CategoryIDs); err != nil {
		return apperrors.ErrInternalError
	}
	if err := s.articleRepo.SetTags(id, req.TagIDs); err != nil {
		return apperrors.ErrInternalError
	}
	return nil
}

// AdminDeleteArticle 后台删除文章（软删）
func (s *ArticleService) AdminDeleteArticle(id uint) error {
	if _, err := s.articleRepo.FindByID(id); err != nil {
		return apperrors.ErrResourceNotFound
	}
	return s.articleRepo.Delete(id)
}

// toListItem 组装文章列表项（含作者/分类/标签 + Redis 增量浏览量）
func (s *ArticleService) toListItem(a *model.Article) dto.ArticleListItem {
	return dto.ArticleListItem{
		ID:          a.ID,
		AuthorID:    a.AuthorID,
		Title:       a.Title,
		Summary:     a.Summary,
		Cover:       a.Cover,
		Views:       s.currentViews(context.Background(), a.ID, a.Views),
		PublishedAt: a.PublishedAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		Author:      s.getAuthor(a.AuthorID),
		Categories:  s.getCategories(a.ID),
		Tags:        s.getTags(a.ID),
	}
}

// toDetail 组装文章详情（含 content）
func (s *ArticleService) toDetail(a *model.Article) dto.ArticleDetail {
	return dto.ArticleDetail{
		ID:          a.ID,
		AuthorID:    a.AuthorID,
		Title:       a.Title,
		Summary:     a.Summary,
		Content:     a.Content,
		Cover:       a.Cover,
		Views:       s.currentViews(context.Background(), a.ID, a.Views),
		PublishedAt: a.PublishedAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		Author:      s.getAuthor(a.AuthorID),
		Categories:  s.getCategories(a.ID),
		Tags:        s.getTags(a.ID),
	}
}

// getAuthor 查询文章作者（CommentUser）
func (s *ArticleService) getAuthor(authorID uint) dto.CommentUser {
	user, err := s.userRepo.FindByID(authorID)
	if err != nil {
		return dto.CommentUser{}
	}
	return dto.CommentUser{ID: user.ID, Nickname: user.Nickname, Avatar: user.Avatar}
}

// getCategories 查询文章分类（CategoryBrief）
func (s *ArticleService) getCategories(articleID uint) []dto.CategoryBrief {
	ids, err := s.articleRepo.GetCategoryIDs(articleID)
	if err != nil || len(ids) == 0 {
		return []dto.CategoryBrief{}
	}
	list := make([]dto.CategoryBrief, 0, len(ids))
	for _, cid := range ids {
		c, err := s.categoryRepo.FindByID(cid)
		if err == nil {
			list = append(list, dto.CategoryBrief{ID: c.ID, Name: c.Name})
		}
	}
	return list
}

// getTags 查询文章标签（TagBrief）
func (s *ArticleService) getTags(articleID uint) []dto.TagBrief {
	ids, err := s.articleRepo.GetTagIDs(articleID)
	if err != nil || len(ids) == 0 {
		return []dto.TagBrief{}
	}
	list := make([]dto.TagBrief, 0, len(ids))
	for _, tid := range ids {
		t, err := s.tagRepo.FindByID(tid)
		if err == nil {
			list = append(list, dto.TagBrief{ID: t.ID, Name: t.Name})
		}
	}
	return list
}

// currentViews 计算当前浏览量 = MySQL views + Redis 未回写增量
func (s *ArticleService) currentViews(ctx context.Context, articleID, baseViews uint) uint {
	incr, err := s.redis.Get(ctx, fmt.Sprintf("article:view:%d", articleID)).Uint64()
	if err != nil {
		return baseViews
	}
	return baseViews + uint(incr)
}

// md5Hash 计算字符串 md5
func md5Hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// SyncViews 回写 Redis 浏览量增量到 MySQL（cron 每 5 分钟调用，GETDEL 原子取数）
func (s *ArticleService) SyncViews() error {
	ctx := context.Background()
	keys, err := s.redis.Keys(ctx, "article:view:*").Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		val, err := s.redis.GetDel(ctx, key).Int64()
		if err != nil || val <= 0 {
			continue
		}
		idStr := strings.TrimPrefix(key, "article:view:")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			continue
		}
		_ = s.articleRepo.IncrViews(uint(id), val)
	}
	return nil
}
