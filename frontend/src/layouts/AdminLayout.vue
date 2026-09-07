<template>
  <div class="admin-shell">
    <!-- 侧边导航 -->
    <aside class="admin-side">
      <router-link to="/" class="side-brand">
        <span class="logo-mark">E</span>Evenstar 管理后台
      </router-link>

      <nav class="side-nav">
        <div class="nav-group">概览</div>
        <router-link to="/admin" :class="{ active: isActive('/admin') }">📊 数据统计</router-link>
        <div class="nav-group">内容</div>
        <router-link to="/admin/articles" :class="{ active: isActive('/admin/articles') || isEdit }">📄 文章管理</router-link>
        <router-link to="/admin/articles/new">✏️ 写文章</router-link>
        <router-link to="/admin/comments" :class="{ active: isActive('/admin/comments') }">💬 评论管理</router-link>
        <router-link to="/admin/categories" :class="{ active: isActive('/admin/categories') }">🗂️ 分类管理</router-link>
        <router-link to="/admin/tags" :class="{ active: isActive('/admin/tags') }">🏷️ 标签管理</router-link>
        <div class="nav-group">系统</div>
        <router-link to="/admin/users" :class="{ active: isActive('/admin/users') }">👥 用户管理</router-link>
        <router-link to="/admin/uploads" :class="{ active: isActive('/admin/uploads') }">📁 上传管理</router-link>
      </nav>

      <div class="side-foot">
        <router-link to="/" class="foot-link">↩ 返回前台</router-link>
        <div class="foot-user">
          <span>{{ store.displayName }}</span>
          <a class="foot-logout" @click="onLogout">退出</a>
        </div>
      </div>
    </aside>

    <!-- 主内容 -->
    <main class="admin-main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const store = useUserStore()

function isActive(path) {
  return route.path === path || route.path.startsWith(path + '/')
}

// 编辑页（/admin/articles/:id/edit）也高亮「文章管理」
const isEdit = computed(() => /^\/admin\/articles\/\d+\/edit$/.test(route.path))

function onLogout() {
  store.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-shell { display: flex; min-height: 100vh; }
.admin-side {
  width: 216px; background: #111827; color: #9ca3af; flex-shrink: 0;
  padding: 18px 12px; display: flex; flex-direction: column;
  position: sticky; top: 0; height: 100vh;
}
.side-brand {
  display: flex; align-items: center; gap: 9px;
  font-weight: 700; font-size: 15px; color: #fff;
  padding: 6px 10px 18px; border-bottom: 1px solid #1f2937; margin-bottom: 10px;
}
.side-brand:hover { color: #fff; }
.logo-mark {
  width: 24px; height: 24px; border-radius: 7px;
  background: linear-gradient(135deg, var(--primary), #7c3aed);
  color: #fff; display: grid; place-items: center; font-size: 13px;
}
.side-nav { flex: 1; overflow-y: auto; }
.nav-group {
  font-size: 11px; text-transform: uppercase; letter-spacing: 1px;
  color: #4b5563; padding: 14px 12px 6px;
}
.side-nav a {
  display: flex; align-items: center; gap: 10px; padding: 9px 12px;
  border-radius: 8px; font-size: 14px; color: #d1d5db; margin-bottom: 2px;
}
.side-nav a:hover { background: #1f2937; color: #fff; }
.side-nav a.active { background: var(--primary); color: #fff; }
.side-foot { border-top: 1px solid #1f2937; padding-top: 12px; }
.foot-link { display: block; padding: 8px 12px; font-size: 13px; color: #9ca3af; border-radius: 8px; }
.foot-link:hover { background: #1f2937; color: #fff; }
.foot-user {
  display: flex; align-items: center; justify-content: space-between;
  padding: 10px 12px 2px; font-size: 13px; color: #d1d5db;
}
.foot-logout { color: #f87171; cursor: pointer; }
.foot-logout:hover { text-decoration: underline; }
.admin-main { flex: 1; padding: 28px 32px; min-width: 0; background: #f5f6f8; }
@media (max-width: 900px) {
  .admin-side { width: 64px; padding: 18px 8px; }
  .side-brand span:not(.logo-mark), .nav-group, .side-nav a { font-size: 0; }
  .side-nav a { justify-content: center; padding: 10px 4px; font-size: 16px; gap: 0; }
  .foot-user, .foot-link { font-size: 0; justify-content: center; padding: 8px 4px; }
}
</style>
