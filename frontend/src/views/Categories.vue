<template>
  <div class="page-container page-block">
    <div class="page-title-bar">
      <h1>分类</h1>
      <div class="sub">按分类浏览文章</div>
    </div>

    <div v-if="loading" class="loading-box">加载中…</div>
    <div v-else-if="!categories.length" class="empty"><div class="big">🗂️</div>暂无分类</div>

    <div v-else class="cat-grid">
      <div v-for="c in categories" :key="c.id" class="card cat-card">
        <div class="name">
          <router-link :to="{ name: 'articles', query: { category_id: c.id } }">{{ c.name }}</router-link>
        </div>
        <div class="count">{{ c.article_count }} 篇文章</div>
        <router-link :to="{ name: 'articles', query: { category_id: c.id } }" class="enter">查看 →</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listCategories } from '@/api/category'

const loading = ref(true)
const categories = ref([])

onMounted(async () => {
  try {
    categories.value = await listCategories()
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.cat-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.cat-card { padding: 20px 22px; display: flex; flex-direction: column; gap: 6px; }
.cat-card .name { font-size: 17px; font-weight: 700; }
.cat-card .name a:hover { color: var(--primary); }
.cat-card .count { font-size: 12.5px; color: var(--text-faint); }
.cat-card .enter { margin-top: 8px; font-size: 13px; color: var(--text-muted); }
.cat-card .enter:hover { color: var(--primary); }
@media (max-width: 900px) { .cat-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 600px) { .cat-grid { grid-template-columns: 1fr; } }
</style>
