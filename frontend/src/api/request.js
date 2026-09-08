import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, clearAuth } from '@/utils/token'

// Axios 实例：所有请求统一经过这里
const service = axios.create({
  // 开发期由 vite proxy 转发 /api -> http://localhost:8080，无需跨域；
  // 生产由 Nginx 同域反代，baseURL 留空即可
  baseURL: import.meta.env.VITE_API_BASE || '',
  timeout: 15000
})

// 请求拦截器：自动携带 JWT
service.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 业务码需要静默处理（不弹默认错误提示）：登录（1007/1008 由页面引导滑块）
const SILENT_URLS = ['/api/auth/login']

function isSilent(config) {
  return config?.silent || SILENT_URLS.some((u) => config?.url?.includes(u))
}

// 响应拦截器：后端统一返回 { code, message, data }
// HTTP 2xx 且 code=0 → resolve(data)；否则 reject（携带 code/message/data）
service.interceptors.response.use(
  (resp) => {
    const body = resp.data
    if (body && typeof body.code === 'number' && body.code === 0) {
      return body.data
    }
    // 极端情况：HTTP 200 但业务失败
    const err = new Error(body?.message || '请求失败')
    err.code = body?.code
    err.data = body?.data
    err.status = resp.status
    if (!isSilent(resp.config)) ElMessage.error(err.message)
    return Promise.reject(err)
  },
  (error) => {
    const resp = error.response
    const body = resp?.data
    const code = body?.code
    const silent = isSilent(error.config)

    // token 失效（401 / 1005 / 1006）：清理本地登录态并跳登录页
    if (!silent && [401, 1005, 1006].includes(code)) {
      clearAuth()
      ElMessage.warning(body?.message || '登录已过期，请重新登录')
      setTimeout(() => {
        if (!location.pathname.startsWith('/login')) location.href = '/login'
      }, 800)
    }

    const err = new Error(body?.message || (resp ? `请求失败 (${resp.status})` : '网络异常，请稍后重试'))
    err.code = code
    err.data = body?.data
    err.status = resp?.status
    if (!silent) ElMessage.error(err.message)
    return Promise.reject(err)
  }
)

export default service
