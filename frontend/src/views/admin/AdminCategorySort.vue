<template>
  <div>
    <router-link to="/admin/categories" class="back-link">← 返回分类列表</router-link>

    <div class="admin-page-head">
      <h2>分类排序</h2>
      <p class="sub">按前台「分类」模块展示位置排列（与前台列表同款式），拖动行调整顺序，松手自动保存并同步前台。</p>
    </div>

    <div class="admin-card sort-card" v-loading="loading">
      <template v-if="list.length">
        <div class="sort-list">
          <div
            v-for="(c, i) in list" :key="c.id"
            class="sort-row" :class="{ dragging: dragIndex === i }"
            draggable="true"
            @dragstart="onDragStart(i)"
            @dragover.prevent="onDragOver"
            @drop="onDrop(i)"
            @dragend="dragIndex = null"
          >
            <span class="badge">{{ c.name.slice(0, 1) }}</span>
            <span class="name">{{ c.name }}</span>
            <span class="count">{{ c.article_count }} 篇</span>
            <span class="grip" title="拖动调整顺序">⠿</span>
          </div>
        </div>
        <div class="save-bar">
          <span class="save-state" :class="{ done: !saving }">
            {{ saving ? '保存中…' : '顺序已自动保存，前台已同步' }}
          </span>
        </div>
      </template>
      <div v-else class="empty">
        <div class="big">🗂️</div>
        暂无分类，先去「分类列表」新建分类再来排序
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listCategories, adminReorderCategories } from '@/api/category'

const list = ref([])
const loading = ref(false)
const saving = ref(false)
const dragIndex = ref(null)

async function load() {
  loading.value = true
  try {
    list.value = await listCategories() // 后端已按 sort_order 升序返回
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
}

function onDragStart(i) { dragIndex.value = i }
function onDragOver() {}
function onDrop(i) {
  const from = dragIndex.value
  dragIndex.value = null
  if (from === null || from === i) return
  const arr = list.value
  const item = arr.splice(from, 1)[0]
  arr.splice(i, 0, item)
  saveOrder()
}

async function saveOrder() {
  saving.value = true
  try {
    await adminReorderCategories(list.value.map(c => c.id))
    ElMessage.success('分类顺序已更新，前台已同步')
  } catch (e) { /* 拦截器提示；失败可重新拖动 */ } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.admin-page-head .sub { font-size: 13px; color: var(--text-muted); margin-top: 4px; }
.sort-card { padding: 18px 22px; }

/* 排序列表：复刻前台侧栏「分类」列表式样式 */
.sort-list { display: flex; flex-direction: column; gap: 6px; }
.sort-row {
  display: flex; align-items: center; gap: 12px;
  padding: 9px 12px; background: #fff;
  border: 1px solid var(--border); border-radius: 10px;
  cursor: grab; user-select: none; color: var(--text);
  transition: border-color .15s, box-shadow .15s, opacity .15s, transform .1s;
}
.sort-row:hover { border-color: var(--primary); box-shadow: var(--shadow-sm); }
.sort-row:active { cursor: grabbing; }
.sort-row.dragging { opacity: .35; transform: scale(.98); }
.sort-row .badge {
  width: 26px; height: 26px; border-radius: 50%; flex-shrink: 0;
  background: var(--primary-100); color: var(--primary-700);
  display: grid; place-items: center; font-size: 13px; font-weight: 600;
}
.sort-row:nth-child(even) .badge { background: #f3e6c8; color: #8a6d2f; }
.sort-row .name {
  flex: 1; min-width: 0; font-size: 14px;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.sort-row .count {
  font-size: 11.5px; color: var(--text-faint); background: var(--bg);
  padding: 2px 9px; border-radius: 999px; flex-shrink: 0;
}
.sort-row .grip { color: var(--text-faint); font-size: 15px; cursor: grab; flex-shrink: 0; }

.save-bar { margin-top: 16px; border-top: 1px dashed var(--border); padding-top: 12px; display: flex; justify-content: flex-end; }
.save-state { font-size: 12.5px; color: var(--text-faint); }
.save-state.done { color: var(--success); }
</style>
