package repository

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"gorm.io/gorm"
)

// AboutRepository 站点设置（KV）数据访问，目前用于 about 页内容
type AboutRepository struct {
	db *gorm.DB
}

// NewAboutRepository 创建站点设置仓库
func NewAboutRepository(db *gorm.DB) *AboutRepository {
	return &AboutRepository{db: db}
}

// Get 按 key 读取（无记录返回 gorm.ErrRecordNotFound）
func (r *AboutRepository) Get(key string) (*model.SiteSetting, error) {
	var s model.SiteSetting
	if err := r.db.Where("setting_key = ?", key).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// Save 按 key 幂等写入（存在则更新 value，不存在则创建）
func (r *AboutRepository) Save(key, value string) error {
	var s model.SiteSetting
	err := r.db.Where("setting_key = ?", key).First(&s).Error
	if err == nil {
		s.Value = value
		return r.db.Save(&s).Error
	}
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(&model.SiteSetting{Key: key, Value: value}).Error
	}
	return err
}
