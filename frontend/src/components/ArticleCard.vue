<template>
  <router-link :to="`/articles/${article.id}`" class="article-card-link">
    <article class="article-card">
      <div class="thumb">
        <img v-if="article.cover" :src="article.cover" :alt="article.title" />
        <div v-else class="thumb-emoji">
          <span>{{ thumbEmoji }}</span>
          <i class="mark"></i>
        </div>
      </div>
      <div class="body">
        <div class="cat-row">
          <span v-for="c in article.categories" :key="c.id" class="badge badge-primary">{{ c.name }}</span>
          <span v-if="!article.categories?.length" class="badge badge-gray">未分类</span>
        </div>
        <h3>{{ article.title }}</h3>
        <p class="summary">{{ article.summary || '（暂无摘要，点击阅读全文）' }}</p>
        <div class="meta">
          <span class="meta-item">
            <i class="dot"></i>{{ authorName }}
          </span>
          <span class="meta-item">{{ dateText }}</span>
          <span class="meta-item views">👁 {{ article.views ?? 0 }}</span>
          <span class="read-more">阅读 →</span>
        </div>
      </div>
    </article>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import { formatDate } from '@/utils/format'

// 文章列表卡片：纸面卡片 + 封面 + 衬线标题 + 悬停上浮
const props = defineProps({
  article: { type: Object, required: true }
})

const EMOJI_POOL = ['📄', '📝', '💡', '⚙️', '🔐', '🌿', '🗂️', '🧩', '🚀', '📚']

const authorName = computed(() => props.article.author?.nickname || '暮星')

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
  display: flex; gap: 16px; padding: 13px 14px;
  margin-bottom: 12px; background: var(--surface);
  border: 1px solid var(--border); border-radius: var(--radius);
  transition: box-shadow .18s ease, transform .18s ease, border-color .18s ease;
}
.article-card:hover {
  box-shadow: var(--shadow); border-color: var(--border-strong);
}
.thumb {
  width: 132px; height: 86px; border-radius: 9px; flex-shrink: 0; overflow: hidden;
  position: relative;
}
.thumb img { width: 100%; height: 100%; object-fit: cover; }
.thumb-emoji {
  width: 100%; height: 100%; display: grid; place-items: center;
  background: linear-gradient(135deg, var(--primary-50), #f3ecdd);
  position: relative; overflow: hidden;
}
.thumb-emoji span { font-size: 26px; }
.thumb-emoji .mark {
  position: absolute; inset: 0; opacity: .5;
  background-image: radial-gradient(circle at 22% 28%, rgba(195, 160, 92, .2) 0 1.5px, transparent 2.5px),
    radial-gradient(circle at 76% 70%, rgba(70, 99, 156, .16) 0 1.5px, transparent 2.5px);
  background-size: 22px 22px, 18px 18px;
}
.body { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.cat-row { display: flex; gap: 5px; margin-bottom: 4px; flex-wrap: wrap; }
h3 {
  font-size: 15.5px; font-weight: 700; line-height: 1.5; margin-bottom: 3px;
  color: #1f2a3f;
  display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden;
}
.article-card-link:hover h3 { color: var(--primary-700); }
.summary {
  font-size: 12.5px; color: var(--text-muted); margin-bottom: 6px; line-height: 1.6;
  display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden;
}
.meta {
  margin-top: auto; font-size: 12px; color: var(--text-faint);
  display: flex; gap: 12px; align-items: center;
}
.meta-item { display: inline-flex; align-items: center; gap: 4px; white-space: nowrap; }
.meta-item .dot { width: 4px; height: 4px; border-radius: 50%; background: var(--accent); }
.read-more { margin-left: auto; color: var(--primary); font-size: 12px; opacity: 0; transform: translateX(-4px); transition: opacity .18s, transform .18s; }
.article-card:hover .read-more { opacity: 1; transform: none; }
@media (max-width: 640px) {
  .thumb { display: none; }
  .read-more { opacity: 1; transform: none; }
}
</style>
