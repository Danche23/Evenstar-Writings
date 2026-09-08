import { defineStore } from 'pinia'
import { getToken, setToken, setStoredUser, getStoredUser, clearAuth } from '@/utils/token'
import { getProfile } from '@/api/user'

// 用户全局状态：登录态 + 当前用户信息（多页面共享）
export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken(),
    user: getStoredUser(), // { id, username, nickname, email, avatar, role, status }
    profileLoaded: false
  }),
  getters: {
    isLoggedIn: (s) => !!s.token,
    isAdmin: (s) => !!s.user && s.user.role === 1,
    displayName: (s) => s.user?.nickname || s.user?.username || ''
  },
  actions: {
    // 登录/注册成功后写入登录态
    setAuth(token, user) {
      this.token = token
      this.user = user
      setToken(token)
      setStoredUser(user)
    },
    setUser(user) {
      this.user = user
      setStoredUser(user)
    },
    logout() {
      this.token = ''
      this.user = null
      this.profileLoaded = false
      clearAuth()
    },
    // 刷新页面后向后端恢复最新用户信息（token 校验 + status + token_version）
    async fetchProfile() {
      if (!this.token) return
      try {
        const data = await getProfile()
        this.user = data
        setStoredUser(data)
      } catch (e) {
        // token 失效由拦截器统一处理（清登录态 + 跳登录页）
        if (e.code === 401 || e.code === 1005 || e.code === 1006) this.logout()
      } finally {
        this.profileLoaded = true
      }
    }
  }
})
