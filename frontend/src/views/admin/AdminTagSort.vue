<template>
  <div>
    <router-link to="/admin/tags" class="back-link">← 返回标签列表</router-link>

    <div class="admin-page-head">
      <h2>标签排序</h2>
      <p class="sub">按前台「标签」模块展示位置排列，拖动调整顺序，松手自动保存并同步前台。</p>
    </div>

    <div class="admin-card sort-card" v-loading="loading">
      <template v-if="list.length">
        <div class="chip-cloud">
          <div
            v-for="(t, i) in list" :key="t.id"
            class="drag-chip" :class="{ dragging: dragIndex === i }"
            draggable="true"
            @dragstart="onDragStart(i)" @dragover.prevent="onDragOver" @drop="onDrop(i)"
          >
            <span class="idx">{{ i + 1 }}</span>
            <span class="label">{{ t.name }}</span>
            <span class="grip">⠿</span>
          </div>
        </div>
        <div class="save-bar">
          <span class="save-state" :class="{ done: !saving }">
            {{ saving ? '保存中…' : '顺序已自动保存' }}
          </span>
        </div>
      </template>
      <div v-else class="empty">
        <div class="big">🏷️</div>
        暂无标签，先去「标签列表」新建标签再来排序
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listTags, adminReorderTags } from '@/api/tag'

const list = ref([])
const loading = ref(false)
const saving = ref(false)
const dragIndex = ref(null)

async function load() {
  loading.value = true
  try {
    list.value = await listTags() // 后端已按 sort_order 升序返回
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
    await adminReorderTags(list.value.map(t => t.id))
    ElMessage.success('标签顺序已更新，前台已同步')
  } catch (e) { /* 拦截器提示；失败可重新拖动 */ } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.admin-page-head .sub { font-size: 13px; color: var(--text-muted); margin-top: 4px; }
.sort-card { padding: 20px 22px; }
.chip-cloud { display: flex; flex-wrap: wrap; gap: 12px; }
.drag-chip {
  display: inline-flex; align-items: center; gap: 9px;
  padding: 9px 16px; background: #fff; border: 1px solid var(--border);
  border-radius: 999px; cursor: grab; user-select: none; font-size: 14px;
  box-shadow: var(--shadow-sm); transition: box-shadow .15s, transform .05s, border-color .15s;
}
.drag-chip:hover { border-color: var(--primary); box-shadow: var(--shadow); }
.drag-chip:active { cursor: grabbing; }
.drag-chip.dragging { opacity: .35; transform: scale(.97); }
.drag-chip .idx {
  width: 22px; height: 22px; border-radius: 50%; background: var(--primary); color: #fff;
  font-size: 12.5px; display: grid; place-items: center; flex-shrink: 0;
}
.drag-chip .label { color: var(--text); }
.drag-chip .grip { color: var(--text-faint); font-size: 15px; }
.save-bar { margin-top: 16px; border-top: 1px dashed var(--border); padding-top: 12px; display: flex; justify-content: flex-end; }
.save-state { font-size: 12.5px; color: var(--text-faint); }
.save-state.done { color: var(--success); }
</style>
