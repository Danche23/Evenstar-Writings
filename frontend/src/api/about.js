import request from './request'

// 前台「关于」页内容
export function getAbout() {
  return request.get('/api/about')
}

// 后台读取（未保存过返回默认文案）
export function adminGetAbout() {
  return request.get('/api/admin/about')
}

// 后台保存
export function adminSaveAbout(data) {
  return request.put('/api/admin/about', data)
}
