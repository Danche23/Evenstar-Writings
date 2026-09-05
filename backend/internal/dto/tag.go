package dto

import "time"

// Tag 标签（含已发布文章数）
type Tag struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	ArticleCount int64     `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TagWriteRequest 标签创建/编辑请求
type TagWriteRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}
