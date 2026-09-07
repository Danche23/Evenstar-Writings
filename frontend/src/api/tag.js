import request from './request'

export function listTags() {
  return request.get('/api/tags')
}

export function adminCreateTag(name) {
  return request.post('/api/admin/tags', { name })
}

export function adminUpdateTag(id, name) {
  return request.put(`/api/admin/tags/${id}`, { name })
}

export function adminDeleteTag(id) {
  return request.delete(`/api/admin/tags/${id}`)
}
