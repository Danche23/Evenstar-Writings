<template>
  <div class="page-container page-block">
    <!-- 首屏卡：浅色花境，文案以技术博客为主 -->
    <section class="poem-hero rise">
      <svg class="ph-art" viewBox="0 0 1200 420" preserveAspectRatio="xMidYMid slice" aria-hidden="true">
        <defs>
          <linearGradient id="phsky" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#faf6ec"/>
            <stop offset="60%" stop-color="#f3eef5"/>
            <stop offset="100%" stop-color="#ece7f1"/>
          </linearGradient>
          <radialGradient id="phmoon" cx="50%" cy="50%" r="50%">
            <stop offset="0%" stop-color="#fff6d8"/><stop offset="100%" stop-color="#eccf8f"/>
          </radialGradient>
        </defs>
        <rect width="1200" height="420" fill="url(#phsky)"/>
        <circle cx="1005" cy="112" r="44" fill="url(#phmoon)" opacity=".92"/>
        <circle cx="1005" cy="112" r="60" fill="none" stroke="#e6cf92" stroke-opacity=".5" stroke-width="1"/>
        <g fill="#e2c77f" opacity=".85">
          <circle cx="976" cy="64" r="2.4"/><circle cx="1076" cy="88" r="2"/><circle cx="1142" cy="158" r="2.4"/>
          <circle cx="928" cy="150" r="1.8"/><circle cx="112" cy="76" r="2.2"/><circle cx="306" cy="56" r="1.9"/><circle cx="520" cy="96" r="1.6"/>
        </g>
        <path d="M0 340 C 240 300, 430 348, 640 324 S 980 300, 1200 336 L 1200 420 L 0 420 Z" fill="#e6e8f3" opacity=".75"/>
        <path d="M0 376 C 340 348, 600 390, 880 366 S 1150 354, 1200 372 L 1200 420 L 0 420 Z" fill="#dfe3ef"/>
        <g stroke="#6f8065" stroke-width="2" fill="none" stroke-linecap="round">
          <path d="M86 402 C 140 352, 206 300, 268 238"/>
          <path d="M150 392 C 210 348, 268 314, 326 280"/>
          <path d="M212 402 C 258 368, 300 346, 352 318"/>
          <path d="M176 342 C 204 326, 236 322, 262 332" stroke-width="1.4"/>
          <path d="M234 310 C 260 296, 290 292, 316 302" stroke-width="1.4"/>
        </g>
        <g>
          <g transform="translate(268 238)">
            <circle r="16" fill="#f9dfd7"/><circle r="11" fill="#f7cabb"/><circle r="5.4" fill="#eea488"/>
            <circle cx="-6" cy="-7" r="4" fill="#fdf6f2"/><circle cx="7" cy="6" r="4" fill="#fdf6f2"/>
          </g>
          <g transform="translate(326 280)">
            <circle r="12.5" fill="#f9dfd7"/><circle r="8.6" fill="#f7cabb"/><circle r="4.2" fill="#e8936f"/>
            <circle cx="-4" cy="-5" r="3.2" fill="#fdf6f2"/><circle cx="5" cy="5" r="3.2" fill="#fdf6f2"/>
          </g>
          <g transform="translate(352 318)"><circle r="9.5" fill="#e3e9fa"/><circle r="6.2" fill="#d2dbf4"/><circle r="3" fill="#c3a05c"/></g>
          <g transform="translate(262 332)"><circle r="8.5" fill="#f6e2ea"/><circle r="5.4" fill="#f0cad8"/><circle r="2.6" fill="#c3a05c"/></g>
          <g transform="translate(316 302)"><circle r="7.5" fill="#f5e6c9"/><circle r="4.8" fill="#eed7a7"/><circle r="2.4" fill="#b98d4a"/></g>
          <g transform="translate(150 306)"><circle r="6.5" fill="#f6e2ea"/><circle r="4" fill="#f0cad8"/><circle r="2" fill="#c3a05c"/></g>
        </g>
      </svg>
      <div class="ph-inner">
        <p class="ph-eyebrow">EVENSTAR WRITINGS · 技术博客</p>
        <h1 class="ph-title">把想明白的事，<br /><span>一个一个写下来。</span></h1>
        <p class="ph-poem poem">记录技术、读书与生活——想写的，都会出现在这里。</p>
        <div class="ph-actions">
          <router-link to="/articles" class="ph-cta">去看文章 <span class="arr">→</span></router-link>
          <button class="ph-cta ghost" @click="randomArticle">❀ 逛一篇</button>
        </div>
      </div>
    </section>

    <!-- 技术主题入口 -->
    <div v-if="!loading && techCats.length" class="tech-top rise rise-1">
      <router-link
        v-for="c in techCats" :key="c.id"
        :to="{ name: 'articles', query: { category_id: c.id } }"
        class="tech-tile"
      >
        <span class="t-badge">{{ c.name.slice(0, 1) }}</span>
        <span class="t-name">{{ c.name }}</span>
        <span class="t-num">{{ c.article_count }} 篇</span>
      </router-link>
    </div>

    <div v-if="loading" class="loading-box">加载中…</div>

    <div v-else class="home-cols">
      <!-- 左：最新文章 -->
      <div class="main-col rise rise-2">
        <div class="section-title">
          <span class="tt">最新文章</span>
          <div class="title-acts">
            <button class="more jump" @click="randomArticle">❀ 逛一篇</button>
            <router-link to="/articles" class="more">全部文章 →</router-link>
          </div>
        </div>
        <ArticleCard v-for="a in latestList" :key="a.id" :article="a" />
      </div>

      <!-- 右：侧栏 -->
      <aside class="side rise rise-3">
        <div v-if="hotList.length" class="side-card">
          <h4 class="side-title">热门文章</h4>
          <ol class="hot-list">
            <li v-for="a in hotList" :key="a.id">
              <router-link :to="`/articles/${a.id}`">{{ a.title }}</router-link>
              <span class="views">{{ a.views ?? 0 }}</span>
            </li>
          </ol>
        </div>

        <div v-if="categories.length" class="side-card">
          <h4 class="side-title">分类</h4>
          <div class="cat-list">
            <router-link
              v-for="c in categories" :key="c.id"
              :to="{ name: 'articles', query: { category_id: c.id } }"
              class="cat-item"
            >
              <span class="ci-badge">{{ c.name.slice(0, 1) }}</span>
              <span class="ci-name">{{ c.name }}</span>
              <span class="ci-num">{{ c.article_count }}</span>
              <span class="ci-arr">→</span>
            </router-link>
          </div>
        </div>

        <div v-if="tags.length" class="side-card">
          <h4 class="side-title">标签</h4>
          <div class="tag-cloud">
            <router-link
              v-for="t in tags" :key="t.id"
              :to="{ name: 'articles', query: { tag_id: t.id } }"
              class="tag-chip flower"
            >{{ t.name }} <em>{{ t.article_count }}</em></router-link>
          </div>
        </div>

        <div class="side-card author-card">
          <h4 class="side-title">关于博主</h4>
          <div class="author-body">
            <img v-if="author.avatar" :src="author.avatar" class="author-avatar" alt="" />
            <div v-else class="author-avatar author-avatar-text">{{ (author.nickname || '暮').slice(0, 1) }}</div>
            <div class="author-meta">
              <div class="author-name">{{ author.nickname || '暮星' }}</div>
              <p class="author-bio">{{ author.bio || '写代码的后端开发者，业余记录与分享。' }}</p>
              <SocialLinks />
              <router-link to="/about" class="author-more">了解博主 →</router-link>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import ArticleCard from '@/components/ArticleCard.vue'
