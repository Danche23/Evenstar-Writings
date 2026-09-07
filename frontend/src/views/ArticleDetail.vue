<template>
  <div class="page-container-narrow page-block" v-if="article">
    <!-- 阅读进度条 -->
    <div class="read-progress" :style="{ width: progress + '%' }" aria-hidden="true"></div>

    <!-- 返回文章列表 -->
    <router-link to="/articles" class="back-link">← 返回文章列表</router-link>

    <!-- 文章头部 -->
    <header class="article-header">
      <div class="cat-row">
        <span v-for="c in article.categories" :key="c.id" class="badge badge-primary">{{ c.name }}</span>
        <span v-for="t in article.tags" :key="t.id" class="badge badge-gray"># {{ t.name }}</span>
      </div>
      <h1>{{ article.title }}</h1>
      <div class="article-meta">
        <UserAvatar :user="article.author || { nickname: '暮星' }" :size="34" />
        <span>{{ article.author?.nickname || '暮星' }}</span>
        <span>发布于 {{ formatDate(article.published_at || article.created_at) }}</span>
        <span>👁 {{ article.views ?? 0 }} 次浏览</span>
      </div>
    </header>

    <!-- 封面 -->
    <div v-if="article.cover" class="article-cover">
      <img :src="article.cover" alt="cover" />
    </div>

    <!-- Markdown 正文 -->
    <MarkdownView :content="article.content" class="markdown-body-pad" />

    <!-- 心情题字：正文后的一行小尾注 -->
    <p v-if="epigraph" class="article-epi-foot poem">❝ {{ epigraph }} ❞</p>

    <!-- 目录（宽屏右侧悬浮） -->
    <nav v-if="toc.length > 1" class="toc" aria-label="目录">
      <div class="toc-title">目录</div>
      <a
        v-for="t in toc" :key="t.id"
        class="toc-item" :class="{ active: t.id === activeId, lv3: t.level === 3 }"
        @click.prevent="jumpToc(t.id)"
      >{{ t.text }}</a>
    </nav>

    <!-- 评论区 -->
    <section ref="commentsSectionRef" class="comments mt-24">
      <div class="comments-head">
        <h2>评论 <span class="count">{{ totalComments }}</span></h2>
      </div>

      <!-- ===== 常态发表卡（登录态：署名 + 输入；游客：登录引导） ===== -->
      <div ref="commentFormRef" class="card composer" :class="{ guest: !store.isLoggedIn }">
        <template v-if="store.isLoggedIn">
          <div class="composer-head">
            <UserAvatar :user="store.user" :size="36" />
            <div class="who">
              <b>{{ store.displayName }}</b>
              <span>写下你的想法，与其他读者交流</span>
            </div>
          </div>
          <el-input
            ref="draftInput"
            v-model="draft"
            type="textarea"
            :rows="3"
            maxlength="400"
            show-word-limit
            placeholder="写下你的评论，支持 Markdown…"
          />
          <div class="composer-foot">
            <span class="hint">支持 Markdown · 友善发言 · 最多 400 字</span>
            <button class="send-btn" :disabled="submitting || !draft.trim()" @click="submitComment">
              {{ submitting ? '发送中…' : '发表评论' }}
            </button>
          </div>
        </template>
        <template v-else>
          <div class="guest-in">
            <span class="guest-mark">💬</span>
            <div class="guest-copy">
              <b>登录后参与讨论</b>
              <p>注册或登录后即可对文章发表评论与回复</p>
            </div>
            <button class="send-btn send-btn-ghost" @click="goLogin">去登录</button>
          </div>
        </template>
      </div>

      <!-- 评论列表（两级树，后端已组装） -->
      <div v-if="commentsLoading" class="loading-box">评论加载中…</div>
      <template v-else>
        <div v-if="!comments.length" class="empty">
          <div class="big">💬</div>
          还没有评论，来抢沙发吧～
        </div>

        <div v-for="c in comments" :key="c.id" class="c-block" :class="{ top: c.is_top === 1 }">
          <!-- 一级评论：主题卡 -->
          <div class="c-main">
            <UserAvatar :user="c.user" :size="40" />
            <div class="c-body">
              <div class="c-bar">
                <b class="nick">{{ c.user ? c.user.nickname : '已注销用户' }}</b>
                <span v-if="c.is_top === 1" class="pin-badge">置顶</span>
                <span class="when">{{ formatDateTime(c.created_at) }}</span>
              </div>
              <div class="comment-md" v-html="renderMd(c.content)"></div>
              <div class="c-ops">
                <button v-if="!replyTo" class="op-btn" @click="startReply(c, null)">回复</button>
                <button v-if="canDelete(c)" class="op-btn danger" @click="removeComment(c.id)">删除</button>
              </div>
            </div>
          </div>

          <!-- 二级回复：挂线线程 -->
          <div v-if="c.replies && c.replies.length" class="c-thread">
            <div v-for="r in c.replies" :key="r.id" class="r-row">
              <UserAvatar :user="r.user" :size="30" />
              <div class="r-body">
                <div class="r-bar">
                  <b>{{ r.user ? r.user.nickname : '已注销用户' }}</b>
                  <span v-if="replyName(c, r)" class="r-to">{{ replyName(c, r) }}</span>
                  <span class="when">{{ formatDateTime(r.created_at) }}</span>
                </div>
                <div class="comment-md" v-html="renderMd(r.content)"></div>
                <div class="c-ops">
                  <button v-if="!replyTo" class="op-btn" @click="startReply(c, r)">回复</button>
                  <button v-if="canDelete(r)" class="op-btn danger" @click="removeComment(r.id)">删除</button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 一级评论分页 -->
        <div v-if="commentTotal > commentQuery.page_size" class="pager">
          <el-pagination
            layout="prev, pager, next"
            :total="commentTotal"
            :page-size="commentQuery.page_size"
            :current-page="commentQuery.page"
            @current-change="(p) => { commentQuery.page = p; loadComments() }"
          />
        </div>
      </template>
    </section>

    <!-- ===== 关联阅读 ===== -->
    <section v-if="relatedList.length" class="related">
      <div class="related-head"><span>{{ relatedLabel }}</span></div>
      <div class="related-grid">
        <router-link v-for="a in relatedList" :key="a.id" :to="`/articles/${a.id}`" class="rel-card">
          <span class="rel-title">{{ a.title }}</span>
          <span class="rel-meta">{{ formatDate(a.published_at || a.created_at) }} · {{ a.views ?? 0 }} 阅读</span>
        </router-link>
      </div>
    </section>

    <!-- ===== 回复浮层（自绘，teleport 到 body）：引用原话 + 输入 ===== -->
    <Teleport to="body">
      <transition name="rmodal">
        <div v-if="replyDialogVisible" class="r-overlay" @click.self="closeReply">
          <div class="r-panel" role="dialog" aria-modal="true" aria-label="回复评论">
            <div class="r-head">
              <div class="r-title">
                <span class="r-badge">回复</span>
                <span class="r-who"><b>{{ replyTo?.name || '评论' }}</b></span>
              </div>
              <button class="r-close" aria-label="关闭" @click="closeReply">
                <el-icon><Close /></el-icon>
              </button>
            </div>
            <div v-if="replyTo?.quote" class="r-quote">{{ replyTo.quote }}</div>
            <el-input
              ref="replyInput"
              v-model="replyDraft"
              type="textarea"
              :rows="3"
              maxlength="400"
              show-word-limit
              placeholder="写下你的回复…（Ctrl+Enter 快捷发送）"
              @keydown.ctrl.enter.prevent="submitReply"
            />
            <div class="r-foot">
              <span class="hint">支持 Markdown · 回复将公开展示在该条下方</span>
              <div class="r-btns">
                <button class="text-btn" @click="closeReply">取消</button>
                <button class="send-btn" :disabled="submitting || !replyDraft.trim()" @click="submitReply">
                  {{ submitting ? '发送中…' : '提交回复' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- ===== 滚动后快捷入口：右下角写评论按钮 ===== -->
    <transition name="fab">
      <button v-show="quickVisible" class="fab" aria-label="写评论" title="写评论" @click="jumpToComment">
        <el-icon :size="21"><ChatDotRound /></el-icon>
      </button>
    </transition>

    <!-- 回到顶部 -->
    <transition name="fab">
      <button v-show="showTop" class="back-top" aria-label="回到顶部" @click="toTop">↑</button>
    </transition>
  </div>
</template>

<script setup>
import { ref, reactive, watch, nextTick, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ChatDotRound, Close } from '@element-plus/icons-vue'
import MarkdownView from '@/components/MarkdownView.vue'
import { getArticle, recordView, listArticles } from '@/api/article'
import { listComments, createComment, deleteComment } from '@/api/comment'
import { useUserStore } from '@/stores/user'
import { formatDate, formatDateTime } from '@/utils/format'
import { renderMarkdown, extractToc } from '@/utils/markdown'
import { epigraphFor } from '@/utils/quotes'
import UserAvatar from '@/components/UserAvatar.vue'

const route = useRoute()
const router = useRouter()
const store = useUserStore()
const articleId = Number(route.params.id)

// 评论正文渲染（markdown-it → DOMPurify 消毒，保留换行/空格）
const renderMd = (src = '') => renderMarkdown(src)

const article = ref(null)
const draft = ref('')
const submitting = ref(false)

const epigraph = computed(() => (article.value ? epigraphFor(article.value.title) : ''))

const comments = ref([])
const commentTotal = ref(0)    // 未删除一级评论数（分页依据）
const commentCount = ref(0)    // 未删除一级+二级总数（头部展示）
const commentsLoading = ref(false)
const commentQuery = reactive({ page: 1, page_size: 10 })

// 回复浮层：{ name 被回复人, quote 原话摘要, parentCommentId, replyToId }
const replyTo = ref(null)
const replyDialogVisible = ref(false)
const replyDraft = ref('')
const replyInput = ref(null)
let viewTimer = null

const totalComments = computed(() => commentCount.value)

// —— 关联阅读：同分类/同标签的其它文章 ——
const relatedList = ref([])
const relatedLabel = ref('继续读')

async function loadRelated() {
  const a = article.value
  if (!a) return
  const cats = a.categories || []
  const tgs = a.tags || []
  const cat = cats[0]
  const tag = tgs[0]
  const params = {}
  if (cat) { params.category_id = cat.id; relatedLabel.value = `更多 · ${cat.name}` }
  else if (tag) { params.tag_id = tag.id; relatedLabel.value = `更多 · ${tag.name}` }
  else return
  try {
    const d = await listArticles({ ...params, page: 1, page_size: 4 })
    relatedList.value = (d.list || []).filter((x) => x.id !== articleId).slice(0, 3)
  } catch (e) { relatedList.value = [] }
}

// —— 浮动「写评论」按钮：主评论框滚出视野后在评论区浮现 ——
const commentFormRef = ref(null)
const commentsSectionRef = ref(null)
const draftInput = ref(null)
const quickVisible = ref(false)
const progress = ref(0)
const showTop = ref(false)
const activeId = ref('')

const toc = computed(() => (article.value ? extractToc(article.value.content) : []))

function onScroll() {
  // 阅读进度
  const doc = document.documentElement
  const max = doc.scrollHeight - window.innerHeight
  progress.value = max > 0 ? Math.min(100, Math.round((window.scrollY / max) * 1000) / 10) : 0

  showTop.value = window.scrollY > 600

  // TOC 当前章节高亮
  if (toc.value.length) {
    let cur = ''
    for (const t of toc.value) {
      const el = document.getElementById(t.id)
      if (el && el.getBoundingClientRect().top <= 130) cur = t.id
    }
    activeId.value = cur
  }

  const form = commentFormRef.value
  const sec = commentsSectionRef.value
  if (!form || !sec) return
  const fr = form.getBoundingClientRect()
  const sr = sec.getBoundingClientRect()
  // 主评论框完全滚出视口上方、且评论区尚未完全滚出底部时显示
  quickVisible.value = fr.bottom <= 0 && sr.bottom > 80
}

function jumpToc(id) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function toTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function goLogin() {
  router.push({ name: 'login', query: { redirect: route.fullPath } })
}

function jumpToComment() {
  if (!store.isLoggedIn) {
    ElMessage.warning('请先登录后再评论')
    goLogin()
    return
  }
  commentFormRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  // 等平滑滚动到位后再聚焦输入框
  setTimeout(() => draftInput.value?.focus(), 500)
}

function replyName(top, r) {
  // 「A → B」：B 为被回复对象（reply_to_id 指向的用户昵称）
  if (!r.reply_to_id) return ''
  const target = [...(top.replies || [])].find((x) => x.id === r.reply_to_id)
  if (target?.user) return target.user.nickname
  // 指向一级或已注销
  if (top.id === r.reply_to_id) return top.user?.nickname || '已注销用户'
  return '已注销用户'
}

function canDelete(c) {
  if (c.deleted) return false
  // 管理员可删任何；用户可删自己的
  return store.isAdmin || (store.user && c.user_id === store.user.id)
}

function startReply(top, target) {
  if (!store.isLoggedIn) {
    ElMessage.warning('请先登录后再回复')
    goLogin()
    return
  }
  // 回复二级时 parent 恒为一级（target 为二级时取其一上级 id）；reply_to 为实际对象
  const targetComment = target || top
  replyTo.value = {
    parentCommentId: target && target.parent_id ? target.parent_id : top.id,
    replyToId: target ? target.id : null,
    name: targetComment.user?.nickname || (targetComment.user ? '' : '已注销用户'),
    quote: (targetComment.content || '').replace(/\s+/g, ' ').trim().slice(0, 80)
  }
  replyDraft.value = ''
  replyDialogVisible.value = true
}

function closeReply() {
  replyDialogVisible.value = false
}

// 浮层打开后聚焦输入框；关闭后复位回复目标（恢复各条「回复」入口）
watch(replyDialogVisible, (v) => {
  if (v) {
    nextTick(() => replyInput.value?.focus())
  } else {
    replyTo.value = null
  }
})

function onKeydown(e) {
  if (e.key === 'Escape' && replyDialogVisible.value) closeReply()
}

// 浮层内提交回复
async function submitReply() {
  if (!replyDraft.value.trim()) return
  const payload = { content: replyDraft.value.trim() }
  if (replyTo.value?.parentCommentId) payload.parent_id = replyTo.value.parentCommentId
  if (replyTo.value?.replyToId) payload.reply_to_id = replyTo.value.replyToId
  submitting.value = true
  try {
    await createComment(articleId, payload)
    ElMessage.success('回复成功')
    replyDialogVisible.value = false
    replyDraft.value = ''
    commentQuery.page = 1
    await loadComments()
  } catch (e) {
    if (e.code === 429) ElMessage.warning('评论太频繁，请稍后再试')
  } finally {
    submitting.value = false
  }
}

async function submitComment() {
  if (!store.isLoggedIn) {
    ElMessage.warning('请先登录后再评论')
    goLogin()
    return
  }
  if (!draft.value.trim()) return
  submitting.value = true
  try {
    await createComment(articleId, { content: draft.value.trim() })
    ElMessage.success('评论成功')
    draft.value = ''
    commentQuery.page = 1
    await loadComments()
  } catch (e) {
    if (e.code === 429) ElMessage.warning('评论太频繁，请稍后再试')
  } finally {
    submitting.value = false
  }
}

async function removeComment(id) {
  try {
    await deleteComment(id)
    ElMessage.success('评论已删除')
    await loadComments()
  } catch (e) { /* 拦截器提示 */ }
}

async function loadComments() {
  commentsLoading.value = true
  try {
    const data = await listComments(articleId, { page: commentQuery.page, page_size: commentQuery.page_size })
    comments.value = data.list || []
    commentTotal.value = Number(data.total || 0)
    commentCount.value = Number(data.count ?? data.total ?? 0)
  } catch (e) { /* 拦截器提示 */ } finally {
    commentsLoading.value = false
  }
}

onMounted(async () => {
  window.addEventListener('scroll', onScroll, { passive: true })
  window.addEventListener('keydown', onKeydown)
  try {
    article.value = await getArticle(articleId)
    await loadRelated()
  } catch (e) {
    ElMessage.error(e.message || '文章不存在或已删除')
  }

  loadComments()

  // 详情停留 5 秒后上报一次浏览（后端做 15 分钟防刷 + IP 限流）
  viewTimer = setTimeout(async () => {
    try {
      const data = await recordView(articleId)
      if (data && article.value) article.value.views = data.views
    } catch (e) { /* 浏览失败静默，不影响阅读 */ }
  }, 5000)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll)
  window.removeEventListener('keydown', onKeydown)
  if (viewTimer) clearTimeout(viewTimer)
})
</script>

<style scoped>
.article-header { padding: 30px 0 8px; }
.read-progress {
  position: fixed; top: 0; left: 0; height: 3px; z-index: 130;
  background: linear-gradient(90deg, var(--primary), var(--accent));
  border-radius: 0 3px 3px 0;
  box-shadow: 0 0 10px rgba(92, 122, 79, 0.35);
}
.cat-row { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }
.article-header h1 { font-size: 30px; line-height: 1.35; font-weight: 700; margin-bottom: 14px; }
.article-meta { display: flex; gap: 18px; align-items: center; font-size: 13px; color: var(--text-muted); flex-wrap: wrap; }
.article-cover { border-radius: 12px; margin: 18px 0 6px; overflow: hidden; }
.article-epi-foot {
  margin: 6px 0 4px; text-align: center;
  font-size: 12.5px; line-height: 1.8; color: var(--text-faint);
}
.article-cover img { width: 100%; max-height: 380px; object-fit: cover; }
.markdown-body-pad { padding: 18px 0 10px; }

/* ============ 评论区：发表卡 / 列表 / 回复浮层 / FAB 统一视觉 ============ */
.comments { margin-top: 10px; }
.comments-head { margin-bottom: 16px; }
.comments-head h2 { font-size: 19px; font-weight: 700; display: flex; align-items: center; gap: 8px; }
.comments-head .count {
  min-width: 22px; height: 22px; padding: 0 7px; border-radius: 999px;
  background: var(--primary-50); color: var(--primary-600);
  font-size: 12.5px; font-weight: 600; display: inline-flex; align-items: center; justify-content: center;
}

/* —— 发送主按钮（发表 / 回复共用）—— */
.send-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  padding: 8px 20px; border: none; border-radius: 999px;
  background: var(--primary); color: #fff; font-size: 13.5px; font-weight: 500;
  box-shadow: 0 4px 12px rgba(92, 122, 79, 0.22);
  transition: background .18s ease, transform .18s ease, box-shadow .18s ease;
}
.send-btn:hover:not(:disabled) { background: var(--primary-600); transform: translateY(-1px); box-shadow: 0 6px 18px rgba(92, 122, 79, 0.3); }
.send-btn:disabled { opacity: .45; cursor: not-allowed; box-shadow: none; }
.send-btn-ghost {
  background: #fff; color: var(--primary); border: 1px solid var(--border-strong); box-shadow: none;
}
.send-btn-ghost:hover:not(:disabled) { background: var(--primary-50); color: var(--primary-600); transform: none; box-shadow: none; }
.hint { font-size: 12.5px; color: var(--text-faint); }

