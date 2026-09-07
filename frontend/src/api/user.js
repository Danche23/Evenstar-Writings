import request from './request'

// 用户模块（前台个人 / 后台用户管理）

export function getProfile() {
  return request.get('/api/user/profile')
}

export function updateProfile(data) {
  // { nickname?, avatar? }
  return request.put('/api/user/profile', data)
}

export function updatePassword(data) {
  // { old_password, new_password }
  return request.put('/api/user/password', data)
}

// —— 后台（AdminOnly）——
export function adminListUsers(params) {
  // { page, page_size, keyword? }
  return request.get('/api/admin/users', { params })
}

export function adminDeleteUser(id) {
  return request.delete(`/api/admin/users/${id}`)
}

export function adminUpdateUserStatus(id, status) {
  // status: 1=正常 2=禁用
  return request.put(`/api/admin/users/${id}/status`, { status })
}
