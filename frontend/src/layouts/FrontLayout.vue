<template>
  <div class="site">
    <!-- 顶部导航：logo + 主导航 + 用户区 -->
    <header class="site-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <span class="logo-mark">E</span>
          <span>Evenstar Writings</span>
        </router-link>

        <nav class="main-nav">
          <router-link to="/">首页</router-link>
          <router-link to="/articles" :class="{ active: isArticleRoute }">文章</router-link>
          <router-link to="/categories">分类</router-link>
          <router-link to="/tags">标签</router-link>
        </nav>

        <div class="header-actions">
          <template v-if="!store.isLoggedIn">
            <router-link to="/login" class="btn btn-outline btn-sm">登录</router-link>
            <router-link to="/register" class="btn btn-primary btn-sm">注册</router-link>
          </template>
          <el-dropdown v-else @command="onCommand">
            <span class="user-chip">
              <img v-if="store.user?.avatar" :src="store.user.avatar" class="u-avatar" alt="" />
              <span v-else class="u-avatar u-avatar-text">{{ store.displayName.slice(0, 1) }}</span>
              {{ store.displayName }}
              <el-icon style="margin-left: 4px"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item v-if="store.isAdmin" command="admin">管理后台</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </header>

    <!-- 页面主体 -->
    <main class="site-main">
      <router-view />
    </main>

    <!-- 底部 -->
    <footer class="site-footer">
      <div class="footer-inner">
        <router-link to="/" class="logo"><span class="logo-mark">E</span>Evenstar Writings</router-link>
        <div class="footer-links">
          <router-link to="/articles">文章</router-link>
          <router-link to="/categories">分类</router-link>
          <router-link to="/tags">标签</router-link>
        </div>
        <div class="footer-copy">© {{ year }} Evenstar Writings · 记录技术与生活</div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDown } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const store = useUserStore()
const year = new Date().getFullYear()

const isArticleRoute = computed(
  () => route.path.startsWith('/articles') || route.path === '/articles'
)

function onCommand(cmd) {
  if (cmd === 'profile') router.push('/profile')
  else if (cmd === 'admin') router.push('/admin')
  else if (cmd === 'logout') {
    store.logout()
    router.push('/')
  }
}
</script>

<style scoped>
.site { display: flex; flex-direction: column; min-height: 100vh; }
.site-header {
  position: sticky; top: 0; z-index: 100;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(8px);
  border-bottom: 1px solid var(--border);
}
.header-inner {
  max-width: 1120px; margin: 0 auto; padding: 0 24px;
  height: 60px; display: flex; align-items: center; gap: 28px;
}
.logo { display: flex; align-items: center; gap: 9px; font-weight: 700; font-size: 17px; color: var(--text); }
.logo:hover { color: var(--text); }
.logo-mark {
  width: 26px; height: 26px; border-radius: 7px;
  background: linear-gradient(135deg, var(--primary), #7c3aed);
  color: #fff; display: grid; place-items: center; font-size: 14px;
}
.main-nav { display: flex; gap: 22px; flex: 1; }
.main-nav a { color: var(--text-muted); font-size: 14px; }
.main-nav a:hover { color: var(--primary); }
.main-nav a.active { color: var(--text); font-weight: 600; }
.header-actions { display: flex; align-items: center; gap: 12px; }
.user-chip {
  display: inline-flex; align-items: center; gap: 8px;
  cursor: pointer; font-size: 14px; color: var(--text);
  outline: none;
}
.u-avatar { width: 28px; height: 28px; border-radius: 50%; object-fit: cover; }
.u-avatar-text {
  display: grid; place-items: center;
  background: var(--primary-100); color: var(--primary-700);
  font-weight: 700; font-size: 13px;
}
.site-main { flex: 1; }
.site-footer { background: var(--footer-bg); color: var(--footer-text); margin-top: 60px; }
.footer-inner {
  max-width: 1120px; margin: 0 auto; padding: 34px 24px;
  display: flex; justify-content: space-between; align-items: center; gap: 20px; flex-wrap: wrap;
}
.footer-inner .logo { color: #fff; }
.footer-links { display: flex; gap: 20px; }
.footer-links a { color: var(--footer-text); font-size: 13px; }
.footer-links a:hover { color: #fff; }
.footer-copy { font-size: 13px; }
@media (max-width: 640px) {
  .main-nav { display: none; }
  .header-inner { gap: 12px; }
}
</style>
