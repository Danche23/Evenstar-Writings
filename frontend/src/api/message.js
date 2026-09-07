import request from './request'

// 留言板
export function listMessages(params) {
  return request.get('/api/messages', { params })
}
export function createMessage(content) {
  return request.post('/api/messages', { content })
}
export function deleteMessage(id) {
  return request.delete(`/api/messages/${id}`)
}
