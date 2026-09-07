package service

import (
	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/model"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"

	"gorm.io/gorm"
)

// StatsService 统计业务逻辑
type StatsService struct {
	db *gorm.DB
}

// NewStatsService 创建统计服务
func NewStatsService(db *gorm.DB) *StatsService {
	return &StatsService{db: db}
}

// Stats 后台统计四卡片
func (s *StatsService) Stats() (*dto.Stats, error) {
	var articleCount, commentCount, userCount, totalViews int64

	if err := s.db.Model(&model.Article{}).Count(&articleCount).Error; err != nil {
		return nil, apperrors.ErrInternalError
	}
	if err := s.db.Model(&model.Comment{}).Count(&commentCount).Error; err != nil {
		return nil, apperrors.ErrInternalError
	}
	if err := s.db.Model(&model.User{}).Count(&userCount).Error; err != nil {
		return nil, apperrors.ErrInternalError
	}
	if err := s.db.Model(&model.Article{}).Select("COALESCE(SUM(views), 0)").Scan(&totalViews).Error; err != nil {
		return nil, apperrors.ErrInternalError
	}

	return &dto.Stats{
		ArticleCount: articleCount,
		CommentCount: commentCount,
		UserCount:    userCount,
		TotalViews:   totalViews,
	}, nil
}

// Public 公开统计（页脚：已发布文章 / 有效评论 / 总浏览量）
func (s *StatsService) Public() (*dto.PublicStats, error) {
	var articles, comments, views int64
	if err := s.db.Model(&model.Article{}).Where("status = 2").Count(&articles).Error; err != nil {
		return nil, apperrors.ErrInternalError
	}
	if err := s.db.Model(&model.Comment{}).Count(&comments).Error; err != nil {
		return nil, apperrors.ErrInternalError
	}
	if err := s.db.Model(&model.Article{}).Where("status = 2").
		Select("COALESCE(SUM(views), 0)").Scan(&views).Error; err != nil {
		return nil, apperrors.ErrInternalError
	}
	return &dto.PublicStats{Articles: articles, Comments: comments, Views: views}, nil
}