import SocialLinks from '@/components/SocialLinks.vue'
import { listArticles, getHotArticles } from '@/api/article'
import { listCategories } from '@/api/category'
import { listTags } from '@/api/tag'
import { getAuthor } from '@/api/user'
import { formatDate } from '@/utils/format'

const router = useRouter()
const loading = ref(true)
const hero = ref(null)
const latestList = ref([])
const hotList = ref([])
const categories = ref([])
const tags = ref([])
const author = ref({ nickname: '', avatar: '', bio: '' })
const postsTotal = ref(0)

const stats = computed(() => ({
  posts: postsTotal.value,
  categories: categories.value.length,
  tags: tags.value.length
}))

// 技术主题（生活/读书类不占首页主体）
const techCats = computed(() =>
  categories.value.filter((c) => !['生活随笔', '读书笔记'].includes(c.name))
)

function goArticle(a) {
  router.push(`/articles/${a.id}`)
}

// 随机翻一篇：先拿总数，再随机页取一篇
async function randomArticle() {
  try {
    const first = await listArticles({ page: 1, page_size: 1 })
    const tp = Number(first.total_page || 1)
    if (tp < 1) {
      ElMessage.info('还没有文章可翻')
      return
    }
    const rp = Math.min(tp, Math.max(1, Math.floor(Math.random() * tp) + 1))
    const got = await listArticles({ page: rp, page_size: 1 })
    const item = got.list?.[0]
    if (item) goArticle(item)
  } catch (e) { /* 拦截器提示 */ }
}

