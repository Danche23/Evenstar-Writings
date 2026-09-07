<template>
  <div class="profile-wrap" v-if="user">
    <!-- 头部信息 -->
    <div class="card profile-head">
      <el-avatar :size="80" :src="avatarUrl || undefined" class="profile-avatar">
        {{ displayName.slice(0, 1) }}
      </el-avatar>
      <div class="info">
        <h2>{{ user.nickname || user.username }}</h2>
        <div class="email">{{ user.email }}</div>
        <div class="roles mt-8">
          <span class="badge" :class="user.role === 1 ? 'badge-primary' : 'badge-gray'">
            {{ user.role === 1 ? '管理员' : '普通用户' }}
          </span>
        </div>
      </div>
    </div>

    <!-- 基本资料 -->
    <div class="card profile-section">
      <h3>基本资料</h3>
      <div class="field">
        <label>昵称</label>
        <el-input v-model="nickname" maxlength="50" placeholder="对外展示的昵称" />
      </div>
      <div class="field">
        <label>头像</label>
        <el-upload
          :show-file-list="false"
          accept="image/jpeg,image/png,image/gif,image/webp"
          :http-request="onAvatarUpload"
          :disabled="uploading"
        >
          <button class="btn btn-outline btn-sm" :disabled="uploading">
            {{ uploading ? '上传中…' : avatarUrl ? '更换头像' : '上传头像' }}
          </button>
        </el-upload>
        <div class="hint">支持 jpg / png / gif / webp，≤ 5MB；每人每月最多更换 5 次（新头像需点「保存资料」生效）</div>
      </div>
      <div class="field">
        <label>邮箱</label>
        <el-input :model-value="user.email" disabled />
        <div class="hint">邮箱不可修改（只读）</div>
      </div>
      <button class="btn btn-primary" :disabled="saving" @click="saveProfile">
        {{ saving ? '保存中…' : '保存资料' }}
      </button>
    </div>

    <!-- 修改密码 -->
    <div class="card profile-section">
      <h3>修改密码</h3>
      <div class="field">
        <label>当前密码</label>
        <el-input v-model="pwd.old_password" type="password" show-password />
      </div>
      <div class="field">
        <label>新密码</label>
        <el-input v-model="pwd.new_password" type="password" show-password placeholder="8 位以上，含字母和数字" />
      </div>
      <div class="field">
        <label>确认新密码</label>
        <el-input v-model="pwd.confirm" type="password" show-password />
      </div>
      <div class="hint" style="margin-bottom: 12px">修改成功后所有设备将退出登录，需重新登录</div>
      <button class="btn btn-outline" :disabled="savingPwd" @click="savePassword">
        {{ savingPwd ? '提交中…' : '修改密码' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { updateProfile, updatePassword } from '@/api/user'
import { uploadImage } from '@/api/upload'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const store = useUserStore()
const user = computed(() => store.user)
const displayName = computed(() => store.displayName || '?')

const nickname = ref(user.value?.nickname || '')
const avatarUrl = ref(user.value?.avatar || '')
const saving = ref(false)
const uploading = ref(false)

const pwd = ref({ old_password: '', new_password: '', confirm: '' })
const savingPwd = ref(false)

// —— 头像上传：el-upload 自定义请求 → 后端中转 OSS ——
async function onAvatarUpload({ file }) {
  uploading.value = true
  try {
    if (file.size > 5 * 1024 * 1024) {
      ElMessage.warning('图片不能超过 5MB')
      return
    }
    const data = await uploadImage(file, 'avatar')
    avatarUrl.value = data.url
    ElMessage.success('头像上传成功，点击「保存资料」生效')
  } catch (e) {
    /* 拦截器已提示（如每月次数超限） */
  } finally {
    uploading.value = false
  }
}

async function saveProfile() {
  saving.value = true
  try {
    await updateProfile({ nickname: nickname.value.trim(), avatar: avatarUrl.value || null })
    store.setUser({ ...user.value, nickname: nickname.value.trim() || user.value.nickname, avatar: avatarUrl.value || user.value.avatar })
    ElMessage.success('资料已保存')
  } catch (e) {
    /* 拦截器提示 */
  } finally {
    saving.value = false
  }
}

async function savePassword() {
  if (!pwd.value.old_password || !pwd.value.new_password) {
    ElMessage.warning('请填写当前密码与新密码')
    return
  }
  if (!/^(?=.*[A-Za-z])(?=.*\d).{8,}$/.test(pwd.value.new_password)) {
    ElMessage.warning('新密码至少 8 位且需包含字母和数字')
    return
  }
  if (pwd.value.new_password !== pwd.value.confirm) {
    ElMessage.warning('两次输入的新密码不一致')
    return
  }
  savingPwd.value = true
  try {
    await updatePassword({ old_password: pwd.value.old_password, new_password: pwd.value.new_password })
    ElMessage.success('密码已修改，请重新登录')
    store.logout()
    router.push('/login')
  } catch (e) {
    /* 拦截器提示（如旧密码错误 1010） */
  } finally {
    savingPwd.value = false
  }
}
</script>

<style scoped>
.profile-wrap { max-width: 720px; margin: 40px auto; padding: 0 20px; }
.profile-head { display: flex; align-items: center; gap: 22px; padding: 24px 26px; margin-bottom: 20px; }
.profile-avatar { background: linear-gradient(135deg, var(--primary-100), #c7d2fe); color: var(--primary-700); font-size: 30px; font-weight: 700; }
.info h2 { font-size: 20px; font-weight: 700; }
.info .email { font-size: 13px; color: var(--text-muted); }
.profile-section { padding: 22px 26px; margin-bottom: 20px; }
.profile-section h3 { font-size: 16px; font-weight: 700; margin-bottom: 18px; }
.hint { font-size: 12px; color: var(--text-faint); margin-top: 6px; }
</style>
