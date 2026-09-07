<template>
  <div class="page-container page-block">
    <div class="page-title-bar">
      <h1>标签</h1>
      <div class="sub">按标签浏览文章</div>
    </div>

    <div v-if="loading" class="loading-box">加载中…</div>
    <div v-else-if="!tags.length" class="empty"><div class="big">🏷️</div>暂无标签</div>

    <div v-else class="card card-pad">
      <div class="tag-cloud">
        <router-link
          v-for="t in tags" :key="t.id"
          :to="{ name: 'articles', query: { tag_id: t.id } }"
          class="tag-chip tag-page"
        >{{ t.name }} <span class="n">({{ t.article_count }})</span></router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listTags } from '@/api/tag'

const loading = ref(true)
const tags = ref([])

onMounted(async () => {
  try {
    tags.value = await listTags()
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.tag-cloud { display: flex; flex-wrap: wrap; gap: 10px; }
.tag-page { padding: 6px 14px; font-size: 14px; }
.tag-page .n { color: var(--text-faint); font-size: 12px; margin-left: 4px; }
</style>
