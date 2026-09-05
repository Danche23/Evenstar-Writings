package dto

// 认证模块 DTO

// SendCodeRequest 发送验证码请求（对应文档 /api/auth/send-code）
type SendCodeRequest struct {
	Email              string `json:"email" binding:"required,email"`
	Type               string `json:"type" binding:"required,oneof=register reset"` // register=注册 reset=找回密码
	CaptchaVerifyParam string `json:"captcha_verify_param" binding:"required"`      // 阿里云验证码 2.0，每次必验
}

// LoginRequest 登录请求（对应文档 /api/auth/login）
type LoginRequest struct {
	Email              string `json:"email" binding:"required,email"`
	Password           string `json:"password" binding:"required"`
	CaptchaVerifyParam string `json:"captcha_verify_param" binding:"omitempty"` // 连续失败 3 次后必填，正常登录可省略
}

// RegisterRequest 注册请求（对应文档 /api/auth/register）
// 注意：password 的「含字母和数字」规则属业务校验，在 service 层实现，DTO 只做长度校验。
type RegisterRequest struct {
	Username           string `json:"username" binding:"required,min=1,max=50"`
	Password           string `json:"password" binding:"required,min=8,max=100"`
	Email              string `json:"email" binding:"required,email"`
	Code               string `json:"code" binding:"required"`                 // 邮箱验证码（5 分钟有效）
	Nickname           string `json:"nickname" binding:"omitempty,max=50"`     // 选填，为空时后端生成默认昵称
	CaptchaVerifyParam string `json:"captcha_verify_param" binding:"required"` // 阿里云验证码 2.0，每次必验
}

// ResetPasswordRequest 找回密码请求（对应文档 /api/auth/reset-password）
type ResetPasswordRequest struct {
	Email              string `json:"email" binding:"required,email"`
	Code               string `json:"code" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=8,max=100"`
	CaptchaVerifyParam string `json:"captcha_verify_param" binding:"required"` // 阿里云验证码 2.0，每次必验
}

// LoginResponse 登录 / 注册成功响应（对应文档 LoginData.data）
type LoginResponse struct {
	Token string       `json:"token"` // JWT，有效期 7 天
	User  UserResponse `json:"user"`
}
