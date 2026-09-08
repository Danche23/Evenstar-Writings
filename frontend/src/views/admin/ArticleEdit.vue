<template>
  <div class="editor-page">
    <router-link to="/admin/articles" class="back-link">← 返回文章列表</router-link>

    <div class="admin-page-head flex-between">
      <h2>{{ isEdit ? '编辑文章' : '写文章' }}</h2>
      <div class="status-switch">
        <span class="badge" :class="status === 2 ? 'badge-success' : 'badge-gray'">
          {{ status === 2 ? '发布状态' : '草稿状态' }}
        </span>
      </div>
    </div>

    <div class="card form-card">
      <!-- 标题 -->
      <div class="field">
        <label class="form-label">标题 <span class="danger-star">*</span></label>
        <el-input v-model="form.title" maxlength="200" placeholder="请输入文章标题" size="large" />
      </div>

      <!-- 摘要 -->
      <div class="field">
        <label class="form-label">摘要（选填，用于列表展示）</label>
        <el-input v-model="form.summary" type="textarea" :rows="2" maxlength="500" placeholder="一句话概括文章内容" />
      </div>

      <!-- 封面 -->
      <div class="field">
        <label class="form-label">封面图（选填，≤5MB）</label>
        <div class="cover-row">
          <img v-if="form.cover" :src="form.cover" class="cover-preview" alt="cover" />
          <el-upload
            :show-file-list="false"
            accept="image/jpeg,image/png,image/gif,image/webp"
            :http-request="onCoverUpload"
            :disabled="uploading"
          >
            <button class="btn btn-outline btn-sm" :disabled="uploading">
              {{ uploading ? '上传中…' : form.cover ? '更换封面' : '上传封面' }}
            </button>
          </el-upload>
          <button v-if="form.cover" class="btn btn-ghost btn-sm" @click="form.cover = ''">移除封面</button>
        </div>
        <div class="hint">留空则前台显示默认封面占位</div>
      </div>

      <!-- 分类 / 标签 多选 -->
      <div class="multi-row">
        <div class="field">
          <label class="form-label">分类</label>
          <el-select v-model="form.category_ids" multiple placeholder="选择分类" style="width: 100%">
            <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </div>
        <div class="field">
          <label class="form-label">标签</label>
          <el-select v-model="form.tag_ids" multiple placeholder="选择标签" style="width: 100%">
            <el-option v-for="t in tags" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
        </div>
      </div>
    </div>

    <!-- Markdown 正文编辑器 -->
    <div class="field mt-16">
      <label class="form-label">正文（Markdown）<span class="danger-star">*</span></label>
      <MdEditor
        v-model="form.content"
        style="height: 520px"
        :toolbars-exclude="['htmlPreview']"
        @onUploadImg="onUploadImg"
      />
    </div>

    <!-- 底部操作 -->
    <div class="actions-bar">
      <router-link to="/admin/articles" class="btn btn-ghost">取消</router-link>
      <button class="btn btn-outline" :disabled="saving" @click="save(1)">保存草稿</button>
      <button class="btn btn-primary" :disabled="saving" @click="save(2)">发布</button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { adminGetArticle, adminCreateArticle, adminUpdateArticle } from '@/api/article'
import { listCategories } from '@/api/category'
import { listTags } from '@/api/tag'
import { uploadImage } from '@/api/upload'

const route = useRoute()
const router = useRouter()
const id = computed(() => (route.params.id ? Number(route.params.id) : null))
const isEdit = computed(() => !!id.value)

const form = reactive({
  title: '',
  summary: '',
  content: '',
  cover: '',
  category_ids: [],
  tag_ids: []
})
const status = ref(1) // 1=草稿 2=发布（决定按钮文案/状态显示）
const saving = ref(false)
const uploading = ref(false)
const categories = ref([])
const tags = ref([])

// —— 正文内插入图片：走后端中转上传（scene=article，仅管理员）——
async function onUploadImg(files, callback) {
  try {
    const urls = []
    for (const file of files) {
      const data = await uploadImage(file, 'article')
      urls.push(data.url)
    }
    callback(urls)
  } catch (e) {
    ElMessage.error(e.message || '图片上传失败')
    callback([])
  }
}

// —— 封面上传 ——
async function onCoverUpload({ file }) {
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.warning('图片不能超过 5MB')
    return
  }
  uploading.value = true
  try {
    const data = await uploadImage(file, 'article')
    form.cover = data.url
    ElMessage.success('封面上传成功')
  } catch (e) { /* 拦截器提示 */ } finally {
    uploading.value = false
  }
}

async function loadOptions() {
  try {
    const [cs, ts] = await Promise.all([listCategories(), listTags()])
    categories.value = cs || []
    tags.value = ts || []
  } catch (e) { /* ignore */ }
}

async function loadArticle() {
  try {
    const d = await adminGetArticle(id.value)
    form.title = d.title || ''
    form.summary = d.summary || ''
    form.content = d.content || ''
    form.cover = d.cover || ''
    form.category_ids = d.category_ids || []
    form.tag_ids = d.tag_ids || []
    status.value = d.status || 1
  } catch (e) {
    ElMessage.error(e.message || '加载文章失败')
    router.push('/admin/articles')
  }
}

async function save(s) {
  if (!form.title.trim()) {
    ElMessage.warning('请填写标题')
    return
  }
  if (!form.content.trim()) {
    ElMessage.warning('请填写正文')
    return
  }
  saving.value = true
  const payload = {
    title: form.title.trim(),
    summary: form.summary.trim() || null,
    content: form.content,
    cover: form.cover || null,
    status: s,
    category_ids: form.category_ids,
    tag_ids: form.tag_ids
  }
  try {
    if (isEdit.value) {
      await adminUpdateArticle(id.value, payload)
    } else {
      await adminCreateArticle(payload)
    }
    ElMessage.success(s === 2 ? '文章已发布' : '草稿已保存')
    router.push('/admin/articles')
  } catch (e) { /* 拦截器提示 */ } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await loadOptions()
  if (isEdit.value) await loadArticle()
})
</script>

<style scoped>
.editor-page { max-width: 920px; margin: 0 auto; }
.form-card { padding: 22px 24px; }
.multi-row { display: grid; grid-template-columns: 1fr 1fr; gap: 18px; }
.cover-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.cover-preview { width: 220px; height: 110px; object-fit: cover; border-radius: 8px; border: 1px solid var(--border); }
.hint { font-size: 12px; color: var(--text-faint); margin-top: 6px; }
.danger-star { color: var(--danger); }
.actions-bar { display: flex; justify-content: flex-end; gap: 12px; margin-top: 20px; }
@media (max-width: 800px) {
  .multi-row { grid-template-columns: 1fr; }
}
</style>
