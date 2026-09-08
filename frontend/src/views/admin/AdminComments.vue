<template>
  <div>
    <div class="admin-page-head">
      <h2>评论管理</h2>
    </div>

    <!-- 按文章标题 / 用户模糊检索 -->
    <div class="toolbar-row">
      <el-input v-model="titleKeyword" placeholder="搜索文章标题…" clearable style="width: 220px"
        @keyup.enter="reload" @clear="reload" />
      <el-input v-model="userKeyword" placeholder="搜索用户（昵称/邮箱）…" clearable style="width: 220px"
        @keyup.enter="reload" @clear="reload" />
      <button class="btn btn-primary btn-sm" @click="reload">搜索</button>
    </div>

    <p class="table-hint">← 左右滑动表格，查看全部列与操作 →</p>
    <div class="admin-card">
      <el-table :data="list" v-loading="loading" style="width: 100%; min-width: 1080px">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="article_title" label="所属文章" min-width="160" show-overflow-tooltip />
        <el-table-column label="昵称" width="120">
          <template #default="{ row }">{{ row.user_nickname || '已注销用户' }}</template>
        </el-table-column>
        <el-table-column label="等级" width="90">
          <template #default="{ row }">
            <span class="badge" :class="row.parent_id ? 'badge-gray' : 'badge-primary'">
              {{ row.parent_id ? '二级' : '一级' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" min-width="220" show-overflow-tooltip />
        <el-table-column label="置顶" width="90">
          <template #default="{ row }">
            <span class="badge" :class="row.is_top === 1 ? 'badge-primary' : 'badge-gray'">
              {{ row.is_top === 1 ? '已置顶' : '未置顶' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="时间" width="120">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <div class="op-cell">
              <div class="op-main">
                <template v-if="!row.parent_id">
                  <span class="op-link" @click="toggleTop(row)">{{ row.is_top === 1 ? '取消置顶' : '置顶' }}</span>
                </template>
                <span v-else class="op-faint">二级不可置顶</span>
              </div>
              <div class="op-del">
                <el-popconfirm title="确定删除该评论？" width="200" @confirm="remove(row)">
                  <template #reference><span class="op-link danger">删除</span></template>
                </el-popconfirm>
              </div>
            </div>
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
import { ElMessage } from 'element-plus'
import { adminListComments, adminDeleteComment, adminToggleCommentTop } from '@/api/comment'
import { formatDateTime } from '@/utils/format'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const titleKeyword = ref('')
const userKeyword = ref('')
const query = reactive({ page: 1, page_size: 10 })

async function load() {
  loading.value = true
  try {
    const params = { page: query.page, page_size: query.page_size }
    if (titleKeyword.value.trim()) params.keyword = titleKeyword.value.trim()
    if (userKeyword.value.trim()) params.user = userKeyword.value.trim()
    const data = await adminListComments(params)
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

async function toggleTop(row) {
  try {
    await adminToggleCommentTop(row.id)
    ElMessage.success(row.is_top === 1 ? '已取消置顶' : '已置顶')
    load()
  } catch (e) { /* 拦截器提示（如二级评论置顶返回 1012） */ }
}

async function remove(row) {
  try {
    await adminDeleteComment(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) { /* 拦截器提示 */ }
}

onMounted(load)
</script>

<style scoped>
.toolbar-row { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.op-cell { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 12px; align-items: center; }
.op-main { min-width: 0; }
.op-del { text-align: right; }
.op-link { color: var(--primary); font-size: 13px; cursor: pointer; white-space: nowrap; }
.op-link.danger { color: var(--danger); }
.op-faint { color: var(--text-faint); font-size: 12px; white-space: nowrap; }
.pager { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
