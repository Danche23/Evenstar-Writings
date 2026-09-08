package service

import (
	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
)

// TagService 标签业务逻辑
type TagService struct {
	tagRepo *repository.TagRepository
}

// NewTagService 创建标签服务
func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

// List 前台标签列表（带已发布文章数）
func (s *TagService) List() ([]dto.Tag, error) {
	tags, err := s.tagRepo.ListAll()
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.Tag, 0, len(tags))
	for _, t := range tags {
		count, _ := s.tagRepo.CountPublishedArticles(t.ID)
		list = append(list, dto.Tag{
			ID:           t.ID,
			Name:         t.Name,
			SortOrder:    t.SortOrder,
			ArticleCount: count,
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
		})
	}
	return list, nil
}

// Create 新建标签（重名返回业务错误）
func (s *TagService) Create(name string) (uint, error) {
	if _, err := s.tagRepo.FindByName(name); err == nil {
		return 0, apperrors.ErrResourceAlreadyExists
	} else if !errorsIsNotFound(err) {
		return 0, apperrors.ErrInternalError
	}
	// 新标签默认排到末尾（当前数量 + 1）
	existing, _ := s.tagRepo.ListAll()
	tag := &model.Tag{Name: name, SortOrder: len(existing) + 1}
	if err := s.tagRepo.Create(tag); err != nil {
		return 0, apperrors.ErrInternalError
	}
	return tag.ID, nil
}

// Reorder 按有序 id 数组重排前台展示顺序
func (s *TagService) Reorder(ids []uint) error {
	if len(ids) == 0 {
		return apperrors.ErrInvalidParam
	}
	if err := s.tagRepo.Reorder(ids); err != nil {
		return apperrors.ErrInternalError
	}
	return nil
}

// Update 编辑标签（重名返回业务错误）
func (s *TagService) Update(id uint, name string) error {
	if _, err := s.tagRepo.FindByID(id); err != nil {
		return apperrors.ErrResourceNotFound
	}
	if existing, err := s.tagRepo.FindByName(name); err == nil && existing.ID != id {
		return apperrors.ErrResourceAlreadyExists
	}
	return s.tagRepo.Update(id, map[string]interface{}{"name": name})
}

// Delete 删除标签（硬删，中间表 CASCADE）
func (s *TagService) Delete(id uint) error {
	if _, err := s.tagRepo.FindByID(id); err != nil {
		return apperrors.ErrResourceNotFound
	}
	return s.tagRepo.Delete(id)
}
