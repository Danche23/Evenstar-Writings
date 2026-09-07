import request from './request'

// 文章模块（前台浏览 / 后台管理）

export function listArticles(params) {
  // { page, page_size, category_id?, tag_id?, keyword? }
  return request.get('/api/articles', { params })
}

export function getArticle(id) {
  return request.get(`/api/articles/${id}`)
}

export function getHotArticles(limit = 10) {
  return request.get('/api/articles/hot', { params: { limit } })
}

export function recordView(id) {
  // 详情页打开 5 秒后调用，防刷由后端控制
  return request.post(`/api/articles/${id}/view`)
}

// —— 后台（AdminOnly）——
export function adminListArticles(params) {
  // { page, page_size, status?, keyword? }
  return request.get('/api/admin/articles', { params })
}

export function adminGetArticle(id) {
  return request.get(`/api/admin/articles/${id}`)
}

export function adminCreateArticle(data) {
  // { title, summary?, content, cover?, status, category_ids[], tag_ids[] }
  return request.post('/api/admin/articles', data)
}

export function adminUpdateArticle(id, data) {
  return request.put(`/api/admin/articles/${id}`, data)
}

export function adminDeleteArticle(id) {
  return request.delete(`/api/admin/articles/${id}`)
}
