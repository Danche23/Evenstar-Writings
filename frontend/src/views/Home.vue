<template>
  <div class="page-container page-block">
    <!-- 置顶展示最新一篇文章 -->
    <section v-if="hero" class="hero" @click="goArticle(hero)">
      <div class="hero-cat">{{ hero.categories?.[0]?.name || '最新发布' }}</div>
      <h1>{{ hero.title }}</h1>
      <p>{{ hero.summary }}</p>
      <div class="hero-meta">
        <span>{{ hero.author?.nickname || '灯影' }}</span>
        <span>{{ formatDate(hero.published_at || hero.created_at) }}</span>
        <span>👁 {{ hero.views ?? 0 }}</span>
        <a class="read" @click.stop="goArticle(hero)">阅读全文 →</a>
      </div>
    </section>

    <div v-if="loading" class="loading-box">加载中…</div>

    <div v-else class="home-cols">
      <!-- 左侧：最新文章 -->
      <div>
        <div class="section-title">
          最新文章
          <router-link to="/articles" class="more">查看全部 →</router-link>
        </div>
        <ArticleCard v-for="a in latestList" :key="a.id" :article="a" />
      </div>

      <!-- 右侧：热门 / 分类 / 标签 / 关于 -->
      <aside>
        <div v-if="hotList.length" class="side-card">
          <h4>热门文章</h4>
          <ol class="hot-list">
            <li v-for="a in hotList" :key="a.id">
              <router-link :to="`/articles/${a.id}`">{{ a.title }}</router-link>
              <span class="views">{{ a.views ?? 0 }}</span>
            </li>
          </ol>
        </div>

        <div v-if="categories.length" class="side-card">
          <h4>分类</h4>
          <div class="tag-cloud">
            <router-link
              v-for="c in categories" :key="c.id"
              :to="{ name: 'articles', query: { category_id: c.id } }"
              class="tag-chip"
            >{{ c.name }} ({{ c.article_count }})</router-link>
          </div>
        </div>

        <div v-if="tags.length" class="side-card">
          <h4>标签</h4>
          <div class="tag-cloud">
            <router-link
              v-for="t in tags" :key="t.id"
              :to="{ name: 'articles', query: { tag_id: t.id } }"
              class="tag-chip"
            >{{ t.name }}</router-link>
          </div>
        </div>

        <div class="side-card">
          <h4>关于</h4>
          <p class="about-text">一个热爱技术与记录的程序员，这里是我的个人博客，欢迎交流。</p>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ArticleCard from '@/components/ArticleCard.vue'
import { listArticles, getHotArticles } from '@/api/article'
import { listCategories } from '@/api/category'
import { listTags } from '@/api/tag'
import { formatDate } from '@/utils/format'

const router = useRouter()
const loading = ref(true)
const hero = ref(null)
const latestList = ref([])
const hotList = ref([])
const categories = ref([])
const tags = ref([])

function goArticle(a) {
  router.push(`/articles/${a.id}`)
}

onMounted(async () => {
  try {
    const [page, hot, cats, tgs] = await Promise.all([
      listArticles({ page: 1, page_size: 6 }),
      getHotArticles(5),
      listCategories(),
      listTags()
    ])
    hero.value = page.list?.[0] || null
    latestList.value = (page.list || []).slice(hero.value ? 1 : 0)
    hotList.value = hot || []
    categories.value = cats || []
    tags.value = (tgs || []).slice(0, 12)
  } catch (e) {
    /* 错误提示已由拦截器统一处理 */
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.hero {
  background: linear-gradient(135deg, #1e1b4b, #312e81 45%, #4f46e5);
  color: #fff; border-radius: 14px; padding: 38px 40px; margin: 8px 0 26px;
  position: relative; overflow: hidden; cursor: pointer;
}
.hero-cat { font-size: 12px; letter-spacing: 1px; text-transform: uppercase; color: #c7d2fe; margin-bottom: 10px; }
.hero h1 { font-size: 26px; line-height: 1.35; font-weight: 700; margin-bottom: 12px; max-width: 640px; }
.hero p { color: #c7d2fe; max-width: 600px; font-size: 14px; margin-bottom: 16px;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.hero-meta { font-size: 13px; color: #a5b4fc; display: flex; gap: 16px; align-items: center; }
.hero a.read { color: #fff; font-weight: 600; }
.home-cols { display: grid; grid-template-columns: 1fr 300px; gap: 24px; align-items: start; }
.section-title {
  font-size: 18px; font-weight: 700; margin: 6px 0 16px;
  display: flex; align-items: center; justify-content: space-between;
}
.section-title .more { font-size: 13px; font-weight: 500; color: var(--text-muted); }
.section-title .more:hover { color: var(--primary); }
.side-card { background: #fff; border: 1px solid var(--border); border-radius: var(--radius); padding: 18px 20px; margin-bottom: 14px; }
.side-card h4 { font-size: 15px; font-weight: 700; margin-bottom: 12px; }
.hot-list { list-style: none; counter-reset: hot; }
.hot-list li { display: flex; gap: 10px; padding: 8px 0; border-bottom: 1px dashed var(--border); font-size: 13.5px; align-items: center; }
.hot-list li:last-child { border-bottom: none; }
.hot-list li::before {
  counter-increment: hot; content: counter(hot);
  flex-shrink: 0; width: 20px; height: 20px; border-radius: 6px;
  background: #f3f4f6; color: var(--text-muted); font-size: 12px;
  display: grid; place-items: center; font-weight: 600;
}
.hot-list li:nth-child(-n+3)::before { background: var(--primary-50); color: var(--primary-600); }
.hot-list a { flex: 1; min-width: 0; color: var(--text);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.hot-list a:hover { color: var(--primary); }
.hot-list .views { color: var(--text-faint); font-size: 12px; white-space: nowrap; }
.tag-cloud { display: flex; flex-wrap: wrap; gap: 8px; }
.about-text { font-size: 13px; color: var(--text-muted); }
@media (max-width: 900px) {
  .home-cols { grid-template-columns: 1fr; }
}
</style>
