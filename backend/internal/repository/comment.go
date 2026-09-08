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

// ListTopLevel 一级评论分页（仅未删除；排序：置顶在前按 top_time 倒序，普通按时间倒序）
func (r *CommentRepository) ListTopLevel(articleID uint, page, size int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	cond := "article_id = ? AND parent_id IS NULL"
	// Count 必须显式指定 Model，否则 GORM 无法确定表名（会报 Table not set）
	if err := r.db.Model(&model.Comment{}).Where(cond, articleID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 列表单独构建链式条件，避免与 Count 复用同一 statement
	if err := r.db.Where(cond, articleID).
		Order("is_top DESC, top_time DESC, created_at DESC").
		Offset((page - 1) * size).Limit(size).Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// ListReplies 二级评论（仅未删除，按时间正序）
func (r *CommentRepository) ListReplies(parentID uint) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.Where("parent_id = ?", parentID).
		Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByID 按 ID 查评论（仅未删除）
func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// CountVisibleAll 统计该文章下未删除评论总数（一级 + 二级）
func (r *CommentRepository) CountVisibleAll(articleID uint) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Comment{}).Where("article_id = ?", articleID).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// DeleteByParent 软删某一级评论下的所有二级回复（级联）
func (r *CommentRepository) DeleteByParent(parentID uint) error {
	return r.db.Where("parent_id = ?", parentID).Delete(&model.Comment{}).Error
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

// AdminList 后台评论列表（分页，支持按文章标题模糊、按用户模糊、按 article_id 精确筛选）
func (r *CommentRepository) AdminList(page, size int, articleID uint, keyword, userKeyword string) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	// 用闭包重建查询，避免 Count 与 Find 复用同一 statement 引发 GORM 条件污染
	// 注意：Select 只能用于 Find，不能用于 Count（否则生成 COUNT(comments.*) 非法 SQL）
	build := func() *gorm.DB {
		q := r.db.Model(&model.Comment{})
		if articleID > 0 {
			q = q.Where("comments.article_id = ?", articleID)
		}
		if keyword != "" {
			// 按文章标题模糊：JOIN articles
			q = q.Joins("JOIN articles ON articles.id = comments.article_id").
				Where("articles.title LIKE ?", "%"+keyword+"%")
		}
		if userKeyword != "" {
			// 按用户模糊：昵称或邮箱命中即纳入（子查询取 user_id）
			sub := r.db.Model(&model.User{}).Select("id").
				Where("nickname LIKE ? OR email LIKE ?", "%"+userKeyword+"%", "%"+userKeyword+"%")
			q = q.Where("comments.user_id IN (?)", sub)
		}
		return q
	}
	if err := build().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// Find 时限定只取评论表字段，避免 JOIN articles 带来的列冲突/字段歧义
	if err := build().Select("comments.*").Order("comments.id DESC").
		Offset((page - 1) * size).Limit(size).Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// DeleteWithReplies 事务内删除一级评论及其全部二级回复（保证数据一致）
func (r *CommentRepository) DeleteWithReplies(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Comment{}, id).Error
	})
}
