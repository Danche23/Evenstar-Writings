import request from './request'

export function listCategories() {
  return request.get('/api/categories')
}

export function adminCreateCategory(name) {
  return request.post('/api/admin/categories', { name })
}

export function adminUpdateCategory(id, name) {
  return request.put(`/api/admin/categories/${id}`, { name })
}

export function adminDeleteCategory(id) {
  return request.delete(`/api/admin/categories/${id}`)
}

// 后台拖拽排序：传入有序 id 数组
export function adminReorderCategories(ids) {
  return request.put('/api/admin/categories/reorder', { ids })
}
