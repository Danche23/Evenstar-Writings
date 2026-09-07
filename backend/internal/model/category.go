package model

import "time"

// Category 分类表（硬删除，不嵌 BaseEntity）
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name;uniqueIndex;size:50;not null" json:"name"`
	SortOrder int       `gorm:"column:sort_order;not null;default:0" json:"sort_order"` // 前台展示排序
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (Category) TableName() string {
	return "categories"
}
