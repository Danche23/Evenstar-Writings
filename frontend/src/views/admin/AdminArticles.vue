<template>
  <div>
    <div class="admin-page-head flex-between">
      <h2>文章管理</h2>
      <router-link to="/admin/articles/new" class="btn btn-primary btn-sm">＋ 新建文章</router-link>
    </div>

    <!-- 筛选 -->
    <div class="toolbar-row">
      <el-select v-model="query.status" placeholder="全部状态" clearable style="width: 140px" @change="reload">
        <el-option label="草稿" :value="1" />
        <el-option label="已发布" :value="2" />
      </el-select>
      <el-input v-model="keyword" placeholder="搜索标题…" clearable style="width: 240px"
        @keyup.enter="reload" @clear="reload" />
      <button class="btn btn-primary btn-sm" @click="reload">搜索</button>
    </div>

    <div class="admin-card">
      <el-table :data="list" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="240" show-overflow-tooltip>
          <template #default="{ row }">
            <router-link :to="`/articles/${row.id}`" target="_blank" class="title-link">{{ row.title }}</router-link>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <span class="badge" :class="row.status === 2 ? 'badge-success' : 'badge-gray'">
              {{ row.status === 2 ? '已发布' : '草稿' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="分类" width="140">
          <template #default="{ row }">
            <span v-for="c in row.categories" :key="c.id" class="muted small" style="margin-right: 6px">{{ c.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="浏览" width="80" />
        <el-table-column label="发布时间" width="120">
          <template #default="{ row }">
            {{ row.published_at ? formatDate(row.published_at) : '—' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <router-link :to="`/admin/articles/${row.id}/edit`" class="op-link">编辑</router-link>
            <el-popconfirm title="确定删除该文章？" width="200" @confirm="remove(row)">
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
import { ElMessage } from 'element-plus'
import { adminListArticles, adminDeleteArticle } from '@/api/article'
import { formatDate } from '@/utils/format'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const keyword = ref('')
const query = reactive({ page: 1, page_size: 10, status: null })

async function load() {
  loading.value = true
  try {
    const params = { page: query.page, page_size: query.page_size }
    if (query.status) params.status = query.status
    if (keyword.value.trim()) params.keyword = keyword.value.trim()
    const data = await adminListArticles(params)
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

async function remove(row) {
  try {
    await adminDeleteArticle(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) { /* 拦截器提示 */ }
}

onMounted(load)
</script>

<style scoped>
.toolbar-row { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.title-link { color: var(--text); font-weight: 500; }
.title-link:hover { color: var(--primary); }
.op-link { color: var(--primary); font-size: 13px; margin-right: 12px; cursor: pointer; }
.op-link.danger { color: var(--danger); }
.pager { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
