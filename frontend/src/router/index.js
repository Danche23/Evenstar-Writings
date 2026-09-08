import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/FrontLayout.vue'),
    children: [
      { path: '', name: 'home', component: () => import('@/views/Home.vue') },
      { path: 'articles', name: 'articles', component: () => import('@/views/ArticleList.vue') },
      { path: 'articles/:id', name: 'article', component: () => import('@/views/ArticleDetail.vue') },
      { path: 'categories', name: 'categories', component: () => import('@/views/Categories.vue') },
      { path: 'tags', name: 'tags', component: () => import('@/views/Tags.vue') },
      { path: 'about', name: 'about', component: () => import('@/views/About.vue') },
      { path: 'archive', name: 'archive', component: () => import('@/views/Archive.vue') },
      { path: 'guestbook', name: 'guestbook', component: () => import('@/views/Guestbook.vue') },
      { path: 'login', name: 'login', component: () => import('@/views/Login.vue') },
      { path: 'register', name: 'register', component: () => import('@/views/Register.vue') },
      { path: 'forgot', name: 'forgot', component: () => import('@/views/ForgotPassword.vue') },
      { path: 'profile', name: 'profile', component: () => import('@/views/Profile.vue'), meta: { auth: true } }
    ]
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { admin: true },
    children: [
      { path: '', name: 'admin-home', component: () => import('@/views/admin/Dashboard.vue') },
      { path: 'articles', name: 'admin-articles', component: () => import('@/views/admin/AdminArticles.vue') },
      { path: 'articles/new', name: 'article-new', component: () => import('@/views/admin/ArticleEdit.vue') },
      { path: 'articles/:id/edit', name: 'article-edit', component: () => import('@/views/admin/ArticleEdit.vue') },
      { path: 'categories', name: 'admin-categories', component: () => import('@/views/admin/AdminCategories.vue') },
      { path: 'categories/sort', name: 'admin-category-sort', component: () => import('@/views/admin/AdminCategorySort.vue') },
      { path: 'tags', name: 'admin-tags', component: () => import('@/views/admin/AdminTags.vue') },
      { path: 'tags/sort', name: 'admin-tag-sort', component: () => import('@/views/admin/AdminTagSort.vue') },
      { path: 'comments', name: 'admin-comments', component: () => import('@/views/admin/AdminComments.vue') },
      { path: 'users', name: 'admin-users', component: () => import('@/views/admin/AdminUsers.vue') },
      { path: 'uploads', name: 'admin-uploads', component: () => import('@/views/admin/AdminUploads.vue') },
      { path: 'about', name: 'admin-about', component: () => import('@/views/admin/AdminAbout.vue') }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// 全局前置守卫：登录态恢复 / 页面访问控制
router.beforeEach(async (to) => {
  const store = useUserStore()

  // 已有 token 且未恢复过资料 → 向后端拉取最新用户信息
  if (store.token && !store.profileLoaded && !['login', 'register', 'forgot'].includes(to.name)) {
    await store.fetchProfile()
  }

  // 需登录的页面
  if (to.meta.auth && !store.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  // 后台页面：需登录 + 管理员
  if (to.matched.some((r) => r.meta.admin)) {
    if (!store.isLoggedIn) return { name: 'login', query: { redirect: to.fullPath } }
    if (!store.isAdmin) return { name: 'home' }
  }

  return true
})

export default router
