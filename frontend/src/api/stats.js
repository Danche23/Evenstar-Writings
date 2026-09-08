import request from './request'

// 后台统计（AdminOnly）
export function getStats() {
  return request.get('/api/admin/stats')
}

// 公开统计（页脚展示）
export function getPublicStats() {
  return request.get('/api/stats')
}
