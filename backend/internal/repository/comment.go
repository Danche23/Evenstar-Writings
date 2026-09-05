package repository

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"gorm.io/gorm"
)

// CommentRepository 评论数据访问
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// ListTopLevel 一级评论分页（含软删，排序：置顶在前按 top_time 倒序，普通按时间正序）
func (r *CommentRepository) ListTopLevel(articleID uint, page, size int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	q := r.db.Unscoped().Where("article_id = ? AND parent_id IS NULL", articleID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("is_top DESC, top_time DESC, created_at ASC").
		Offset((page - 1) * size).Limit(size).Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// ListReplies 二级评论（含软删，按时间正序）
func (r *CommentRepository) ListReplies(parentID uint) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.Unscoped().Where("parent_id = ?", parentID).
		Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByID 按 ID 查评论（含软删，用于层级归位/置顶/删除判断）
func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.Unscoped().First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// Delete 软删评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

// Update 更新评论指定字段（置顶切换等）
func (r *CommentRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Updates(updates).Error
}

// AdminList 后台评论列表（分页，可按 article_id 筛选）
func (r *CommentRepository) AdminList(page, size int, articleID uint) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	q := r.db.Model(&model.Comment{})
	if articleID > 0 {
		q = q.Where("article_id = ?", articleID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}
