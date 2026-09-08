import request from './request'

// 上传：需登录；scene=article(仅管理员) | avatar(所有登录用户)
export function uploadImage(file, scene) {
  const form = new FormData()
  form.append('file', file)
  form.append('scene', scene)
  return request.post('/api/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// —— 后台（AdminOnly）——
export function adminListUploads(params) {
  // { page, page_size, scene? }
  return request.get('/api/admin/uploads', { params })
}

export function adminDeleteUpload(id, force = false) {
  // silent: 409 引用冲突需页面弹确认框，由页面自行提示
  return request.delete(`/api/admin/uploads/${id}`, { params: { force }, silent: true })
}
