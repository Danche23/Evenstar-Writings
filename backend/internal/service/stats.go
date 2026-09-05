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
