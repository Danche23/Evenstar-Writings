<template>
  <div class="page-container page-block">
    <!-- 卡1：博主 -->
    <section class="block author-block">
      <img v-if="author.avatar" :src="author.avatar" class="big-avatar" alt="" />
      <div v-else class="big-avatar big-avatar-text">{{ (author.nickname || '暮').slice(0, 1) }}</div>
      <div class="a-meta">
        <div class="a-name">{{ author.nickname || '暮星' }}</div>
        <p class="a-bio">{{ author.bio || '写代码的技术人，业余记录与分享。' }}</p>
      </div>
      <div class="a-social"><SocialLinks /></div>
    </section>

    <!-- 卡2：介绍（内容由后台「关于设置」管理） -->
    <section class="block">
      <h2 class="block-title">关于这个站点</h2>
      <div class="intro-md" v-html="renderIntro()"></div>
      <p class="mini-stats">文章 {{ statsItems.articles }} · 评论 {{ statsItems.comments }} · 分类 {{ statsItems.categories }} · 标签 {{ statsItems.tags }}</p>
      <div v-if="doc.stack && doc.stack.length" class="chips">
        <span v-for="t in doc.stack" :key="t" class="chip">{{ t }}</span>
      </div>
    </section>

    <!-- 卡3：交流 -->
    <section class="block">
      <h2 class="block-title">如何与我交流</h2>
      <ul class="contact-list">
        <li><router-link to="/guestbook">留言板</router-link><span>——想说点什么，都可以留在这里。</span></li>
        <li><router-link :to="{ name: 'articles' }">文章评论区</router-link><span>——对某篇内容有疑问或指正，欢迎在文末讨论。</span></li>
      </ul>
      <p v-if="doc.footnote" class="footnote poem">{{ doc.footnote }}</p>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAuthor } from '@/api/user'
import { getPublicStats } from '@/api/stats'
import { listCategories } from '@/api/category'
import { listTags } from '@/api/tag'
import { getAbout } from '@/api/about'
import { renderMarkdown } from '@/utils/markdown'
import SocialLinks from '@/components/SocialLinks.vue'

const author = ref({ nickname: '', avatar: '', bio: '' })
const statsItems = ref({ articles: '—', comments: '—', categories: '—', tags: '—' })

// 后台未配置时后端会返回默认文案；此处本地兜底保证离线也完整
const fallbackDoc = {
  intro: '「暮星随笔」是一个以技术为主、内容不限的个人博客。整站由我独立实现：后端 Go + Gin + GORM，MySQL + Redis，JWT 鉴权，接入阿里云 OSS 与验证码；前端 Vue3 + Element Plus，Markdown 写作与渲染。\n\n不追热点，只写实践过、想明白的东西。偶尔也写点生活，放在角落，不影响主线。',
  stack: ['Go', 'Gin', 'GORM', 'MySQL', 'Redis', 'JWT', 'Vue3', 'Element Plus', 'Markdown'],
  footnote: '站名 Evenstar（暮星）取自萨福笔下那颗把白昼散落之物带回家的星——愿每篇文章，都把你带回来。'
}
const doc = ref({ ...fallbackDoc })

function renderIntro() {
  return renderMarkdown(doc.value.intro || '')
}

onMounted(async () => {
  try {
    const [authorResp, stats, catsResp, tagsResp, aboutResp] = await Promise.all([
      getAuthor(), getPublicStats(), listCategories(), listTags(), getAbout()
    ])
    if (authorResp) author.value = authorResp
    statsItems.value = {
      articles: stats?.articles ?? '—',
      comments: stats?.comments ?? '—',
      categories: (catsResp || []).length || '—',
      tags: (tagsResp || []).length || '—'
    }
    if (aboutResp) doc.value = aboutResp
  } catch (e) { /* 拦截器提示 */ }
})
</script>

<style scoped>
/* 内容列整体收窄（780）+ 卡片内部留足呼吸感，避免“矮胖” */
.block {
  max-width: 780px; margin: 0 auto 18px;
  background: var(--surface); border: 1px solid var(--border);
  border-radius: 16px; padding: 30px 36px;
  box-shadow: var(--shadow-sm);
}
.block-title {
  font-size: 18px; font-weight: 700; color: #24345a; margin-bottom: 16px;
  display: flex; align-items: center; gap: 9px;
}
.block-title::before { content: ""; width: 4px; height: 16px; border-radius: 2px; background: var(--primary); }

/* 卡1：博主 */
.author-block { display: flex; align-items: center; gap: 20px; padding: 24px 36px; }
.big-avatar { width: 66px; height: 66px; border-radius: 50%; object-fit: cover; flex-shrink: 0; }
.big-avatar-text { display: grid; place-items: center; background: linear-gradient(135deg, var(--primary-100), #d8cdb2); color: var(--primary-700); font-weight: 700; font-size: 24px; }
.a-meta { flex: 1; min-width: 0; }
.a-name { font-family: var(--font-display); font-size: 20px; font-weight: 700; color: #24345a; }
.a-bio { font-size: 13.5px; color: var(--text-muted); line-height: 1.7; margin: 6px 0 0; }
.a-social { margin-left: auto; flex-shrink: 0; }
.a-social :deep(.social-links) { margin-top: 0; }

/* 卡2：介绍（Markdown 渲染） */
.intro-md { font-size: 14.5px; line-height: 2.05; color: var(--text); }
.intro-md :deep(p) { margin: 0 0 14px; }
.intro-md :deep(p:last-child) { margin-bottom: 0; }
.intro-md :deep(a) { color: var(--primary-600); }
.intro-md :deep(a:hover) { text-decoration: underline; }
.intro-md :deep(code) {
  background: var(--surface-2); border-radius: 4px; padding: 1px 6px;
  font-family: var(--mono); font-size: 13px;
}
.mini-stats { font-size: 12.5px; color: var(--text-faint); margin: 18px 0 0 !important; }
.chips { display: flex; flex-wrap: wrap; gap: 7px; margin-top: 14px; }
.chip { font-size: 12px; color: #55617a; background: var(--surface-2); border-radius: 999px; padding: 3px 12px; }

/* 卡3：交流 */
.contact-list { list-style: none; margin: 0; padding: 0; }
.contact-list li { font-size: 14.5px; line-height: 2.15; color: var(--text); }
.contact-list a { color: var(--primary-600); font-weight: 600; }
.contact-list a:hover { text-decoration: underline; }
.contact-list span { color: var(--text-muted); }
.footnote { font-size: 12.5px; color: var(--text-faint); line-height: 1.9; margin: 18px 0 0; border-top: 1px dashed var(--border); padding-top: 14px; }

@media (max-width: 640px) {
  .block { padding: 22px 20px; }
  .author-block { padding: 20px; flex-direction: column; align-items: flex-start; gap: 14px; }
  .a-meta, .a-social { width: 100%; }
  .a-social { margin-left: 0; }
}
</style>
