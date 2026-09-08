<template>
  <div class="page-container page-block">
    <header class="page-title-bar">
      <h1>标签</h1>
      <div class="sub">共 {{ tags.length }} 个标签 · 热度按文章数排列</div>
    </header>

    <div v-if="loading" class="loading-box">加载中…</div>
    <div v-else-if="!tags.length" class="empty"><div class="big">🏷️</div>暂无标签</div>

    <div v-else class="tag-page-body">
      <div class="tag-head-strip">
        <router-link
          v-for="t in tags.slice(0, 8)" :key="t.id"
          :to="{ name: 'articles', query: { tag_id: t.id } }"
          class="top-chip"
        >{{ t.name }} <em>{{ t.article_count }}</em></router-link>
      </div>

      <div class="tag-grid">
        <router-link
          v-for="(t, i) in tags" :key="t.id"
          :to="{ name: 'articles', query: { tag_id: t.id } }"
          class="tag-row"
        >
          <span class="tr-idx" :class="{ hot: i < 3 }">{{ String(i + 1).padStart(2, '0') }}</span>
          <span class="tr-name">{{ t.name }}</span>
          <span class="tr-bar"><i :style="{ width: bar(t.article_count) }"></i></span>
          <span class="tr-num">{{ t.article_count }}</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listTags } from '@/api/tag'

const loading = ref(true)
const tags = ref([])
const maxCount = ref(0)

// 热度条宽度：相对最高文章数的标签
function bar(count) {
  return maxCount.value > 0 ? Math.max(8, Math.round((count / maxCount.value) * 100)) + '%' : '8%'
}

onMounted(async () => {
  try {
    tags.value = await listTags()
    maxCount.value = tags.value.reduce((m, t) => Math.max(m, t.article_count || 0), 0)
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-title-bar {
  max-width: 1020px; margin: 10px auto 18px; padding: 0 0 14px;
  border-bottom: 1px solid var(--border);
}
.page-title-bar h1 { font-size: 28px; }
.page-title-bar .sub { font-size: 13px; color: var(--text-muted); margin-top: 4px; }

.tag-page-body { max-width: 1020px; margin: 0 auto; }

/* 顶部高频 chips */
.tag-head-strip {
  display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 16px;
  padding: 12px 14px; background: var(--surface);
  border: 1px solid var(--border); border-radius: var(--radius); box-shadow: var(--shadow-sm);
}
.top-chip {
  display: inline-flex; align-items: center; gap: 5px;
  padding: 5px 13px; border-radius: 999px; font-size: 13px;
  background: var(--primary-50); color: var(--primary-700);
  transition: background .15s;
}
.top-chip:hover { background: var(--primary-100); }
.top-chip em { font-style: normal; font-size: 11px; opacity: .7; }

/* 主区两列热度表 */
.tag-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px 14px; }
.tag-row {
  display: flex; align-items: center; gap: 12px;
  padding: 10px 16px; background: var(--surface);
  border: 1px solid var(--border); border-radius: 10px; box-shadow: var(--shadow-sm);
  transition: transform .15s, border-color .15s, box-shadow .15s;
}
.tag-row:hover { transform: translateY(-1px); border-color: var(--primary-200); box-shadow: var(--shadow); }
.tr-idx {
  flex-shrink: 0; font-family: var(--font-display); font-size: 13px; color: var(--text-faint);
  width: 22px;
}
.tr-idx.hot { color: var(--accent); font-weight: 700; }
.tr-name { flex-shrink: 0; font-size: 14.5px; font-weight: 600; color: #24345a; min-width: 76px; }
.tr-bar { flex: 1; height: 6px; border-radius: 3px; background: var(--surface-2); overflow: hidden; }
.tr-bar i { display: block; height: 100%; border-radius: 3px; background: linear-gradient(90deg, var(--primary), var(--accent)); }
.tr-num { flex-shrink: 0; min-width: 20px; text-align: right; font-size: 12px; color: var(--text-faint); font-variant-numeric: tabular-nums; }
@media (max-width: 720px) {
  .tag-grid { grid-template-columns: 1fr; }
}
</style>
