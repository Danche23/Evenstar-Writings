<template>
  <div>
    <div class="admin-page-head">
      <h2>上传管理</h2>
    </div>

    <!-- 按场景筛选 -->
    <div class="toolbar-row">
      <el-select v-model="sceneFilter" placeholder="全部场景" clearable style="width: 160px" @change="reload">
        <el-option label="article（文章图）" value="article" />
        <el-option label="avatar（头像）" value="avatar" />
      </el-select>
    </div>

    <p class="table-hint">← 左右滑动表格，查看全部列与操作 →</p>
    <div class="admin-card">
      <el-table :data="list" v-loading="loading" style="width: 100%; min-width: 920px">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="预览" width="90">
          <template #default="{ row }">
            <el-image v-if="row.url" :src="row.url" fit="cover" style="width: 56px; height: 40px; border-radius: 6px"
              :preview-src-list="[row.url]" preview-teleported />
            <span v-else class="faint">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="filename" label="文件名" min-width="160" show-overflow-tooltip />
        <el-table-column label="场景" width="110">
          <template #default="{ row }">
            <span class="badge" :class="row.scene === 'article' ? 'badge-primary' : 'badge-gray'">{{ row.scene }}</span>
          </template>
        </el-table-column>
        <el-table-column label="大小" width="100">
          <template #default="{ row }">{{ formatSize(row.size) }}</template>
        </el-table-column>
        <el-table-column prop="user_id" label="上传者ID" width="100" />
        <el-table-column label="时间" width="130">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <span class="op-link" @click="copyLink(row)">复制链接</span>
            <el-popconfirm title="确定删除该文件？" width="200" @confirm="remove(row)">
              <template #reference><span class="op-link danger">删除</span></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="pager">
      <el-pagination
        layout="prev, pager, next, total"
        :total="total"
        :page-size="query.page_size"
        :current-page="query.page"
        @current-change="(p) => { query.page = p; load() }"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminListUploads, adminDeleteUpload } from '@/api/upload'
import { formatDateTime, formatSize } from '@/utils/format'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const sceneFilter = ref(null)
const query = reactive({ page: 1, page_size: 10 })

async function load() {
  loading.value = true
  try {
    const params = { page: query.page, page_size: query.page_size }
    if (sceneFilter.value) params.scene = sceneFilter.value
    const data = await adminListUploads(params)
    list.value = data.list || []
    total.value = Number(data.total || 0)
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
}

function reload() {
  query.page = 1
  load()
}

async function copyLink(row) {
  try {
    await navigator.clipboard.writeText(row.url)
    ElMessage.success('链接已复制')
  } catch (e) {
    ElMessage.error('复制失败，请手动复制：' + row.url)
  }
}

// 删除：被文章引用时后端返回 409 + 引用文章列表 → 弹确认后 force 删除
async function remove(row) {
  try {
    await adminDeleteUpload(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    if (e.code === 409) {
      const articles = e.data?.articles || []
      const names = articles.map((a) => `《${a.title}》（${a.status === 2 ? '已发布' : '草稿'}）`).join('\n')
      try {
        await ElMessageBox.confirm(
          `该文件被 ${articles.length} 篇文章引用，强制删除后图片将无法显示：\n${names}`,
          '文件被引用',
          { confirmButtonText: '强制删除', cancelButtonText: '取消', type: 'warning' }
        )
        await adminDeleteUpload(row.id, true)
        ElMessage.success('已强制删除')
        load()
      } catch (confirmErr) {
        if (confirmErr !== 'cancel' && confirmErr !== 'close') ElMessage.error(confirmErr.message || '删除失败')
      }
    } else {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar-row { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.op-link { color: var(--primary); font-size: 13px; margin-right: 12px; cursor: pointer; }
.op-link.danger { color: var(--danger); }
.pager { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