/* —— 常态发表卡 —— */
.composer {
  padding: 18px 20px 16px; border-radius: 14px; margin-bottom: 20px;
  scroll-margin-top: 84px;
  box-shadow: var(--shadow-sm);
}
.composer-head { display: flex; align-items: center; gap: 10px; margin-bottom: 12px; }
.composer-head .who { display: flex; flex-direction: column; line-height: 1.35; }
.composer-head .who b { font-size: 14px; }
.composer-head .who span { font-size: 12px; color: var(--text-faint); }
.composer :deep(.el-textarea__inner) {
  border-radius: 10px; padding: 10px 12px; font-size: 14px; line-height: 1.7;
  background: #f8f3e8; border-color: var(--border); transition: border-color .15s, box-shadow .15s, background .15s;
}
.composer :deep(.el-textarea__inner:focus) {
  background: #fff; border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-100);
}
.composer :deep(.el-input__count) { font-size: 12px; color: var(--text-faint); }
.composer-foot { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 12px; }

/* 游客引导态 */
.composer.guest { background: linear-gradient(180deg, #fffdf7, #f6f0e1); }
.guest-in { display: flex; align-items: center; gap: 14px; }
.guest-mark {
  width: 42px; height: 42px; border-radius: 12px; flex-shrink: 0;
  background: var(--primary-50); color: var(--primary-600);
  display: grid; place-items: center; font-size: 19px;
}
.guest-copy { flex: 1; min-width: 0; line-height: 1.5; }
.guest-copy b { font-size: 14.5px; display: block; }
.guest-copy p { font-size: 12.5px; color: var(--text-muted); margin: 2px 0 0; }

/* —— 评论列表：对话式卡片（一级主题卡 / 二级挂线回复线程） —— */
.c-block {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: 14px; padding: 16px 18px; margin-bottom: 14px;
  box-shadow: var(--shadow-sm); transition: box-shadow .15s ease;
}
.c-block:hover { box-shadow: var(--shadow); }
.c-block.top {
  border-color: var(--primary-200);
  background: linear-gradient(180deg, var(--primary-50) 0%, var(--surface) 58%);
}
.c-main { display: flex; gap: 12px; align-items: flex-start; }
.c-body { flex: 1; min-width: 0; }
.c-bar { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.c-bar .nick { font-size: 14.5px; font-weight: 700; color: #2c261d; }
.pin-badge {
  font-size: 11px; padding: 1px 9px; border-radius: 999px;
  background: var(--primary); color: #fff; letter-spacing: .06em;
}
.c-bar .when { margin-left: auto; font-size: 12px; color: var(--text-faint); white-space: nowrap; }

/* 评论正文（Markdown）：紧凑行内排版（两级共用） */
.comment-md { font-size: 14px; line-height: 1.75; color: var(--text); word-break: break-word; margin: 6px 0 2px; }
.comment-md > :first-child { margin-top: 0 !important; }
.comment-md > :last-child { margin-bottom: 0 !important; }
.comment-md p { margin: 4px 0; }
.comment-md a { color: var(--primary); text-decoration: none; }
.comment-md a:hover { text-decoration: underline; }
.comment-md code {
  background: #efe9db; border-radius: 4px; padding: 1px 6px;
  font-family: var(--mono); font-size: 12.5px; color: #9a4b26;
}
.comment-md pre { background: #23281e; border-radius: 8px; padding: 10px 12px; overflow-x: auto; margin: 8px 0; }
.comment-md pre code { background: none; color: #e6e1d2; padding: 0; }
.comment-md ul, .comment-md ol { margin: 4px 0; padding-left: 22px; }
.comment-md blockquote { border-left: 3px solid var(--primary); padding: 4px 12px; margin: 6px 0; background: var(--primary-50); border-radius: 0 8px 8px 0; color: #6b6253; }
.comment-md img { max-width: 100%; border-radius: 8px; margin: 6px 0; }

.c-ops { display: flex; align-items: center; gap: 2px; margin-top: 2px; }
.op-btn {
  border: none; background: transparent; padding: 3px 10px; border-radius: 8px;
  font-size: 12.5px; color: var(--text-muted); cursor: pointer;
  transition: color .15s, background .15s;
}
.op-btn:hover { color: var(--primary); background: var(--primary-50); }
.op-btn.danger:hover { color: var(--danger); background: var(--danger-50); }

/* —— 二级回复线程：浅纸面挂线列表 —— */
.c-thread {
  margin: 12px -8px -4px; padding: 12px 14px 2px;
  background: #f1ebdc; border-radius: 12px;
}
.r-row { display: flex; gap: 10px; padding: 8px 2px 10px; }
.r-row + .r-row { border-top: 1px dashed rgba(151, 136, 110, .3); }
.r-body { flex: 1; min-width: 0; }
.r-bar { display: flex; align-items: baseline; gap: 8px; flex-wrap: wrap; }
.r-bar b { font-size: 13.5px; color: #3a3329; }
.r-to { font-size: 12px; color: var(--text-faint); }
.r-to::before { content: "↳ 回复 "; opacity: .85; }
.r-bar .when { margin-left: auto; font-size: 11.5px; color: var(--text-faint); white-space: nowrap; }
.r-body .c-ops { margin-top: 0; }
.r-body .op-btn { font-size: 12px; padding: 2px 8px; }

.pager { display: flex; justify-content: center; margin: 16px 0; }

/* —— 关联阅读 —— */
.related { margin: 30px 0 4px; }
.related-head {
  font-family: var(--font-display); font-size: 16px; font-weight: 700; color: #24345a;
  display: flex; align-items: center; gap: 8px; margin-bottom: 12px;
}
.related-head::before { content: "❀"; color: var(--accent); font-size: 14px; }
.related-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.rel-card {
  display: flex; flex-direction: column; gap: 6px;
  padding: 14px 15px; border: 1px solid var(--border); border-radius: 12px;
  background: var(--surface); transition: box-shadow .16s, border-color .16s;
}
.rel-card:hover { border-color: var(--primary-200); box-shadow: var(--shadow); }
.rel-title { font-size: 14px; font-weight: 600; line-height: 1.5; color: #24345a;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.rel-card:hover .rel-title { color: var(--primary-700); }
.rel-meta { margin-top: auto; font-size: 12px; color: var(--text-faint); }
@media (max-width: 640px) { .related-grid { grid-template-columns: 1fr; } }

/* —— 回复浮层 —— */
.r-overlay {
  position: fixed; inset: 0; z-index: 200;
  background: rgba(15, 23, 42, 0.42); backdrop-filter: blur(3px);
  display: grid; place-items: center; padding: 20px;
}
.r-panel {
  width: min(500px, 100%); background: #fff; border-radius: 18px;
  padding: 18px 20px 16px; box-shadow: 0 24px 64px rgba(15, 23, 42, 0.22);
  display: flex; flex-direction: column; gap: 12px;
}
.r-head { display: flex; align-items: center; justify-content: space-between; }
.r-title { display: flex; align-items: center; gap: 8px; min-width: 0; }
.r-badge {
  flex-shrink: 0; padding: 2px 10px; border-radius: 999px;
  background: var(--primary); color: #fff; font-size: 12px; font-weight: 500;
}
.r-who { font-size: 14.5px; color: var(--text); min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.r-who b { font-weight: 600; }
.r-close {
  width: 28px; height: 28px; border: none; border-radius: 8px; background: transparent;
  color: var(--text-faint); display: grid; place-items: center; cursor: pointer;
  transition: background .15s, color .15s;
}
.r-close:hover { background: #f3f4f6; color: var(--text); }
.r-quote {
  position: relative; padding: 9px 12px 9px 32px; border-radius: 10px;
  background: #f6f7f9; color: var(--text-muted); font-size: 13px; line-height: 1.6;
  overflow: hidden; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical;
}
.r-quote::before { content: "“"; position: absolute; left: 12px; top: 2px; font-size: 22px; color: var(--primary); font-weight: 700; }
.r-panel :deep(.el-textarea__inner) {
  border-radius: 10px; padding: 10px 12px; font-size: 14px; line-height: 1.7;
  background: #f8f3e8; border-color: var(--border); transition: border-color .15s, box-shadow .15s, background .15s;
}
.r-panel :deep(.el-textarea__inner:focus) {
  background: #fff; border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-100);
}
.r-panel :deep(.el-input__count) { font-size: 12px; color: var(--text-faint); }
.r-foot { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.r-btns { display: flex; align-items: center; gap: 10px; }
.text-btn { border: none; background: transparent; color: var(--text-muted); font-size: 13.5px; padding: 7px 10px; border-radius: 8px; cursor: pointer; transition: color .15s, background .15s; }
.text-btn:hover { color: var(--text); background: #f3f4f6; }

/* 浮层过渡：遮罩淡入 + 面板轻升起 */
.rmodal-enter-active, .rmodal-leave-active { transition: opacity .22s ease; }
.rmodal-enter-from, .rmodal-leave-to { opacity: 0; }
.rmodal-enter-active .r-panel, .rmodal-leave-active .r-panel { transition: transform .26s cubic-bezier(.16, 1, .3, 1), opacity .22s ease; }
.rmodal-enter-from .r-panel, .rmodal-leave-to .r-panel { transform: translateY(14px) scale(.97); opacity: 0; }

/* —— 右下角「写评论」FAB —— */
.fab {
  position: fixed; right: 26px; bottom: 26px; z-index: 90;
  width: 52px; height: 52px; border-radius: 50%; border: none;
  background: var(--primary); color: #fff; cursor: pointer;
  display: grid; place-items: center;
  box-shadow: 0 10px 26px rgba(70, 99, 156, 0.35);
  transition: transform .18s ease, box-shadow .18s ease, background .18s ease;
}
.fab:hover { background: var(--primary-600); transform: translateY(-2px); box-shadow: 0 14px 32px rgba(70, 99, 156, 0.45); }
/* 回到顶部（FAB 上方小圆钮） */
.back-top {
  position: fixed; right: 30px; bottom: 92px; z-index: 89;
  width: 38px; height: 38px; border-radius: 50%; cursor: pointer;
  background: var(--surface); color: var(--text-muted); border: 1px solid var(--border);
  display: grid; place-items: center; font-size: 17px; line-height: 1;
  box-shadow: var(--shadow-sm); transition: color .15s, border-color .15s, transform .15s;
}
.back-top:hover { color: var(--primary); border-color: var(--primary); transform: translateY(-2px); }
.fab-enter-active, .fab-leave-active, .back-top-enter-active, .back-top-leave-active { transition: opacity .22s ease, transform .22s ease; }
.fab-enter-from, .fab-leave-to, .back-top-enter-from, .back-top-leave-to { opacity: 0; transform: scale(.6); }

/* —— 文章目录（超宽屏右侧悬浮） —— */
.toc {
  position: fixed; left: calc(50% + 424px); top: 96px; z-index: 50;
  width: 184px; max-height: 62vh; overflow-y: auto;
  border-left: 1px solid var(--border); padding: 2px 0 2px 16px;
}
.toc-title { font-size: 12px; letter-spacing: .18em; color: var(--text-faint); margin-bottom: 8px; text-transform: uppercase; }
.toc-item {
  display: block; padding: 4px 8px; border-radius: 7px 0 0 7px; margin: 1px 0;
  font-size: 13px; color: var(--text-muted); line-height: 1.55; cursor: pointer;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  transition: color .15s, background .15s;
}
.toc-item:hover { color: var(--primary); background: var(--primary-50); }
.toc-item.lv3 { padding-left: 20px; font-size: 12.5px; }
.toc-item.active { color: var(--primary-700); font-weight: 600; background: var(--primary-50); box-shadow: inset 2px 0 0 var(--accent); }
@media (max-width: 1380px) { .toc { display: none; } }

@media (max-width: 640px) {
  .c-block { padding: 14px 14px; }
  .c-thread { margin: 10px -6px -2px; padding: 10px 10px 2px; }
  .composer-foot, .r-foot { flex-direction: column; align-items: stretch; }
  .composer-foot .send-btn, .r-foot .send-btn { width: 100%; }
  .r-btns { justify-content: flex-end; }
  .fab { right: 16px; bottom: 16px; width: 48px; height: 48px; }
}
</style>