onMounted(async () => {
  try {
    const [page, hot, cats, tgs, authorResp] = await Promise.all([
      listArticles({ page: 1, page_size: 6 }),
      getHotArticles(5),
      listCategories(),
      listTags(),
      getAuthor()
    ])
    postsTotal.value = Number(page.total || 0)
    latestList.value = page.list || []
    hotList.value = hot || []
    categories.value = cats || []
    tags.value = (tgs || []).slice(0, 12)
    if (authorResp) author.value = authorResp
  } catch (e) {
    /* 错误提示已由拦截器统一处理 */
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
/* —— 首屏卡：浅色花境插画（装饰低调，文案技术向） —— */
.poem-hero {
  position: relative; overflow: hidden;
  border-radius: 20px; margin: 6px 0 22px;
  border: 1px solid var(--border); box-shadow: var(--shadow);
}
.ph-art { position: absolute; inset: 0; width: 100%; height: 100%; }
.ph-inner { position: relative; padding: 40px 44px 38px; max-width: 780px; }
.ph-eyebrow { font-size: 12px; letter-spacing: .24em; text-transform: uppercase; color: #a8883f; margin-bottom: 12px; }
.ph-title {
  font-family: var(--font-display);
  font-size: 36px; line-height: 1.45; font-weight: 700;
  color: #24345a; letter-spacing: .01em; margin-bottom: 12px;
}
.ph-title span { color: #6b7891; font-weight: 400; font-size: .86em; }
.ph-poem { font-size: 14px; color: #5d6783; margin: 0 0 20px; line-height: 1.9; font-family: var(--font-poem); }
.ph-cta {
  display: inline-flex; align-items: center; gap: 8px;
  color: var(--primary-700); font-size: 14px; font-weight: 600;
  padding: 8px 22px; border: 1px solid var(--primary-200); border-radius: 999px;
  transition: background .18s, border-color .18s;
}
.ph-cta:hover { background: var(--primary-50); border-color: var(--primary); }
.ph-cta .arr { display: inline-block; transition: transform .18s; }
.ph-cta:hover .arr { transform: translateX(4px); }
.ph-actions { display: flex; gap: 12px; flex-wrap: wrap; }
.ph-cta.ghost { background: transparent; color: var(--text-muted); border-color: var(--border-strong); font-weight: 500; }
.ph-cta.ghost:hover { color: var(--primary); border-color: var(--primary); background: var(--primary-50); }
@media (max-width: 640px) {
  .ph-inner { padding: 24px 20px; }
  .ph-title { font-size: 24px; }
}

/* —— 技术主题入口 —— */
.tech-top { display: grid; grid-template-columns: repeat(5, 1fr); gap: 12px; margin-bottom: 24px; }
.tech-tile {
  display: flex; align-items: center; gap: 10px;
  padding: 12px 14px; background: var(--surface);
  border: 1px solid var(--border); border-radius: 12px; box-shadow: var(--shadow-sm);
  transition: transform .16s ease, border-color .16s ease, box-shadow .16s ease;
}
.tech-tile:hover { transform: translateY(-2px); border-color: var(--primary-200); box-shadow: var(--shadow); }
.t-badge {
  width: 30px; height: 30px; border-radius: 9px; flex-shrink: 0;
  background: var(--primary-50); color: var(--primary-700);
  display: grid; place-items: center; font-weight: 700;
}
.tech-tile:nth-child(even) .t-badge { background: #f3e6c8; color: #8a6d2f; }
.t-name { flex: 1; min-width: 0; font-size: 14px; font-weight: 600; color: #24345a; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.t-num { font-size: 12px; color: var(--text-faint); white-space: nowrap; }
@media (max-width: 900px) { .tech-top { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 480px) { .tech-top { grid-template-columns: 1fr; } }

/* —— 统计条 —— */
.stats {
  display: flex; align-items: center; justify-content: center; gap: 30px;
  background: var(--surface); border: 1px solid var(--border); border-radius: 14px;
  padding: 16px 20px; margin-bottom: 26px; box-shadow: var(--shadow-sm);
}
.stat { display: flex; align-items: baseline; gap: 8px; }
.stat .num { font-family: var(--font-display); font-size: 24px; font-weight: 700; color: var(--primary-700); }
.stat .lbl { font-size: 12.5px; color: var(--text-muted); }
.stats .divider { width: 1px; height: 22px; background: var(--border); }

.home-cols { display: grid; grid-template-columns: minmax(0, 1fr) 300px; gap: 26px; align-items: start; min-width: 0; }
.main-col { min-width: 0; max-width: 100%; }
.side { min-width: 0; max-width: 100%; }
.section-title {
  display: flex; align-items: baseline; justify-content: space-between;
  margin: 4px 0 16px; padding-bottom: 10px; border-bottom: 1px solid var(--border);
}
.section-title .tt {
  position: relative; font-family: var(--font-display);
  font-size: 20px; font-weight: 700; color: #2c261d;
}
.section-title .tt::after {
  content: ""; position: absolute; left: 0; bottom: -11px; width: 34px; height: 3px;
  border-radius: 3px; background: var(--accent);
}
.section-title .more { font-size: 13px; color: var(--text-muted); }
.section-title .more:hover { color: var(--primary); }
.title-acts { display: flex; align-items: center; gap: 12px; }
.section-title .more.jump { border: none; background: none; font-family: inherit; padding: 0; cursor: pointer; }

/* —— 侧栏 —— */
.side-card {
  background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius);
  padding: 18px 20px; margin-bottom: 16px; box-shadow: var(--shadow-sm);
}
.side-title {
  font-family: var(--font-display);
  font-size: 16px; font-weight: 700; margin-bottom: 12px; color: #2c261d;
  display: flex; align-items: center; gap: 8px;
}
.side-title::before { content: "❀ "; color: var(--accent); font-size: 15px; }
.hot-list { list-style: none; counter-reset: hot; }
.hot-list li {
  display: flex; gap: 10px; padding: 9px 0; border-bottom: 1px dashed var(--border);
  font-size: 13.5px; align-items: center;
}
.hot-list li:last-child { border-bottom: none; }
.hot-list li::before {
  counter-increment: hot; content: counter(hot, decimal-leading-zero);
  flex-shrink: 0; font-family: var(--font-display); font-size: 14px; color: var(--text-faint);
}
.hot-list li:nth-child(-n+3)::before { color: var(--accent); font-weight: 700; }
.hot-list a { flex: 1; min-width: 0; color: var(--text);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.hot-list a:hover { color: var(--primary-600); }
.hot-list .views { color: var(--text-faint); font-size: 12px; white-space: nowrap; }
.tag-cloud { display: flex; flex-wrap: wrap; gap: 8px; }
.tag-chip em { font-style: normal; color: var(--text-faint); font-size: 12px; margin-left: 3px; }
.tag-chip.flower::before { content: "❀"; font-size: 11px; color: var(--accent); margin-right: 3px; }
.cat-list { display: flex; flex-direction: column; gap: 3px; }
.cat-item {
  display: flex; align-items: center; gap: 10px; padding: 6px 8px;
  border-radius: 9px; color: var(--text); transition: background .15s;
}
.cat-item:hover { background: var(--primary-50); color: var(--primary-700); }
.ci-badge {
  width: 24px; height: 24px; border-radius: 50%; flex-shrink: 0;
  background: var(--primary-100); color: var(--primary-700);
  display: grid; place-items: center; font-size: 12.5px; font-weight: 600;
}
.cat-item:nth-child(even) .ci-badge { background: #f3e6c8; color: #8a6d2f; }
.ci-name { flex: 1; min-width: 0; font-size: 13.5px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ci-num { font-size: 11.5px; color: var(--text-faint); background: var(--bg); padding: 1px 7px; border-radius: 999px; }
.ci-arr { color: var(--accent); font-size: 12px; opacity: 0; transform: translateX(-3px); transition: .15s; }
.cat-item:hover .ci-arr { opacity: 1; transform: none; }
.author-body { display: flex; gap: 12px; align-items: center; }
.author-avatar { width: 52px; height: 52px; border-radius: 50%; object-fit: cover; flex-shrink: 0; box-shadow: var(--shadow-sm); }
.author-avatar-text { display: grid; place-items: center; background: linear-gradient(135deg, var(--primary-100), #d8cdb2); color: var(--primary-700); font-weight: 700; font-size: 18px; }
.author-meta { min-width: 0; }
.author-name { font-family: var(--font-display); font-weight: 700; font-size: 15px; margin-bottom: 4px; color: #2c261d; }
.author-bio { font-size: 12.5px; color: var(--text-muted); line-height: 1.6; margin: 0; }
.author-more { display: inline-block; margin-top: 6px; font-size: 12.5px; color: var(--primary); }
@media (max-width: 900px) {
  .home-cols { grid-template-columns: 1fr; }
}
</style>
