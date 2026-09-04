package dto

import "time"

// Category 分类（含已发布文章数）
type Category struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	ArticleCount int64     `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CategoryWriteRequest 分类创建/编辑请求
type CategoryWriteRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}
