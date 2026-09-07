<template>
  <div class="page-container page-block">
    <div class="page-title-bar flex-between">
      <div>
        <h1>{{ title }}</h1>
        <div class="sub">{{ total > 0 ? `共 ${total} 篇` : ' ' }}</div>
      </div>
      <!-- 筛选：分类 + 关键词 -->
      <div class="search-row">
        <el-select v-model="query.category_id" placeholder="全部分类" clearable style="width: 150px" @change="onFilter">
          <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
        </el-select>
        <el-input v-model="keyword" placeholder="搜索标题 / 摘要…" clearable style="width: 220px"
          @keyup.enter="onSearch" @clear="onSearch" />
        <button class="btn btn-primary" @click="onSearch">搜索</button>
      </div>
    </div>

    <div v-if="loading" class="loading-box">加载中…</div>
    <div v-else-if="!list.length" class="empty">
      <div class="big">📭</div>
      暂无文章
    </div>
    <template v-else>
      <ArticleCard v-for="a in list" :key="a.id" :article="a" />
      <div class="pager" v-if="totalPage > 1">
        <el-pagination
          layout="prev, pager, next, total"
          :total="total"
          :page-size="query.page_size"
          :current-page="query.page"
          @current-change="onPage"
        />
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ArticleCard from '@/components/ArticleCard.vue'
import { listArticles } from '@/api/article'
import { listCategories } from '@/api/category'

const route = useRoute()
const loading = ref(true)
const list = ref([])
const total = ref(0)
const totalPage = ref(1)
const categories = ref([])
const keyword = ref('')

const query = reactive({ page: 1, page_size: 10, category_id: null, tag_id: null })

const title = computed(() => {
  if (query.category_id) return categories.value.find((c) => c.id === Number(query.category_id))?.name || '文章列表'
  if (query.tag_id) return '标签相关文章'
  if (keyword.value) return `搜索：${keyword.value}`
  return '全部文章'
})

async function load() {
  loading.value = true
  try {
    const params = { page: query.page, page_size: query.page_size }
    if (query.category_id) params.category_id = query.category_id
    if (query.tag_id) params.tag_id = query.tag_id
    if (keyword.value) params.keyword = keyword.value
    const data = await listArticles(params)
    list.value = data.list || []
    total.value = Number(data.total || 0)
    totalPage.value = data.total_page || 1
  } catch (e) {
    /* 拦截器已提示 */
  } finally {
    loading.value = false
  }
}

function onPage(p) {
  query.page = p
  load()
}

function onFilter() {
  query.page = 1
  load()
}

function onSearch() {
  query.page = 1
  load()
}

// 从分类页 / 标签页带参数跳转过来时同步筛选
watch(
  () => route.query,
  (q) => {
    query.category_id = q.category_id ? Number(q.category_id) : null
    query.tag_id = q.tag_id ? Number(q.tag_id) : null
    if (q.keyword) keyword.value = String(q.keyword)
    query.page = 1
    load()
  },
  { immediate: true }
)

onMounted(async () => {
  try {
    categories.value = await listCategories()
  } catch (e) { /* ignore */ }
})
</script>

<style scoped>
.pager { display: flex; justify-content: center; margin: 22px 0 8px; }
</style>
