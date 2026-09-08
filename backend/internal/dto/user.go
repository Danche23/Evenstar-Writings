package dto

import (
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/model"
)

// UserResponse 对外返回的用户信息
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Role      int8      `json:"role"`   // 1=管理员 2=普通用户
	Status    int8      `json:"status"` // 1=正常 2=禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse
func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UpdateProfileRequest 修改个人资料请求（nickname/avatar 可选，用指针区分未传）
type UpdateProfileRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,max=50"`
	Avatar   *string `json:"avatar" binding:"omitempty,max=255"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=100"`
}

// UpdateUserStatusRequest 后台禁用/解禁用户请求
type UpdateUserStatusRequest struct {
	Status int8 `json:"status" binding:"required,oneof=1 2"` // 1=正常 2=禁用
}

// AuthorResponse 前台展示的博主信息（公开，不含敏感字段）
type AuthorResponse struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// AdminUpdateUserRequest 后台编辑用户资料（nickname/bio 可选，用指针区分未传）
type AdminUpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,max=50"`
	Bio      *string `json:"bio" binding:"omitempty,max=500"`
}
