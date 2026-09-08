<template>
  <div class="page-container page-block">
    <header class="page-title-bar">
      <h1>归档</h1>
      <div class="sub">{{ total > 0 ? `共 ${total} 篇文章，按时间回溯` : ' ' }}</div>
    </header>

    <div v-if="loading" class="loading-box">翻找旧信札中…</div>
    <div v-else-if="!groups.length" class="empty"><div class="big">✉️</div>还没有写过一篇文章</div>

    <div v-else class="archive">
      <section v-for="g in groups" :key="g.month" class="arch-month">
        <div class="arch-month-head">
          <b>{{ g.month }}</b>
          <span>{{ g.items.length }} 篇</span>
        </div>
        <ul class="arch-list">
          <li v-for="a in g.items" :key="a.id">
            <span class="arch-date">{{ formatDate(a.published_at || a.created_at) }}</span>
            <router-link :to="`/articles/${a.id}`" class="arch-title">{{ a.title }}</router-link>
            <span v-if="a.categories?.[0]" class="arch-cat">{{ a.categories[0].name }}</span>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { listArticles } from '@/api/article'
import { formatDate } from '@/utils/format'

const loading = ref(true)
const list = ref([])
const total = ref(0)

// 分组：YYYY-MM（列表接口按时间倒序返回）
const groups = computed(() => {
  const map = new Map()
  for (const a of list.value) {
    const d = a.published_at || a.created_at
    const key = d ? String(d).slice(0, 7) : '未标注'
    if (!map.has(key)) map.set(key, [])
    map.get(key).push(a)
  }
  return [...map.entries()].map(([month, items]) => ({ month, items }))
})

onMounted(async () => {
  try {
    const size = 50
    let page = 1
    const all = []
    // 逐页拉取直至取完（小站文章量小，最多保护 30 页）
    for (; page <= 30; page++) {
      const data = await listArticles({ page, page_size: size })
      all.push(...(data.list || []))
      total.value = Number(data.total || 0)
      if (all.length >= total.value) break
    }
    list.value = all
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-title-bar {
  max-width: 900px; margin: 10px auto 26px; padding-bottom: 14px;
  border-bottom: 1px solid var(--border);
}
.page-title-bar h1 { font-size: 30px; }
.page-title-bar .sub { font-size: 13px; color: var(--text-muted); margin-top: 4px; }
.archive { max-width: 900px; margin: 0 auto; }
.arch-month { margin-bottom: 30px; }
.arch-month-head {
  display: flex; align-items: baseline; gap: 12px; margin-bottom: 8px;
  padding-left: 10px; position: relative;
}
.arch-month-head b {
  font-family: var(--font-display); font-size: 21px; font-weight: 700; color: #24345a;
}
.arch-month-head b::before { content: "❀ "; color: var(--accent); font-size: 15px; }
.arch-month-head span { font-size: 12px; color: var(--text-faint); }
.arch-list { list-style: none; border-left: 1px solid var(--border); margin-left: 10px; padding-left: 22px; position: relative; }
.arch-list::before {
  content: ""; position: absolute; left: -4px; top: 14px; width: 7px; height: 7px;
  border-radius: 50%; background: var(--accent);
}
.arch-list li {
  display: flex; align-items: baseline; gap: 14px; padding: 7px 0; font-size: 14px;
}
.arch-date {
  flex-shrink: 0; font-size: 12.5px; color: var(--text-faint);
  font-variant-numeric: tabular-nums; width: 84px;
}
.arch-title { flex: 1; min-width: 0; color: var(--text); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.arch-title:hover { color: var(--primary); }
.arch-cat { flex-shrink: 0; font-size: 12px; color: var(--text-faint); background: var(--surface-2); padding: 1px 9px; border-radius: 999px; }
@media (max-width: 640px) {
  .arch-cat { display: none; }
  .arch-date { width: 78px; }
}
</style>
