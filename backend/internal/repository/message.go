package repository

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"gorm.io/gorm"
)

// MessageRepository 留言数据访问
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建留言仓库
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// List 分页留言（时间倒序，不含软删）
func (r *MessageRepository) List(page, size int) ([]model.Message, int64, error) {
	var list []model.Message
	var total int64
	q := r.db.Model(&model.Message{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Model(&model.Message{}).Order("created_at DESC").
		Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Create 新增留言
func (r *MessageRepository) Create(m *model.Message) error {
	return r.db.Create(m).Error
}

// FindByID 查询留言（不含软删）
func (r *MessageRepository) FindByID(id uint) (*model.Message, error) {
	var m model.Message
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Delete 删除留言（软删）
func (r *MessageRepository) Delete(id uint) error {
	return r.db.Delete(&model.Message{}, id).Error
}
