<template>
  <div class="site">
    <!-- 顶部导航：纸感毛玻璃 + 衬线 logo -->
    <header class="site-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <span class="logo-mark">暮</span>
          <span class="logo-text">暮星随笔<em>Evenstar Writings</em></span>
        </router-link>

        <nav class="main-nav">
          <router-link to="/" :class="{ active: isActive('/') }">首页</router-link>
          <router-link to="/articles" :class="{ active: isArticleRoute }">文章</router-link>
          <router-link to="/categories" :class="{ active: isActive('/categories') }">分类</router-link>
          <router-link to="/tags" :class="{ active: isActive('/tags') }">标签</router-link>
          <router-link to="/archive" :class="{ active: isActive('/archive') }">归档</router-link>
          <router-link to="/guestbook" :class="{ active: isActive('/guestbook') }">留言板</router-link>
          <router-link to="/about" :class="{ active: isActive('/about') }">关于</router-link>
        </nav>

        <div class="header-actions">
          <template v-if="!store.isLoggedIn">
            <router-link to="/login" class="btn btn-outline btn-sm">登录</router-link>
            <router-link to="/register" class="btn btn-primary btn-sm">注册</router-link>
          </template>
          <el-dropdown v-else @command="onCommand">
            <span class="user-chip">
              <span class="u-avatar-wrap">
                <img v-if="store.user?.avatar" :src="store.user.avatar" class="u-avatar" alt="" />
                <span v-else class="u-avatar u-avatar-text">{{ store.displayName.slice(0, 1) }}</span>
              </span>
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

    <!-- 底部：夜暮深蓝 -->
    <footer class="site-footer">
      <div class="footer-poem">
        <p class="fp-line">
          <span class="fp-greet">❀ {{ greeting }}</span>
          <span class="fp-quote poem">“{{ quote }}”</span>
        </p>
      </div>
      <div class="footer-inner">
        <div class="footer-brand">
          <router-link to="/" class="logo"><span class="logo-mark">暮</span>暮星随笔</router-link>
          <p>暮色四合，星子初亮 —— 记录技术、生活与此刻的想法。</p>
        </div>
        <div class="footer-links">
          <router-link to="/articles">文章</router-link>
          <router-link to="/categories">分类</router-link>
          <router-link to="/tags">标签</router-link>
          <router-link to="/archive">归档</router-link>
          <router-link to="/guestbook">留言板</router-link>
          <router-link to="/about">关于</router-link>
        </div>
        <div class="footer-copy">© {{ year }} Evenstar Writings · 暮星随笔</div>
      </div>
      <div class="footer-stats">
        <span v-if="stats">共 <b>{{ stats.articles }}</b> 篇文章 · <b>{{ stats.comments }}</b> 条评论 · <b>{{ stats.views }}</b> 次阅读</span>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDown } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { greetNow, dailyQuote } from '@/utils/quotes'
import { getPublicStats } from '@/api/stats'

const route = useRoute()
const router = useRouter()
const store = useUserStore()
const year = new Date().getFullYear()
const greeting = greetNow()
const quote = dailyQuote()
const stats = ref(null)

onMounted(async () => {
  try {
    stats.value = await getPublicStats()
  } catch (e) { /* 页脚统计失败静默 */ }
})

