<template>
  <div>
    <div class="admin-page-head flex-between">
      <h2>数据统计</h2>
      <router-link to="/admin/articles/new" class="btn btn-primary btn-sm">✏️ 写文章</router-link>
    </div>

    <!-- 统计卡片 -->
    <div class="stat-grid">
      <div class="card stat-card">
        <div class="label">文章总数</div>
        <div class="value">{{ stats.article_count ?? '-' }}</div>
        <div class="sub">含草稿</div>
      </div>
      <div class="card stat-card">
        <div class="label">评论总数</div>
        <div class="value">{{ stats.comment_count ?? '-' }}</div>
        <div class="sub">全部评论</div>
      </div>
      <div class="card stat-card">
        <div class="label">注册用户</div>
        <div class="value">{{ stats.user_count ?? '-' }}</div>
        <div class="sub">全部用户</div>
      </div>
      <div class="card stat-card">
        <div class="label">总浏览量</div>
        <div class="value">{{ stats.total_views ?? '-' }}</div>
        <div class="sub">累计浏览</div>
      </div>
    </div>

    <!-- 最近文章 -->
    <div class="section-title">最近文章
      <router-link to="/admin/articles" class="more">管理全部 →</router-link>
    </div>
    <div class="admin-card">
      <el-table :data="recent" v-loading="loading" style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="220" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <span class="badge" :class="row.status === 2 ? 'badge-success' : 'badge-gray'">
              {{ row.status === 2 ? '已发布' : '草稿' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="浏览量" width="90" />
        <el-table-column label="发布时间" width="130">
          <template #default="{ row }">{{ formatDate(row.published_at || row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <router-link :to="`/admin/articles/${row.id}/edit`" class="op-link">编辑</router-link>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getStats } from '@/api/stats'
import { adminListArticles } from '@/api/article'
import { formatDate } from '@/utils/format'

const stats = ref({})
const recent = ref([])
const loading = ref(false)

onMounted(async () => {
  try {
    stats.value = await getStats()
  } catch (e) { /* 拦截器提示 */ }
  loading.value = true
  try {
    const data = await adminListArticles({ page: 1, page_size: 6 })
    recent.value = data.list || []
  } catch (e) { /* ignore */ } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.stat-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 18px; }
.stat-card { padding: 20px 22px; }
.stat-card .label { font-size: 13px; color: var(--text-muted); margin-bottom: 8px; }
.stat-card .value { font-size: 28px; font-weight: 700; line-height: 1.1; }
.stat-card .sub { font-size: 12px; margin-top: 8px; color: var(--text-faint); }
.section-title {
  font-size: 16px; font-weight: 700; margin: 26px 0 14px;
  display: flex; justify-content: space-between; align-items: center;
}
.section-title .more { font-size: 13px; color: var(--text-muted); font-weight: 400; }
.op-link { color: var(--primary); font-size: 13px; }
@media (max-width: 900px) {
  .stat-grid { grid-template-columns: repeat(2, 1fr); }
}
</style>
