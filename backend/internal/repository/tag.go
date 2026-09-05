package repository

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"gorm.io/gorm"
)

// TagRepository 标签数据访问
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓库
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// ListAll 全量列表（按 id 升序）
func (r *TagRepository) ListAll() ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.Order("id ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByID 按 ID 查标签
func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindByName 按名称查标签（重名校验）
func (r *TagRepository) FindByName(name string) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// Create 创建标签
func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// Update 更新标签名称
func (r *TagRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 硬删标签（中间表 CASCADE 由数据库外键处理）
func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

// CountPublishedArticles 统计该标签下已发布且未删除的文章数
func (r *TagRepository) CountPublishedArticles(tagID uint) (int64, error) {
	var count int64
	err := r.db.Table("article_tags AS at").
		Joins("JOIN articles AS a ON at.article_id = a.id").
		Where("at.tag_id = ? AND a.status = 2 AND a.deleted_at IS NULL", tagID).
		Count(&count).Error
	return count, err
}
