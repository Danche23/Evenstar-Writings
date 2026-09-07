<template>
  <div class="auth-wrap">
    <div class="card auth-card">
      <h1>注册</h1>
      <div class="sub">注册后可发表评论（管理员不开放注册）</div>

      <el-form label-position="top" @submit.prevent>
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="用于展示的用户名" size="large" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" placeholder="you@example.com" size="large" />
        </el-form-item>
        <el-form-item label="邮箱验证码">
          <div class="code-row">
            <el-input v-model="form.code" placeholder="6 位验证码" size="large" />
            <button class="btn btn-outline" :disabled="countdown > 0 || sending" style="white-space: nowrap" @click="handleSendCode">
              {{ countdown > 0 ? `${countdown}s 后重发` : '发送验证码' }}
            </button>
          </div>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password placeholder="8 位以上，含字母和数字" size="large" />
        </el-form-item>
        <el-form-item label="昵称（选填）">
          <el-input v-model="form.nickname" placeholder="不填将自动生成默认昵称" size="large" />
        </el-form-item>
      </el-form>

      <button class="btn btn-primary btn-block" :disabled="submitting" @click="handleRegister">
        {{ submitting ? '注册中…' : '注册' }}
      </button>
      <div class="auth-footer">已有账号？<router-link to="/login">去登录</router-link></div>
    </div>

    <CaptchaVerify ref="captchaRef" title="请完成滑块验证" />
  </div>
</template>

<script setup>
import { ref, reactive, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import CaptchaVerify from '@/components/CaptchaVerify.vue'
import { sendCode, register } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const store = useUserStore()

const form = reactive({ username: '', email: '', code: '', password: '', nickname: '' })
const submitting = ref(false)
const sending = ref(false)
const countdown = ref(0)
const captchaRef = ref(null)
let timer = null

const emailValid = () => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())

function startCountdown() {
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value -= 1
    if (countdown.value <= 0) clearInterval(timer)
  }, 1000)
}

// 发送验证码：需先过滑块（后端 send-code 每次必验）
async function handleSendCode() {
  if (!emailValid()) {
    ElMessage.warning('请先填写正确的邮箱')
    return
  }
  sending.value = true
  try {
    const param = await captchaRef.value.verify()
    await sendCode(form.email.trim(), 'register', param)
    ElMessage.success('验证码已发送，请查收邮箱（mock 模式下查看后端日志）')
    startCountdown()
  } catch (e) {
    // 拦截器已统一提示后端错误；仅忽略用户主动取消滑块
    if (e?.message === '已取消滑块验证') return
  } finally {
    sending.value = false
  }
}

// 注册：每次必验滑块
async function handleRegister() {
  if (!form.username.trim() || !emailValid() || !form.code.trim()) {
    ElMessage.warning('请完整填写用户名、邮箱与验证码')
    return
  }
  if (!/^(?=.*[A-Za-z])(?=.*\d).{8,}$/.test(form.password)) {
    ElMessage.warning('密码至少 8 位且需包含字母和数字')
    return
  }
  submitting.value = true
  try {
    const param = await captchaRef.value.verify()
    const data = await register({
      username: form.username.trim(),
      password: form.password,
      email: form.email.trim(),
      code: form.code.trim(),
      nickname: form.nickname.trim() || undefined,
      captcha_verify_param: param
    })
    store.setAuth(data.token, data.user)
    ElMessage.success('注册成功，已自动登录')
    router.push(route.query.redirect ? String(route.query.redirect) : '/')
  } catch (e) {
    // 拦截器已统一提示后端错误（含滑块失败 1008）；仅忽略用户主动取消
    if (e?.message === '已取消滑块验证') return
  } finally {
    submitting.value = false
  }
}

onBeforeUnmount(() => timer && clearInterval(timer))
</script>

<style scoped>
.auth-wrap { max-width: 440px; margin: 50px auto; padding: 0 16px; }
.auth-card { padding: 32px 30px; }
.auth-card h1 { font-size: 22px; font-weight: 700; margin-bottom: 4px; }
.auth-card .sub { font-size: 14px; color: var(--text-muted); margin-bottom: 22px; }
.code-row { display: flex; gap: 10px; width: 100%; }
.code-row .el-input { flex: 1; }
.auth-footer { text-align: center; font-size: 13.5px; color: var(--text-muted); margin-top: 18px; }
.auth-footer a { color: var(--primary); font-weight: 500; }
</style>
