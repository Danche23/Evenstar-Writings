import request from './request'

// 后台统计（AdminOnly）
export function getStats() {
  return request.get('/api/admin/stats')
}
