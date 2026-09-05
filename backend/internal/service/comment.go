package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"

	"github.com/redis/go-redis/v9"
)

// CommentService 评论业务逻辑
type CommentService struct {
	commentRepo *repository.CommentRepository
	userRepo    *repository.UserRepository
	articleRepo *repository.ArticleRepository
	redis       *redis.Client
}

// NewCommentService 创建评论服务
func NewCommentService(commentRepo *repository.CommentRepository, userRepo *repository.UserRepository, articleRepo *repository.ArticleRepository, rdb *redis.Client) *CommentService {
	return &CommentService{commentRepo: commentRepo, userRepo: userRepo, articleRepo: articleRepo, redis: rdb}
}

// ListComments 两级评论树（一级分页，二级全量带出）
func (s *CommentService) ListComments(articleID uint, page, size int) (*dto.PageData[dto.Comment], error) {
	topLevel, total, err := s.commentRepo.ListTopLevel(articleID, page, size)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.Comment, 0, len(topLevel))
	for i := range topLevel {
		c := s.toComment(&topLevel[i])
		replies, _ := s.commentRepo.ListReplies(c.ID)
		for j := range replies {
			c.Replies = append(c.Replies, s.toComment(&replies[j]))
		}
		list = append(list, c)
	}
	return dto.NewPageData(list, total, page, size), nil
}

// CreateComment 发表评论（限流 + 层级归位）
func (s *CommentService) CreateComment(userID, articleID uint, req dto.CommentWriteRequest) (*dto.Comment, error) {
	// 1. 限流：同用户 1 分钟 5 条
	ctx := context.Background()
	limitKey := fmt.Sprintf("comment:limit:%d", userID)
	count, err := s.redis.Incr(ctx, limitKey).Result()
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	if count == 1 {
		_ = s.redis.Expire(ctx, limitKey, time.Minute)
	}
	if count > 5 {
		return nil, apperrors.New(apperrors.CodeTooManyRequests, "评论太频繁，请稍后再试")
	}

	// 2. 校验文章存在
	if _, err := s.articleRepo.FindByID(articleID); err != nil {
		return nil, apperrors.ErrResourceNotFound
	}

	// 3. 层级归位：回复二级评论时，parent_id 自动取一级评论 id
	parentID := req.ParentID
	if parentID != nil {
		parent, err := s.commentRepo.FindByID(*parentID)
		if err != nil {
			return nil, apperrors.ErrResourceNotFound
		}
		if parent.ParentID != nil {
			parentID = parent.ParentID
		}
	}

	// 4. 创建评论
	comment := &model.Comment{
		ArticleID: articleID,
		UserID:    &userID,
		ParentID:  parentID,
		ReplyToID: req.ReplyToID,
		Content:   req.Content,
	}
	if err := s.commentRepo.Create(comment); err != nil {
		return nil, apperrors.ErrInternalError
	}

	result := s.toComment(comment)
	return &result, nil
}

// DeleteComment 删除评论（用户删自己的，管理员删任意）
func (s *CommentService) DeleteComment(userID, commentID uint, isAdmin bool) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return apperrors.ErrResourceNotFound
	}
	if !isAdmin && (comment.UserID == nil || *comment.UserID != userID) {
		return apperrors.ErrForbidden
	}
	return s.commentRepo.Delete(commentID)
}

// AdminListComments 后台评论列表
func (s *CommentService) AdminListComments(page, size int, articleID uint) (*dto.PageData[dto.AdminComment], error) {
	comments, total, err := s.commentRepo.AdminList(page, size, articleID)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.AdminComment, 0, len(comments))
	for i := range comments {
		c := &comments[i]
		item := dto.AdminComment{
			ID:        c.ID,
			ArticleID: c.ArticleID,
			UserID:    c.UserID,
			ParentID:  c.ParentID,
			ReplyToID: c.ReplyToID,
			Content:   c.Content,
			IsTop:     c.IsTop,
			TopTime:   c.TopTime,
			CreatedAt: c.CreatedAt,
		}
		if a, err := s.articleRepo.FindByID(c.ArticleID); err == nil {
			item.ArticleTitle = a.Title
		}
		if c.UserID != nil {
			if u, err := s.userRepo.FindByID(*c.UserID); err == nil {
				item.UserNickname = u.Nickname
			}
		}
		list = append(list, item)
	}
	return dto.NewPageData(list, total, page, size), nil
}

// AdminDeleteComment 后台删除评论（软删）
func (s *CommentService) AdminDeleteComment(commentID uint) error {
	if _, err := s.commentRepo.FindByID(commentID); err != nil {
		return apperrors.ErrResourceNotFound
	}
	return s.commentRepo.Delete(commentID)
}

// AdminToggleTop 置顶/取消置顶（仅一级评论）
func (s *CommentService) AdminToggleTop(commentID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return apperrors.ErrResourceNotFound
	}
	if comment.ParentID != nil {
		return apperrors.ErrCommentTopNotAllowed
	}
	if comment.IsTop == 1 {
		return s.commentRepo.Update(commentID, map[string]interface{}{"is_top": 0, "top_time": nil})
	}
	now := time.Now()
	return s.commentRepo.Update(commentID, map[string]interface{}{"is_top": 1, "top_time": &now})
}

// toComment 组装单条评论 DTO（软删显示「已删除」，注销用户显示 null）
func (s *CommentService) toComment(c *model.Comment) dto.Comment {
	comment := dto.Comment{
		ID:        c.ID,
		ArticleID: c.ArticleID,
		UserID:    c.UserID,
		ParentID:  c.ParentID,
		ReplyToID: c.ReplyToID,
		Content:   c.Content,
		IsTop:     c.IsTop,
		TopTime:   c.TopTime,
		CreatedAt: c.CreatedAt,
		Replies:   []dto.Comment{},
	}
	if c.DeletedAt.Valid {
		comment.Deleted = true
		comment.Content = ""
		comment.User = nil
		return comment
	}
	if c.UserID != nil {
		if user, err := s.userRepo.FindByID(*c.UserID); err == nil {
			comment.User = &dto.CommentUser{ID: user.ID, Nickname: user.Nickname, Avatar: user.Avatar}
		}
	}
	return comment
}
