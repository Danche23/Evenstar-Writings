package dto

import "time"

// Message 留言
type Message struct {
	ID        uint         `json:"id"`
	UserID    *uint        `json:"user_id"`
	User      *CommentUser `json:"user"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"created_at"`
}

// MessageWriteRequest 留言请求
type MessageWriteRequest struct {
	Content string `json:"content" binding:"required,max=400"`
}

// MessageListResponse 留言列表响应
type MessageListResponse struct {
	List     []Message `json:"list"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
