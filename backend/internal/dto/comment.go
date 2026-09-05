package dto

import "time"

// Comment 评论（含二级回复）
type Comment struct {
	ID        uint         `json:"id"`
	ArticleID uint         `json:"article_id"`
	UserID    *uint        `json:"user_id"`
	ParentID  *uint        `json:"parent_id"`
	ReplyToID *uint        `json:"reply_to_id"`
	Content   string       `json:"content"`
	IsTop     int8         `json:"is_top"`
	TopTime   *time.Time   `json:"top_time"`
	CreatedAt time.Time    `json:"created_at"`
	User      *CommentUser `json:"user"`
	Deleted   bool         `json:"deleted"`
	Replies   []Comment    `json:"replies"`
}

// CommentWriteRequest 发表评论请求
type CommentWriteRequest struct {
	Content   string `json:"content" binding:"required,max=400"`
	ParentID  *uint  `json:"parent_id"`  // 一级评论 id（回复时恒指一级）
	ReplyToID *uint  `json:"reply_to_id"` // 实际回复对象 id，仅展示「张三 → 李四」
}

// AdminComment 后台评论列表项（扁平，不含嵌套）
type AdminComment struct {
	ID           uint       `json:"id"`
	ArticleID    uint       `json:"article_id"`
	ArticleTitle string     `json:"article_title"`
	UserID       *uint      `json:"user_id"`
	UserNickname string     `json:"user_nickname"`
	ParentID     *uint      `json:"parent_id"`
	ReplyToID    *uint      `json:"reply_to_id"`
	Content      string     `json:"content"`
	IsTop        int8       `json:"is_top"`
	TopTime      *time.Time `json:"top_time"`
	CreatedAt    time.Time  `json:"created_at"`
}
