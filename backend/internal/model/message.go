package model

// Message 留言板（一张无主题的便签墙，需登录留言；软删除）
type Message struct {
	BaseEntity
	UserID  *uint  `gorm:"column:user_id" json:"user_id"` // 可空，用户注销后为 NULL
	Content string `gorm:"column:content;size:400;not null" json:"content"`
}

// TableName 表名
func (Message) TableName() string {
	return "messages"
}
