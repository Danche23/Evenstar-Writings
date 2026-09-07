package service

import (
	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"

	"gorm.io/gorm"
)

// CategoryService 分类业务逻辑
type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// List 前台分类列表（带已发布文章数）
func (s *CategoryService) List() ([]dto.Category, error) {
	categories, err := s.categoryRepo.ListAll()
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.Category, 0, len(categories))
	for _, c := range categories {
		count, _ := s.categoryRepo.CountPublishedArticles(c.ID)
		list = append(list, dto.Category{
			ID:           c.ID,
			Name:         c.Name,
			SortOrder:    c.SortOrder,
			ArticleCount: count,
			CreatedAt:    c.CreatedAt,
			UpdatedAt:    c.UpdatedAt,
		})
	}
	return list, nil
}

// Create 新建分类（重名返回业务错误）
func (s *CategoryService) Create(name string) (uint, error) {
	if _, err := s.categoryRepo.FindByName(name); err == nil {
		return 0, apperrors.ErrResourceAlreadyExists
	} else if !errorsIsNotFound(err) {
		return 0, apperrors.ErrInternalError
	}
	// 新分类默认排到末尾（当前数量 + 1）
	existing, _ := s.categoryRepo.ListAll()
	category := &model.Category{Name: name, SortOrder: len(existing) + 1}
	if err := s.categoryRepo.Create(category); err != nil {
		return 0, apperrors.ErrInternalError
	}
	return category.ID, nil
}

// Reorder 按有序 id 数组重排前台展示顺序
func (s *CategoryService) Reorder(ids []uint) error {
	if len(ids) == 0 {
		return apperrors.ErrInvalidParam
	}
	if err := s.categoryRepo.Reorder(ids); err != nil {
		return apperrors.ErrInternalError
	}
	return nil
}

// Update 编辑分类（重名返回业务错误）
func (s *CategoryService) Update(id uint, name string) error {
	if _, err := s.categoryRepo.FindByID(id); err != nil {
		return apperrors.ErrResourceNotFound
	}
	if existing, err := s.categoryRepo.FindByName(name); err == nil && existing.ID != id {
		return apperrors.ErrResourceAlreadyExists
	}
	return s.categoryRepo.Update(id, map[string]interface{}{"name": name})
}

// Delete 删除分类（硬删，中间表 CASCADE）
func (s *CategoryService) Delete(id uint) error {
	if _, err := s.categoryRepo.FindByID(id); err != nil {
		return apperrors.ErrResourceNotFound
	}
	return s.categoryRepo.Delete(id)
}

// errorsIsNotFound 判断是否是「记录不存在」
func errorsIsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
