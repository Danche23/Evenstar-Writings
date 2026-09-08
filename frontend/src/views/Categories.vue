<template>
  <div class="page-container page-block">
    <header class="page-title-bar">
      <h1>分类</h1>
      <div class="sub">共 {{ categories.length }} 个分类 · 点击进入对应文章</div>
    </header>

    <div v-if="loading" class="loading-box">加载中…</div>
    <div v-else-if="!categories.length" class="empty"><div class="big">🗂️</div>暂无分类</div>

    <div v-else class="cat-grid">
      <div v-for="(c, i) in categories" :key="c.id" class="card cat-card rise" :style="{ animationDelay: (i % 8) * 50 + 'ms' }">
        <span class="idx">{{ String(i + 1).padStart(2, '0') }}</span>
        <router-link :to="{ name: 'articles', query: { category_id: c.id } }" class="name">{{ c.name }}</router-link>
        <span class="count">{{ c.article_count }} 篇文章</span>
        <router-link :to="{ name: 'articles', query: { category_id: c.id } }" class="enter">浏览该分类 →</router-link>
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
.page-title-bar {
  max-width: 1000px; margin: 10px auto 22px; padding-bottom: 14px;
  border-bottom: 1px solid var(--border);
}
.page-title-bar h1 { font-size: 30px; }
.page-title-bar .sub { font-size: 13px; color: var(--text-muted); margin-top: 4px; }
.cat-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; max-width: 1000px; margin: 0 auto; }
.cat-card {
  position: relative; padding: 22px 24px; display: flex; flex-direction: column; gap: 6px;
  overflow: hidden; transition: transform .18s ease, box-shadow .18s ease;
}
.cat-card:hover { transform: translateY(-3px); box-shadow: var(--shadow-lg); }
.cat-card .idx {
  position: absolute; top: 8px; right: 14px;
  font-family: var(--font-display); font-size: 42px; font-weight: 700;
  color: var(--primary-100); line-height: 1; pointer-events: none;
}
.cat-card .name {
  position: relative; font-family: var(--font-display); font-size: 20px; font-weight: 700;
  color: #2c261d; display: inline-block;
}
.cat-card .name::after {
  content: ""; display: block; width: 0; height: 2px; background: var(--accent);
  transition: width .22s ease; margin-top: 2px; border-radius: 2px;
}
.cat-card:hover .name::after { width: 42px; }
.cat-card .count { font-size: 12.5px; color: var(--text-faint); }
.cat-card .enter { margin-top: 8px; font-size: 13px; color: var(--text-muted); }
.cat-card .enter:hover { color: var(--primary); }
@media (max-width: 900px) { .cat-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 600px) { .cat-grid { grid-template-columns: 1fr; } }
</style>
