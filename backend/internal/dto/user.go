package dto

// 用户模块 DTO
//
// 建议包含的 DTO（字段参考 docs/openapi.yaml 的 user 相关接口）：
//   - UserResponse         用户信息（脱敏后，绝不包含 password / token_version）
//   - UpdateProfileRequest 修改个人资料（nickname / email / avatar）
//   - ChangePasswordRequest 修改密码（old_password / new_password）
//   - UserListResponse     后台用户管理列表（分页 + 角色/状态）
