package dto

import "time"

// ArticleListItem 文章列表项（前台，不含 content）
type ArticleListItem struct {
	ID          uint            `json:"id"`
	AuthorID    uint            `json:"author_id"`
	Title       string          `json:"title"`
	Summary     string          `json:"summary"`
	Cover       string          `json:"cover"`
	Views       uint            `json:"views"`
	PublishedAt *time.Time      `json:"published_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Author      CommentUser     `json:"author"`
	Categories  []CategoryBrief `json:"categories"`
	Tags        []TagBrief      `json:"tags"`
}

// ArticleDetail 文章详情（前台，含 content）
type ArticleDetail struct {
	ID          uint            `json:"id"`
	AuthorID    uint            `json:"author_id"`
	Title       string          `json:"title"`
	Summary     string          `json:"summary"`
	Content     string          `json:"content"`
	Cover       string          `json:"cover"`
	Views       uint            `json:"views"`
	PublishedAt *time.Time      `json:"published_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Author      CommentUser     `json:"author"`
	Categories  []CategoryBrief `json:"categories"`
	Tags        []TagBrief      `json:"tags"`
}

// AdminArticleListItem 后台文章列表项（含 status）
type AdminArticleListItem struct {
	ID          uint       `json:"id"`
	AuthorID    uint       `json:"author_id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Cover       string     `json:"cover"`
	Status      int8       `json:"status"`
	Views       uint       `json:"views"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AdminArticleDetail 后台文章详情（含 category_ids/tag_ids 供编辑页预选）
type AdminArticleDetail struct {
	ArticleDetail
	Status      int8   `json:"status"`
	CategoryIDs []uint `json:"category_ids"`
	TagIDs      []uint `json:"tag_ids"`
}

// ArticleWriteRequest 文章创建/编辑请求（status 1=草稿 2=发布）
type ArticleWriteRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Summary     string `json:"summary" binding:"omitempty,max=500"`
	Content     string `json:"content" binding:"required"`
	Cover       string `json:"cover" binding:"omitempty,max=255"`
	Status      int8   `json:"status" binding:"omitempty,oneof=1 2"`
	CategoryIDs []uint `json:"category_ids"`
	TagIDs      []uint `json:"tag_ids"`
}
