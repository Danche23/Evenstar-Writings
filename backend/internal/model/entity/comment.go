package entity

import "time"

// Comment 评论表（最多两级：parent_id 恒指一级评论；reply_to_id 仅展示用）
type Comment struct {
	BaseEntity
	ArticleID uint       `gorm:"column:article_id;not null;index" json:"article_id"`
	UserID    *uint      `gorm:"column:user_id" json:"user_id"`           // 可空，用户注销后为 NULL
	ParentID  *uint      `gorm:"column:parent_id;index" json:"parent_id"` // 一级=NULL；二级=所属一级评论 id
	ReplyToID *uint      `gorm:"column:reply_to_id" json:"reply_to_id"`   // 实际回复对象，仅展示「张三 → 李四」，无外键
	Content   string     `gorm:"column:content;size:400;not null" json:"content"`
	IsTop     int8       `gorm:"column:is_top;not null;default:0" json:"is_top"` // 1=置顶（仅一级评论可置顶）
	TopTime   *time.Time `gorm:"column:top_time" json:"top_time"`                // 置顶时间，多条置顶按它倒序
}

// TableName 表名
func (Comment) TableName() string {
	return "comments"
}
