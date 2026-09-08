<template>
  <!-- 博主社交链接：有地址才显示，空则不占布局 -->
  <div v-if="visible.length" class="social-links" :class="{ center }">
    <a
      v-for="s in visible" :key="s.key" :href="s.url" target="_blank" rel="noopener"
      class="social-btn" :class="s.key" :aria-label="s.name"
    >
      <!-- GitHub 图标 -->
      <svg v-if="s.key === 'github'" viewBox="0 0 16 16" width="16" height="16" fill="currentColor" aria-hidden="true">
        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82a7.7 7.7 0 0 1 2-.27c.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8z"/>
      </svg>
      <span v-else class="csdn-word">CSDN</span>
      <span class="lbl">{{ s.name }}</span>
    </a>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { authorLinks } from '@/utils/site'

const props = defineProps({
  center: { type: Boolean, default: false }
})

const visible = computed(() => {
  const out = []
  if (authorLinks.github) out.push({ key: 'github', name: 'GitHub', url: authorLinks.github })
  if (authorLinks.csdn) out.push({ key: 'csdn', name: 'CSDN', url: authorLinks.csdn })
  return out
})
</script>

<style scoped>
.social-links { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 10px; }
.social-links.center { justify-content: center; }
.social-btn {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 5px 12px; border-radius: 999px; font-size: 12.5px;
  border: 1px solid var(--border); background: var(--surface); color: var(--text-muted);
  transition: color .15s, border-color .15s, background .15s;
}
.social-btn:hover { color: var(--primary); border-color: var(--primary); background: var(--primary-50); }
.social-btn .lbl { font-weight: 500; }
.social-btn .csdn-word { font-size: 11px; letter-spacing: .02em; color: #c5542a; font-weight: 700; }
.social-btn.csdn:hover { color: #c5542a; border-color: #d98a6b; background: #fdf3ef; }
</style>
