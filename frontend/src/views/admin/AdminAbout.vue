<template>
  <div>
    <div class="admin-page-head flex-between">
      <div>
        <h2>关于设置</h2>
        <p class="sub">编辑前台「关于」页的文字内容：站点介绍（支持 Markdown）、技术栈标签、文末小注。保存后前台立即生效。</p>
      </div>
      <a href="/about" target="_blank" class="btn btn-outline btn-sm">预览前台关于页 ↗</a>
    </div>

    <div class="admin-card about-form" v-loading="loading">
      <el-form label-position="top" class="about-fields">
        <el-form-item label="站点介绍">
          <el-input
            v-model="doc.intro"
            type="textarea"
            :rows="7"
            placeholder="介绍你的站点。支持 Markdown：空行分段、**加粗**、[链接](url)。"
          />
          <p class="hint">前端会按 Markdown 渲染，空行即分段。</p>
        </el-form-item>

        <el-form-item label="技术栈标签">
          <el-input
            v-model="stackText"
            type="textarea"
            :rows="3"
            placeholder="每行一个，如：&#10;Go&#10;Gin&#10;MySQL"
          />
          <p class="hint">每行一个标签，展示为一行小胶囊。</p>
        </el-form-item>

        <el-form-item label="文末小注">
          <el-input
            v-model="doc.footnote"
            type="textarea"
            :rows="2"
            maxlength="300"
            show-word-limit
            placeholder="关于站名的一句话小注（显示在最下方）"
          />
        </el-form-item>
      </el-form>

      <div class="form-foot">
        <button class="btn btn-primary" :disabled="saving" @click="save">
          {{ saving ? '保存中…' : '保存修改' }}
        </button>
        <button class="btn btn-ghost" :disabled="saving" @click="reload">放弃修改</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminGetAbout, adminSaveAbout } from '@/api/about'

const doc = ref({ intro: '', stack: [], footnote: '' })
const stackText = ref('')
const loading = ref(false)
const saving = ref(false)

function syncStackText() {
  stackText.value = (doc.value.stack || []).join('\n')
}

function parseStack() {
  return stackText.value.split('\n').map((s) => s.trim()).filter(Boolean)
}

async function load() {
  loading.value = true
  try {
    doc.value = (await adminGetAbout()) || { intro: '', stack: [], footnote: '' }
    syncStackText()
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    const body = { ...doc.value, stack: parseStack() }
    doc.value = await adminSaveAbout(body)
    syncStackText()
    ElMessage.success('已保存，前台关于页已同步')
  } catch (e) { /* 拦截器提示 */ } finally {
    saving.value = false
  }
}

function reload() {
  load()
  ElMessage.info('已还原为服务器上的内容')
}

onMounted(load)
</script>

<style scoped>
.about-form { padding: 24px 26px; }
.about-fields { max-width: 760px; }
.about-fields :deep(.el-form-item__label) { font-weight: 600; }
.hint { font-size: 12px; color: var(--text-faint); margin-top: 4px; line-height: 1.6; }
.form-foot { display: flex; gap: 10px; margin-top: 4px; }
</style>
