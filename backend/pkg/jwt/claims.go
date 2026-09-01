package jwt

import "github.com/golang-jwt/jwt/v5"

// CustomClaims 自定义 JWT Claims
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	TV       uint   `json:"tv"` // token_version，改密/重置/禁用/删除时 +1，用于旧 token 失效
	jwt.RegisteredClaims
}

// GetUserID 获取用户 ID
func (c *CustomClaims) GetUserID() uint {
	return c.UserID
}

// GetUsername 获取用户名
func (c *CustomClaims) GetUsername() string {
	return c.Username
}

// GetTokenVersion 获取 token 版本号
func (c *CustomClaims) GetTokenVersion() uint {
	return c.TV
}
