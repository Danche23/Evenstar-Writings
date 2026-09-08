<template>
  <div class="admin-shell">
    <!-- ✅ 手机端顶栏（仅 <=900px 显示） -->
    <div class="admin-topbar">
      <button class="admin-menu-btn" aria-label="菜单" @click="sideOpen = true">☰</button>
      <span class="admin-topbar-title">管理后台</span>
    </div>

    <!-- ✅ 手机端遮罩 -->
    <div v-if="sideOpen" class="admin-overlay" @click="sideOpen = false"></div>

    <!-- 侧边导航 -->
    <aside class="admin-side" :class="{ open: sideOpen }">
      <button class="side-close" aria-label="关闭菜单" @click="sideOpen = false">✕</button>
      <router-link to="/" class="side-brand">
        <span class="logo-mark">E</span>Evenstar 管理后台
      </router-link>

      <nav class="side-nav">
        <div class="nav-group">概览</div>
        <router-link to="/admin" :class="{ active: isActive('/admin') }">📊 数据统计</router-link>

        <div class="nav-group">内容</div>
        <router-link to="/admin/articles" :class="{ active: isActive('/admin/articles') || isEdit }">📄 文章管理</router-link>
        <router-link to="/admin/comments" :class="{ active: isActive('/admin/comments') }">💬 评论管理</router-link>

        <!-- 分类管理：父级可展开，下挂 列表 / 排序 -->
        <div class="nav-parent" :class="{ open: catOpen }">
          <button class="nav-toggle" @click="toggleGroup('cat')">
            <span>🗂️ 分类管理</span>
            <span class="arrow">▾</span>
          </button>
          <div v-show="catOpen" class="nav-children">
            <router-link to="/admin/categories" :class="{ active: isCatListActive }">分类列表</router-link>
            <router-link to="/admin/categories/sort" :class="{ active: isCatSortActive }">分类排序</router-link>
          </div>
        </div>

        <!-- 标签管理：父级可展开，下挂 列表 / 排序 -->
        <div class="nav-parent" :class="{ open: tagOpen }">
          <button class="nav-toggle" @click="toggleGroup('tag')">
            <span>🏷️ 标签管理</span>
            <span class="arrow">▾</span>
          </button>
          <div v-show="tagOpen" class="nav-children">
            <router-link to="/admin/tags" :class="{ active: isTagListActive }">标签列表</router-link>
            <router-link to="/admin/tags/sort" :class="{ active: isTagSortActive }">标签排序</router-link>
          </div>
        </div>

        <div class="nav-group">系统</div>
        <router-link to="/admin/users" :class="{ active: isActive('/admin/users') }">👥 用户管理</router-link>
        <router-link to="/admin/uploads" :class="{ active: isActive('/admin/uploads') }">📁 上传管理</router-link>
        <router-link to="/admin/about" :class="{ active: isActive('/admin/about') }">🖋️ 关于设置</router-link>
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
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const store = useUserStore()

// ✅ 手机端侧边栏抽屉开关；路由切换时自动关闭
const sideOpen = ref(false)
watch(() => route.path, () => { sideOpen.value = false })

function isActive(path) {
  return route.path === path || route.path.startsWith(path + '/')
}

// 编辑页（/admin/articles/:id/edit）也高亮「文章管理」
const isEdit = computed(() => /^\/admin\/articles\/\d+\/edit$/.test(route.path))

// 二级菜单：分类 / 标签 的展开状态（进入相关页面时默认展开对应分组）
const catOpen = ref(route.path.startsWith('/admin/categories'))
const tagOpen = ref(route.path.startsWith('/admin/tags'))

const isCatListActive = computed(() => route.path === '/admin/categories')
const isCatSortActive = computed(() => route.path === '/admin/categories/sort')
const isTagListActive = computed(() => route.path === '/admin/tags')
const isTagSortActive = computed(() => route.path === '/admin/tags/sort')

// 手风琴：展开一组时收起另一组
function toggleGroup(g) {
  if (g === 'cat') {
    catOpen.value = !catOpen.value
    if (catOpen.value) tagOpen.value = false
  } else {
    tagOpen.value = !tagOpen.value
    if (tagOpen.value) catOpen.value = false
  }
}

