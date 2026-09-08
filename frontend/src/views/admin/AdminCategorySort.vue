<template>
  <div>
    <router-link to="/admin/categories" class="back-link">← 返回分类列表</router-link>

    <div class="admin-page-head">
      <h2>分类排序</h2>
      <p class="sub">按前台「分类」模块展示位置排列（与前台列表同款式，方形编号卡片），拖动卡片调整顺序，松手自动保存并同步前台。</p>
    </div>

    <div class="admin-card sort-card" v-loading="loading">
      <template v-if="list.length">
        <div class="cat-sort-grid">
          <div
            v-for="(c, i) in list" :key="c.id"
            class="cat-tile" :class="{ dragging: dragIndex === i }"
            draggable="true"
            @dragstart="onDragStart(i)"
            @dragover.prevent="onDragOver"
            @drop="onDrop(i)"
            @dragend="dragIndex = null"
          >
            <span class="tile-idx">{{ String(i + 1).padStart(2, '0') }}</span>
            <div class="tile-body">
              <span class="tile-name">{{ c.name }}</span>
              <span class="tile-count">{{ c.article_count }} 篇</span>
            </div>
            <div class="tile-foot">
              <span class="tile-hint">拖动排序</span>
              <span class="grip" title="拖动调整顺序">⠿</span>
            </div>
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
.sort-card { padding: 20px 22px; }

/* 排序栅格：与前台「分类」页完全一致 —— 桌面 3 列 / 平板 2 列 / 手机 1 列，
   前台断点：>900px 3列 → ≤900px 2列 → ≤600px 1列（Categories.vue） */
.cat-sort-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  max-width: 1000px;
}
.cat-tile {
  position: relative; overflow: hidden; min-width: 0;
  display: flex; flex-direction: column; justify-content: space-between;
  min-height: 118px; padding: 18px 20px;
  background: #fff; border: 1px solid var(--border); border-radius: 12px;
  cursor: grab; user-select: none; color: var(--text);
  box-shadow: var(--shadow-sm);
  transition: border-color .15s, box-shadow .15s, opacity .15s, transform .1s;
}
.cat-tile:hover { border-color: var(--primary-200); box-shadow: var(--shadow); }
.cat-tile:active { cursor: grabbing; }
.cat-tile.dragging { opacity: .35; transform: scale(.98); }
.tile-idx {
  position: absolute; top: 4px; right: 12px;
  font-family: var(--font-display); font-size: 40px; font-weight: 700;
  color: var(--primary-100); line-height: 1; pointer-events: none;
}
.tile-body { position: relative; display: flex; flex-direction: column; gap: 4px; }
.tile-name {
  font-family: var(--font-display); font-size: 19px; font-weight: 700;
  color: #2c261d; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.tile-count { font-size: 12.5px; color: var(--text-faint); }
.tile-foot {
  position: relative; display: flex; align-items: center; justify-content: space-between;
  margin-top: 10px; padding-top: 10px; border-top: 1px dashed var(--border);
}
.tile-hint { font-size: 12px; color: var(--text-faint); }
.cat-tile .grip { color: var(--text-faint); font-size: 16px; cursor: grab; flex-shrink: 0; }

.save-bar { margin-top: 18px; border-top: 1px dashed var(--border); padding-top: 12px; display: flex; justify-content: flex-end; }
.save-state { font-size: 12.5px; color: var(--text-faint); }
.save-state.done { color: var(--success); }

@media (max-width: 900px) {
  .cat-sort-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 600px) {
  .sort-card { padding: 14px; }
  .cat-sort-grid { grid-template-columns: 1fr; }
}
</style>
