<template>
  <div>
    <div class="admin-page-head flex-between">
      <h2>标签列表</h2>
      <div class="head-actions">
        <router-link to="/admin/tags/sort" class="btn btn-outline btn-sm">↕ 排序布局</router-link>
        <button class="btn btn-primary btn-sm" @click="openDialog()">＋ 新建标签</button>
      </div>
    </div>

    <p class="table-hint">← 左右滑动表格，查看全部列与操作 →</p>
    <div class="admin-card">
      <el-table :data="list" v-loading="loading" style="width: 100%; min-width: 720px">
        <el-table-column type="index" label="排名" width="80" :index="rankIndex" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="160" />
        <el-table-column prop="article_count" label="文章数" width="100" />
        <el-table-column label="创建时间" width="140">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <span class="op-link" @click="openDialog(row)">编辑</span>
            <el-popconfirm title="删除后该标签关联将清空，确定？" width="220" @confirm="remove(row)">
              <template #reference><span class="op-link danger">删除</span></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 新建 / 编辑 -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑标签' : '新建标签'" width="400px">
      <el-form label-position="top" @submit.prevent>
        <el-form-item label="名称">
          <el-input v-model="name" maxlength="50" placeholder="标签名称" @keyup.enter="submit" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="btn btn-ghost" @click="dialogVisible = false">取消</button>
        <button class="btn btn-primary" :disabled="savingDlg" @click="submit">{{ savingDlg ? '保存中…' : '保存' }}</button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listTags, adminCreateTag, adminUpdateTag, adminDeleteTag } from '@/api/tag'
import { formatDateTime } from '@/utils/format'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(null)
const name = ref('')
const savingDlg = ref(false)

// 排名 = 列表序号（列表已按 sort_order 升序返回，即前台展示顺序）
function rankIndex(i) { return i + 1 }

async function load() {
  loading.value = true
  try {
    list.value = await listTags()
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
}

function openDialog(row) {
  editingId.value = row?.id || null
  name.value = row?.name || ''
  dialogVisible.value = true
}

async function submit() {
  if (!name.value.trim()) {
    ElMessage.warning('请输入标签名称')
    return
  }
  savingDlg.value = true
  try {
    if (editingId.value) await adminUpdateTag(editingId.value, name.value.trim())
    else await adminCreateTag(name.value.trim())
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } catch (e) { /* 拦截器提示 */ } finally {
    savingDlg.value = false
  }
}

async function remove(row) {
  try {
    await adminDeleteTag(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) { /* 拦截器提示 */ }
}

onMounted(load)
</script>

<style scoped>
.head-actions { display: flex; align-items: center; gap: 12px; }
.op-link { color: var(--primary); font-size: 13px; margin-right: 12px; cursor: pointer; }
.op-link.danger { color: var(--danger); }
</style>
