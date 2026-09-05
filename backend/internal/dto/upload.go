package dto

import "time"

// Upload 上传文件信息
type Upload struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Scene     string    `json:"scene"`
	Filename  string    `json:"filename"`
	URL       string    `json:"url"`
	Size      int64     `json:"size"`
	Mime      string    `json:"mime"`
	CreatedAt time.Time `json:"created_at"`
}

// RefArticle 引用该文件的文章（上传删除引用检查用）
type RefArticle struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status int8   `json:"status"` // 1=草稿 2=发布
}
