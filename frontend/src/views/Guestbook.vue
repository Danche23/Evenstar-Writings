<template>
  <div class="page-container page-block">
    <header class="gb-head">
      <h1>留言板</h1>
      <p class="sub poem">没有主题的一张墙——想说的话，都可以留在这里。</p>
    </header>

    <!-- 写留言 -->
    <div class="card gb-composer">
      <template v-if="store.isLoggedIn">
        <div class="gc-bar">
          <UserAvatar :user="store.user" :size="34" />
          <b class="gc-name">{{ store.displayName }}</b>
        </div>
        <el-input
          v-model="draft" type="textarea" :rows="3" maxlength="400" show-word-limit
          placeholder="留下点什么吧：一句话、一个问题、或近来的心情…"
        />
        <div class="gc-foot">
          <span class="hint">友善发言 · 最多 400 字</span>
          <button class="send-btn" :disabled="sending || !draft.trim()" @click="submit">
            {{ sending ? '放飞中…' : '❀ 放飞便签' }}
          </button>
        </div>
      </template>
      <div v-else class="gc-guest">
        <span class="gc-mark">✉️</span>
        <div class="gc-copy"><b>登录后即可留言</b><p>与博主和其他读者在留言板聊聊吧</p></div>
        <button class="send-btn ghost" @click="goLogin">去登录</button>
      </div>
    </div>

    <div v-if="loading" class="loading-box">墙上正在风干墨水…</div>
    <div v-else-if="!list.length" class="empty"><div class="big">🪶</div>墙上还空着，来写第一张吧</div>

    <div v-else class="gb-list">
      <div v-for="m in list" :key="m.id" class="gb-item">
        <UserAvatar :user="m.user" :size="40" />
        <div class="gb-body">
          <div class="gb-bar">
            <b>{{ m.user?.nickname || '已注销用户' }}</b>
            <span class="when">{{ formatDateTime(m.created_at) }}</span>
            <button v-if="canDelete(m)" class="op-btn danger" @click="remove(m)">删除</button>
          </div>
          <p class="gb-text">{{ m.content }}</p>
        </div>
      </div>

      <div v-if="total > query.page_size" class="pager">
        <el-pagination
          layout="prev, pager, next" :total="total" :page-size="query.page_size"
          :current-page="query.page"
          @current-change="(p) => { query.page = p; load() }"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import UserAvatar from '@/components/UserAvatar.vue'
import { listMessages, createMessage, deleteMessage } from '@/api/message'
import { useUserStore } from '@/stores/user'
import { formatDateTime } from '@/utils/format'

const router = useRouter()
const route = useRoute()
const store = useUserStore()

const list = ref([])
const total = ref(0)
const loading = ref(true)
const sending = ref(false)
const draft = ref('')
const query = reactive({ page: 1, page_size: 10 })

function canDelete(m) {
  return store.isAdmin || (store.user && m.user_id === store.user.id)
}
function goLogin() {
  router.push({ name: 'login', query: { redirect: route.fullPath } })
}

async function load() {
  loading.value = true
  try {
    const d = await listMessages({ page: query.page, page_size: query.page_size })
    list.value = d.list || []
    total.value = Number(d.total || 0)
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
}

async function submit() {
  sending.value = true
  try {
    await createMessage(draft.value.trim())
    ElMessage.success('留言已贴上墙')
    draft.value = ''
    query.page = 1
    load()
  } catch (e) { /* 拦截器提示 */ } finally {
    sending.value = false
  }
}

async function remove(m) {
  try {
    await deleteMessage(m.id)
    ElMessage.success('已删除')
    load()
  } catch (e) { /* 拦截器提示 */ }
}

onMounted(load)
</script>

<style scoped>
.gb-head { margin: 10px 0 22px; padding-bottom: 14px; border-bottom: 1px solid var(--border); max-width: 760px; margin-left: auto; margin-right: auto; }
.gb-head h1 { font-size: 30px; }
.gb-head .sub { font-size: 14px; color: var(--text-muted); margin-top: 6px; }
.gb-composer { padding: 18px 20px 16px; border-radius: 14px; margin: 0 auto 22px; max-width: 760px; box-shadow: var(--shadow-sm); }
.gc-bar { display: flex; align-items: center; gap: 10px; margin-bottom: 12px; }
.gc-name { font-size: 14px; }
.gb-composer :deep(.el-textarea__inner) {
  border-radius: 10px; padding: 10px 12px; font-size: 14px; line-height: 1.7;
  background: var(--bg); border-color: var(--border);
  transition: border-color .15s, box-shadow .15s, background .15s;
}
.gb-composer :deep(.el-textarea__inner:focus) { background: var(--surface); border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-100); }
.gc-foot { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 12px; }
.hint { font-size: 12.5px; color: var(--text-faint); }
.send-btn {
  padding: 8px 20px; border: none; border-radius: 999px; background: var(--primary); color: #fff;
  font-size: 13.5px; font-weight: 500; box-shadow: 0 4px 12px rgba(70, 99, 156, .22);
  transition: background .18s, transform .18s;
}
.send-btn:hover:not(:disabled) { background: var(--primary-600); transform: translateY(-1px); }
.send-btn:disabled { opacity: .45; cursor: not-allowed; }
.send-btn.ghost { background: #fff; color: var(--primary); border: 1px solid var(--border-strong); box-shadow: none; }
.send-btn.ghost:hover { background: var(--primary-50); color: var(--primary-600); }
.gc-guest { display: flex; align-items: center; gap: 14px; }
.gc-mark { width: 42px; height: 42px; border-radius: 12px; background: var(--primary-50); display: grid; place-items: center; font-size: 19px; }
.gc-copy { flex: 1; }
.gc-copy b { font-size: 14.5px; display: block; }
.gc-copy p { font-size: 12.5px; color: var(--text-muted); margin: 2px 0 0; }

.gb-list { display: flex; flex-direction: column; max-width: 760px; margin: 0 auto; }
.gb-item {
  display: flex; gap: 12px; padding: 16px 4px;
  border-bottom: 1px solid var(--border);
}
.gb-body { flex: 1; min-width: 0; }
.gb-bar { display: flex; align-items: baseline; gap: 10px; flex-wrap: wrap; }
.gb-bar b { font-size: 14px; color: #24345a; }
.gb-bar .when { font-size: 12px; color: var(--text-faint); }
.op-btn { margin-left: auto; border: none; background: none; padding: 2px 8px; font-size: 12.5px; color: var(--text-muted); border-radius: 8px; cursor: pointer; }
.op-btn:hover { color: var(--danger); background: var(--danger-50); }
.gb-text {
  margin: 6px 0 0; font-size: 14px; line-height: 1.85; color: var(--text);
  white-space: pre-wrap; word-break: break-word;
}
.pager { display: flex; justify-content: center; margin: 18px 0 6px; }
</style>
