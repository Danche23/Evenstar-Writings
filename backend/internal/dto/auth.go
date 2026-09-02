package dto

// 认证模块 DTO
//
// 建议包含的 DTO（字段参考 docs/openapi.yaml 的 auth 相关接口）：
//   - SendCodeRequest      发送邮箱验证码（email）
//   - RegisterRequest      注册（username / email / password / code）
//   - LoginRequest         登录（email / password）
//   - LoginResponse        登录成功返回（token + user 信息）
//   - ResetPasswordRequest 找回密码（email / code / new_password）
//   - CaptchaVerifyRequest 滑块验证（如需要）
