<template>
  <div class="page-container-narrow page-block" v-if="article">
    <!-- 文章头部 -->
    <header class="article-header">
      <div class="cat-row">
        <span v-for="c in article.categories" :key="c.id" class="badge badge-primary">{{ c.name }}</span>
        <span v-for="t in article.tags" :key="t.id" class="badge badge-gray"># {{ t.name }}</span>
      </div>
      <h1>{{ article.title }}</h1>
      <div class="article-meta">
        <span class="avatar">{{ (article.author?.nickname || '灯影').slice(0, 1) }}</span>
        <span>{{ article.author?.nickname || '灯影' }}</span>
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

    <!-- 评论区 -->
    <section class="comments mt-24">
      <div class="comments-head">
        <h2>评论 <span class="faint" style="font-size: 14px">({{ totalComments }})</span></h2>
      </div>

      <!-- 发表评论 -->
      <div class="card comment-form">
        <div class="form-label" style="font-weight: 600">发表评论</div>
        <el-input
          v-model="draft"
          type="textarea"
          :rows="3"
          maxlength="400"
          show-word-limit
          :placeholder="replyTo ? `回复「${replyTo.name}」：` : '写下你的评论…'"
        />
        <div class="actions">
          <div class="login-tip">
            <template v-if="replyTo">
              正在回复 <b>{{ replyTo.name }}</b>
              <a class="cancel-reply" @click="cancelReply">取消回复</a>
            </template>
            <template v-else>支持 Markdown · 文明评论 · 最多 400 字</template>
          </div>
          <button class="btn btn-primary btn-sm" :disabled="submitting || !draft.trim()" @click="submitComment">
            {{ replyTo ? '提交回复' : '发表评论' }}
          </button>
        </div>
      </div>

      <!-- 评论列表（两级树，后端已组装） -->
      <div v-if="commentsLoading" class="loading-box">评论加载中…</div>
      <template v-else>
        <div v-if="!comments.length" class="empty">
          <div class="big">💬</div>
          还没有评论，来抢沙发吧～
        </div>

        <div v-for="c in comments" :key="c.id" class="comment-wrap">
          <!-- 一级评论 -->
          <div class="comment-item" :class="{ top: c.is_top === 1 }">
            <span class="avatar" :class="{ gray: !c.user }">{{ avatarText(c.user) }}</span>
            <div class="content">
              <div class="name">
                {{ c.user ? c.user.nickname : '已注销用户' }}
                <span v-if="c.is_top === 1" class="badge badge-primary" style="margin-left: 6px">📌 置顶</span>
              </div>
              <div class="text" :class="{ deleted: c.deleted }">
                {{ c.deleted ? '该评论已删除' : c.content }}
              </div>
              <div class="foot">
                <span>{{ formatDateTime(c.created_at) }}</span>
                <a v-if="!c.deleted && !replyTo" class="reply-link" @click="startReply(c, null)">回复</a>
                <a v-if="canDelete(c)" class="reply-link danger" @click="removeComment(c.id)">删除</a>
              </div>
            </div>
          </div>

          <!-- 二级回复（最多两层，后端保证） -->
          <div v-if="c.replies && c.replies.length" class="comment-replies">
            <div v-for="r in c.replies" :key="r.id" class="comment-item">
              <span class="avatar" :class="{ gray: !r.user }">{{ avatarText(r.user) }}</span>
              <div class="content">
                <div class="name">
                  {{ r.user ? r.user.nickname : '已注销用户' }}
                  <span v-if="replyName(c, r)" class="reply-arrow">→ {{ replyName(c, r) }}</span>
                </div>
                <div class="text" :class="{ deleted: r.deleted }">
                  {{ r.deleted ? '该评论已删除' : r.content }}
                </div>
                <div class="foot">
                  <span>{{ formatDateTime(r.created_at) }}</span>
                  <a v-if="!r.deleted && !replyTo" class="reply-link" @click="startReply(c, r)">回复</a>
                  <a v-if="canDelete(r)" class="reply-link danger" @click="removeComment(r.id)">删除</a>
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import MarkdownView from '@/components/MarkdownView.vue'
import { getArticle, recordView } from '@/api/article'
import { listComments, createComment, deleteComment } from '@/api/comment'
import { useUserStore } from '@/stores/user'
import { formatDate, formatDateTime, avatarText } from '@/utils/format'

const route = useRoute()
const store = useUserStore()
const articleId = Number(route.params.id)

const article = ref(null)
const draft = ref('')
const submitting = ref(false)

