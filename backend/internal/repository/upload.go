package repository

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"gorm.io/gorm"
)

// UploadRepository 上传文件数据访问
type UploadRepository struct {
	db *gorm.DB
}

// NewUploadRepository 创建上传仓库
func NewUploadRepository(db *gorm.DB) *UploadRepository {
	return &UploadRepository{db: db}
}

// Create 创建上传记录
func (r *UploadRepository) Create(upload *model.Upload) error {
	return r.db.Create(upload).Error
}

// FindByID 按 ID 查上传记录
func (r *UploadRepository) FindByID(id uint) (*model.Upload, error) {
	var upload model.Upload
	if err := r.db.First(&upload, id).Error; err != nil {
		return nil, err
	}
	return &upload, nil
}

// List 分页查询上传记录，scene 可选筛选
func (r *UploadRepository) List(page, size int, scene string) ([]model.Upload, int64, error) {
	var uploads []model.Upload
	var total int64
	q := r.db.Model(&model.Upload{})
	if scene != "" {
		q = q.Where("scene = ?", scene)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&uploads).Error; err != nil {
		return nil, 0, err
	}
	return uploads, total, nil
}

// Delete 硬删上传记录
func (r *UploadRepository) Delete(id uint) error {
	return r.db.Delete(&model.Upload{}, id).Error
}
