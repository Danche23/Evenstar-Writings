package entity

import "time"

// Article 文章表
type Article struct {
	BaseEntity
	AuthorID    uint       `gorm:"column:author_id;not null;index" json:"author_id"`
	Title       string     `gorm:"column:title;size:200;not null" json:"title"`
	Summary     string     `gorm:"column:summary;size:500" json:"summary"`
	Content     string     `gorm:"column:content;type:longtext" json:"content"`                                          // Markdown 原文
	Cover       string     `gorm:"column:cover;size:255" json:"cover"`                                                   // 可空，空则前端默认封面
	Status      int8       `gorm:"column:status;not null;default:1;index:idx_status_published,priority:1" json:"status"` // 1=草稿 2=发布
	Views       uint       `gorm:"column:views;not null;default:0" json:"views"`
	PublishedAt *time.Time `gorm:"column:published_at;index:idx_status_published,priority:2" json:"published_at"` // 可空，status 首次变 2 时写入
}

// TableName 表名
func (Article) TableName() string {
	return "articles"
}