const comments = ref([])
const commentTotal = ref(0)
const commentsLoading = ref(false)
const commentQuery = reactive({ page: 1, page_size: 10 })

// 当前回复目标：{ parentComment(一级), target(实际回复对象或 null) }
const replyTo = ref(null)
let viewTimer = null

const totalComments = computed(() => commentTotal.value)

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
    return
  }
  // 回复二级时 parent 恒为一级（target 为二级时取其一上级 id）；reply_to 为实际对象
  const targetComment = target || top
  replyTo.value = {
    parentCommentId: target && target.parent_id ? target.parent_id : top.id,
    replyToId: target ? target.id : null,
    name: (targetComment.user?.nickname) || (targetComment.user ? '' : '已注销用户')
  }
}

function cancelReply() {
  replyTo.value = null
  draft.value = ''
}

async function submitComment() {
  if (!store.isLoggedIn) {
    ElMessage.warning('请先登录后再评论')
    return
  }
  if (!draft.value.trim()) return
  const payload = { content: draft.value.trim() }
  if (replyTo.value) {
    if (replyTo.value.parentCommentId) payload.parent_id = replyTo.value.parentCommentId
    if (replyTo.value.replyToId) payload.reply_to_id = replyTo.value.replyToId
  }
  submitting.value = true
  try {
    await createComment(articleId, payload)
    ElMessage.success('评论成功')
    draft.value = ''
    cancelReply()
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
  } catch (e) { /* 拦截器提示 */ } finally {
    commentsLoading.value = false
  }
}

onMounted(async () => {
  try {
    article.value = await getArticle(articleId)
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
  if (viewTimer) clearTimeout(viewTimer)
})
</script>

<style scoped>
.article-header { padding: 30px 0 8px; }
.cat-row { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }
.article-header h1 { font-size: 30px; line-height: 1.35; font-weight: 700; margin-bottom: 14px; }
.article-meta { display: flex; gap: 18px; align-items: center; font-size: 13px; color: var(--text-muted); flex-wrap: wrap; }
.avatar { width: 34px; height: 34px; border-radius: 50%; background: var(--primary-100); color: var(--primary-700); display: grid; place-items: center; font-weight: 700; flex-shrink: 0; }
.avatar.gray { background: #e5e7eb; color: #6b7280; }
.article-cover { border-radius: 12px; margin: 18px 0 6px; overflow: hidden; }
.article-cover img { width: 100%; max-height: 380px; object-fit: cover; }
.markdown-body-pad { padding: 18px 0 10px; }

/* —— 评论区（两层级联，后端组装）—— */
.comments { margin-top: 10px; }
.comments-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.comments-head h2 { font-size: 19px; font-weight: 700; }
.comment-form { padding: 16px 18px; }
.comment-form .actions { display: flex; justify-content: space-between; align-items: center; margin-top: 12px; }
.comment-form .login-tip { font-size: 13px; color: var(--text-muted); }
.cancel-reply { color: var(--primary); cursor: pointer; margin-left: 8px; }
.comment-wrap { margin-bottom: 4px; }
.comment-item { display: flex; gap: 12px; padding: 15px 0; border-bottom: 1px solid var(--border); }
.comment-item .avatar { width: 38px; height: 38px; }
.comment-item .content { flex: 1; min-width: 0; }
.comment-item .name { font-weight: 600; font-size: 14px; }
.comment-item .name .reply-arrow { color: var(--text-faint); font-weight: 400; margin-left: 6px; }
.comment-item .text { font-size: 14px; margin: 4px 0 6px; color: var(--text); word-break: break-word; }
.comment-item .text.deleted { color: var(--text-faint); font-style: italic; }
.comment-item .foot { display: flex; gap: 16px; font-size: 12.5px; color: var(--text-faint); align-items: center; }
.comment-item .foot .reply-link { color: var(--text-muted); cursor: pointer; }
.comment-item .foot .reply-link:hover { color: var(--primary); }
.comment-item .foot .reply-link.danger { color: var(--danger); }
.comment-item.top { background: var(--primary-50); border-radius: 8px; padding: 15px 14px; border-bottom: none; margin-top: 8px; }
.comment-replies { margin-left: 50px; border-left: 2px solid var(--border); padding-left: 18px; }
.comment-replies .comment-item { border-bottom: none; padding: 10px 0; }
.pager { display: flex; justify-content: center; margin: 16px 0; }
@media (max-width: 640px) {
  .comment-replies { margin-left: 16px; padding-left: 10px; }
}
</style>