const isArticleRoute = computed(
  () => route.path.startsWith('/articles') && route.path !== '/articles'
)
const isActive = (p) => (p === '/' ? route.path === '/' : route.path.startsWith(p))

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
  background: rgba(246, 241, 231, 0.88);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border);
}
.header-inner {
  max-width: 1120px; margin: 0 auto; padding: 0 24px;
  height: 62px; display: flex; align-items: center; gap: 30px;
}
.logo { display: flex; align-items: center; gap: 10px; color: var(--text); }
.logo:hover { color: var(--text); }
.logo-mark {
  width: 30px; height: 30px; border-radius: 9px;
  background: linear-gradient(135deg, #2e4168, #c3a05c);
  color: #fffdf6; display: grid; place-items: center;
  font-family: var(--font-display); font-size: 16px; font-weight: 700;
  box-shadow: 0 2px 8px rgba(46, 65, 104, 0.32);
}
.logo-text { font-family: var(--font-display); font-size: 17px; font-weight: 700; letter-spacing: 0.04em; line-height: 1.2; }
.logo-text em {
  display: block; font-style: normal; font-family: var(--font);
  font-size: 9.5px; font-weight: 400; letter-spacing: 0.18em;
  color: var(--text-faint); text-transform: uppercase; margin-top: 1px;
}
.main-nav { display: flex; gap: 26px; flex: 1; }
.main-nav a {
  position: relative; color: var(--text-muted); font-size: 14.5px; padding: 4px 0;
}
.main-nav a::after {
  content: ""; position: absolute; left: 0; right: 100%; bottom: -2px; height: 2px;
  background: var(--primary); border-radius: 2px; transition: right .22s ease;
}
.main-nav a:hover { color: var(--primary); }
.main-nav a:hover::after, .main-nav a.active::after { right: 0; }
.main-nav a.active { color: var(--text); font-weight: 600; }
.header-actions { display: flex; align-items: center; gap: 12px; }
.user-chip {
  display: inline-flex; align-items: center; gap: 8px;
  cursor: pointer; font-size: 14px; color: var(--text); outline: none;
}
.u-avatar-wrap { display: inline-flex; }
.u-avatar { width: 28px; height: 28px; border-radius: 50%; object-fit: cover; }
.u-avatar-text {
  display: grid; place-items: center;
  background: var(--primary-100); color: var(--primary-700);
  font-weight: 700; font-size: 13px;
}
.site-main { flex: 1; }

/* 页脚：暖墨深绿 */
/* 页脚：夜暮深蓝（压缩高度） */
.site-footer { background: var(--footer-bg); color: var(--footer-text); margin-top: 70px; border-top: 1px solid rgba(195, 160, 92, 0.35); }
.footer-poem {
  max-width: 1120px; margin: 0 auto; padding: 12px 24px 12px;
  text-align: center; border-bottom: 1px dashed rgba(170, 179, 198, .16);
}
.fp-line { display: flex; align-items: baseline; justify-content: center; gap: 14px; flex-wrap: wrap; margin: 0; }
.fp-greet { font-size: 12.5px; color: #aab3c6; }
.fp-quote { font-size: 13px; color: #e2d3a6; margin: 0; }
.footer-inner {
  max-width: 1120px; margin: 0 auto; padding: 20px 24px 16px;
  display: flex; align-items: flex-start; justify-content: space-between; gap: 24px; flex-wrap: wrap;
}
.footer-brand .logo { color: #f2ecdd; font-family: var(--font-display); font-size: 17px; font-weight: 700; margin-bottom: 8px; }
.footer-brand .logo .logo-mark { background: linear-gradient(135deg, #46639c, #c3a05c); }
.footer-brand p { font-size: 13px; color: #8f897a; margin: 0; max-width: 260px; line-height: 1.7; }
.footer-links { display: flex; gap: 22px; padding-top: 6px; }
.footer-links a { color: var(--footer-text); font-size: 13.5px; }
.footer-links a:hover { color: #f2ecdd; }
.footer-copy { font-size: 12.5px; color: #7c7667; padding-top: 6px; }
.footer-stats {
  max-width: 1120px; margin: 0 auto; padding: 10px 24px 18px; text-align: center;
  font-size: 12.5px; color: #7f8aa2;
}
.footer-stats b { color: #d9c390; font-weight: 600; }
@media (max-width: 700px) {
  .main-nav { display: none; }
  .header-inner { gap: 12px; }
  .footer-inner { flex-direction: column; gap: 18px; }
}
</style>
