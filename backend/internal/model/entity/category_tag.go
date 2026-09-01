package entity

import "time"

// Category 分类表（硬删除，不嵌 BaseEntity）
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name;uniqueIndex;size:50;not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (Category) TableName() string {
	return "categories"
}

// Tag 标签表（硬删除，不嵌 BaseEntity）
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name;uniqueIndex;size:50;not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (Tag) TableName() string {
	return "tags"
}
