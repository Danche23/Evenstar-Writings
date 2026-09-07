package model

import "time"

// SiteSetting 站点设置（通用 KV 表，如 about 页内容；硬删除）
type SiteSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"column:setting_key;uniqueIndex;size:64;not null" json:"setting_key"`
	Value     string    `gorm:"column:setting_value;type:text" json:"setting_value"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (SiteSetting) TableName() string {
	return "site_settings"
}
