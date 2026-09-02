package model

// User 用户表（单管理员 + 注册访客）
type User struct {
	BaseEntity
	Username     string `gorm:"column:username;uniqueIndex;size:50;not null" json:"username"`
	Password     string `gorm:"column:password;size:100;not null" json:"-"` // bcrypt，永不输出
	Nickname     string `gorm:"column:nickname;size:50" json:"nickname"`    // 可空，空则后端生成默认昵称
	Email        string `gorm:"column:email;uniqueIndex;size:100;not null" json:"email"`
	Avatar       string `gorm:"column:avatar;size:255" json:"avatar"`             // 可空，空则前端默认头像
	Role         int8   `gorm:"column:role;not null;default:2" json:"role"`       // 1=管理员 2=普通用户
	Status       int8   `gorm:"column:status;not null;default:1" json:"status"`   // 1=正常 2=禁用
	TokenVersion uint   `gorm:"column:token_version;not null;default:0" json:"-"` // 改密/重置/禁用/删除时+1
}

// TableName 表名
func (User) TableName() string {
	return "users"
}
