package repository

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 分类数据访问
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// ListAll 全量列表（按 id 升序）
func (r *CategoryRepository) ListAll() ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByID 按 ID 查分类
func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// FindByName 按名称查分类（重名校验）
func (r *CategoryRepository) FindByName(name string) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Create 创建分类
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类名称
func (r *CategoryRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.Category{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 硬删分类（中间表 CASCADE 由数据库外键处理）
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// CountPublishedArticles 统计该分类下已发布且未删除的文章数
func (r *CategoryRepository) CountPublishedArticles(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Table("article_categories AS ac").
		Joins("JOIN articles AS a ON ac.article_id = a.id").
		Where("ac.category_id = ? AND a.status = 2 AND a.deleted_at IS NULL", categoryID).
		Count(&count).Error
	return count, err
}
