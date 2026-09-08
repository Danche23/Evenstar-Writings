import request from './request'

export function listComments(articleId, params) {
  // params: { page, page_size }
  return request.get(`/api/articles/${articleId}/comments`, { params })
}

export function createComment(articleId, data) {
  // { content, parent_id?, reply_to_id? }
  return request.post(`/api/articles/${articleId}/comments`, data)
}

export function deleteComment(id) {
  return request.delete(`/api/comments/${id}`)
}

// —— 后台（AdminOnly）——
export function adminListComments(params) {
  // { page, page_size, article_id?, keyword?, user? }
  return request.get('/api/admin/comments', { params })
}

export function adminDeleteComment(id) {
  return request.delete(`/api/admin/comments/${id}`)
}

export function adminToggleCommentTop(id) {
  return request.put(`/api/admin/comments/${id}/top`)
}
