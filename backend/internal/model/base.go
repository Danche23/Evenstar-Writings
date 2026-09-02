package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity 基础实体（软删除表使用：users / articles / comments）
type BaseEntity struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// BeforeCreate GORM hook - 创建前
func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate GORM hook - 更新前
func (b *BaseEntity) BeforeUpdate(tx *gorm.DB) error {
	return nil
}
