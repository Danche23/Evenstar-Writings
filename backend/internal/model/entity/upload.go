package entity

import "time"

// Upload 上传文件表（硬删除，删除时先删 OSS 再删记录）
type Upload struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"user_id"`
	Scene     string    `gorm:"column:scene;size:16;not null;index" json:"scene"` // article / avatar
	Filename  string    `gorm:"column:filename;size:255;not null" json:"filename"`
	URL       string    `gorm:"column:url;size:512;not null" json:"url"`
	Size      int64     `gorm:"column:size;not null;default:0" json:"size"` // 字节
	Mime      string    `gorm:"column:mime;size:100" json:"mime"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName 表名
func (Upload) TableName() string {
	return "uploads"
}
