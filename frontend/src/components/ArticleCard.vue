<template>
  <router-link :to="`/articles/${article.id}`" class="article-card-link">
    <article class="article-card">
      <div class="thumb">
        <img v-if="article.cover" :src="article.cover" alt="" />
        <span v-else class="thumb-emoji">{{ thumbEmoji }}</span>
      </div>
      <div class="body">
        <div class="cat-row">
          <span v-for="c in article.categories" :key="c.id" class="badge badge-primary">{{ c.name }}</span>
          <span v-if="!article.categories?.length" class="badge badge-gray">未分类</span>
        </div>
        <h3>{{ article.title }}</h3>
        <p class="summary">{{ article.summary || '（无摘要）' }}</p>
        <div class="meta">
          <span>{{ authorName }}</span>
          <span>{{ dateText }}</span>
          <span>👁 {{ article.views ?? 0 }}</span>
        </div>
      </div>
    </article>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import { formatDate } from '@/utils/format'

// 文章列表卡片：展示封面/分类/标题/摘要/元信息（与原型 article-card 一致）
const props = defineProps({
  article: { type: Object, required: true }
})

const EMOJI_POOL = ['📄', '📝', '💡', '⚙️', '🔐', '🐳', '🗂️', '🧩', '🚀', '📚']

const authorName = computed(() => {
  const a = props.article.author
  return a?.nickname || '灯影'
})

const dateText = computed(() =>
  formatDate(props.article.published_at || props.article.created_at)
)

// 无封面时用一个与标题绑定的稳定 emoji 占位
const thumbEmoji = computed(() => {
  const title = props.article.title || ''
  let hash = 0
  for (let i = 0; i < title.length; i++) hash = (hash * 31 + title.charCodeAt(i)) >>> 0
  return EMOJI_POOL[hash % EMOJI_POOL.length]
})
</script>

<style scoped>
.article-card-link { display: block; color: inherit; }
.article-card {
  display: flex; gap: 18px; padding: 18px 20px; margin-bottom: 14px;
  background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius);
  transition: box-shadow 0.15s, transform 0.15s;
}
.article-card:hover { box-shadow: var(--shadow); transform: translateY(-1px); }
.thumb {
  width: 150px; height: 96px; border-radius: 8px; flex-shrink: 0; overflow: hidden;
  display: grid; place-items: center;
  background: linear-gradient(135deg, #e0e7ff, #c7d2fe);
}
.thumb img { width: 100%; height: 100%; object-fit: cover; }
.thumb-emoji { font-size: 28px; color: var(--primary-600); }
.body { flex: 1; min-width: 0; }
.cat-row { display: flex; gap: 6px; margin-bottom: 6px; flex-wrap: wrap; }
h3 { font-size: 16px; font-weight: 600; line-height: 1.4; margin-bottom: 6px; }
h3:hover { color: var(--primary); }
.summary {
  font-size: 13px; color: var(--text-muted); margin-bottom: 10px;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}
.meta { font-size: 12px; color: var(--text-faint); display: flex; gap: 14px; align-items: center; }
@media (max-width: 640px) {
  .thumb { display: none; }
}
</style>