function onLogout() {
  store.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-shell { display: flex; min-height: 100vh; }
.admin-side {
  width: 216px; background: #121a2b; color: #98a3bd; flex-shrink: 0;
  padding: 18px 12px; display: flex; flex-direction: column;
  position: sticky; top: 0; height: 100vh;
}
.side-brand {
  display: flex; align-items: center; gap: 9px;
  font-weight: 700; font-size: 15px; color: #fff;
  padding: 6px 10px 18px; border-bottom: 1px solid #1f2940; margin-bottom: 10px;
}
.side-brand:hover { color: #fff; }
.logo-mark {
  width: 24px; height: 24px; border-radius: 7px;
  background: linear-gradient(135deg, var(--primary), #c3a05c);
  color: #fff; display: grid; place-items: center; font-size: 13px;
}
.side-nav { flex: 1; overflow-y: auto; }
.nav-group {
  font-size: 11px; text-transform: uppercase; letter-spacing: 1px;
  color: #7c8bb0; padding: 14px 12px 6px;
}
.side-nav > a {
  display: flex; align-items: center; gap: 10px; padding: 9px 12px;
  border-radius: 8px; font-size: 14px; color: #cfd8ea; margin-bottom: 2px;
}
.side-nav > a:hover { background: #1f2940; color: #fff; }
.side-nav > a.active { background: var(--primary); color: #fff; }

/* —— 可展开父级 —— */
.nav-parent { margin-bottom: 2px; }
.nav-toggle {
  width: 100%; display: flex; align-items: center; justify-content: space-between; gap: 8px;
  padding: 9px 12px; border: none; border-radius: 8px; background: transparent;
  font-size: 14px; color: #cfd8ea; cursor: pointer; text-align: left;
}
.nav-toggle:hover { background: #1f2940; color: #fff; }
.nav-parent.open .nav-toggle { color: #fff; }
.nav-toggle .arrow {
  color: #6f7d9e; font-size: 11px; transition: transform 0.18s ease; flex-shrink: 0;
}
.nav-parent.open .arrow { transform: rotate(180deg); color: #98a3bd; }

.nav-children { display: flex; flex-direction: column; gap: 2px; margin: 2px 0 4px; }
.nav-children a {
  display: flex; align-items: center; padding: 7px 12px 7px 36px;
  border-radius: 8px; font-size: 13px; color: #98a3bd;
}
.nav-children a:hover { background: #1f2940; color: #fff; }
.nav-children a.active { background: var(--primary); color: #fff; }

.side-foot { border-top: 1px solid #1f2940; padding-top: 12px; }
.foot-link { display: block; padding: 8px 12px; font-size: 13px; color: #98a3bd; border-radius: 8px; }
.foot-link:hover { background: #1f2940; color: #fff; }
.foot-user {
  display: flex; align-items: center; justify-content: space-between;
  padding: 10px 12px 2px; font-size: 13px; color: #cfd8ea;
}
.foot-logout { color: #e07a60; cursor: pointer; }
.foot-logout:hover { text-decoration: underline; }
.admin-main { flex: 1; padding: 28px 32px; min-width: 0; background: var(--bg); }

/* ✅ 手机端：侧边栏改为 off-canvas 抽屉（<=900px 生效） */
.admin-topbar { display: none; }
.admin-overlay { display: none; }
.side-close { display: none; }
@media (max-width: 900px) {
  .admin-shell { flex-direction: column; }
  .admin-topbar {
    display: flex; align-items: center; gap: 12px;
    height: 52px; padding: 0 14px; background: #121a2b; color: #fff;
    position: sticky; top: 0; z-index: 200;
  }
  .admin-menu-btn {
    display: inline-flex; align-items: center; justify-content: center;
    width: 36px; height: 36px; border-radius: 8px; font-size: 20px; line-height: 1;
    background: transparent; border: 1px solid #2a3550; color: #fff; cursor: pointer;
  }
  .admin-topbar-title { font-weight: 700; font-size: 15px; }
  .admin-overlay {
    display: block; position: fixed; inset: 0; z-index: 150;
    background: rgba(8, 12, 22, 0.5);
  }
  .admin-side {
    position: fixed; top: 0; left: 0; bottom: 0; z-index: 160;
    width: 240px; max-width: 80vw; transform: translateX(-100%);
    transition: transform 0.25s ease; box-shadow: 2px 0 16px rgba(0,0,0,.3);
  }
  .admin-side.open { transform: translateX(0); }
  .side-close {
    display: block; position: absolute; top: 14px; right: 10px;
    width: 30px; height: 30px; border-radius: 8px; border: none; cursor: pointer;
    background: #1f2940; color: #cfd8ea; font-size: 15px; line-height: 1;
  }
  .side-close:hover { background: #2a3550; color: #fff; }
  .admin-main { padding: 16px 14px; }
}
</style>
