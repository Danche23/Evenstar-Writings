<template>
  <div class="auth-wrap">
    <div class="card auth-card">
      <h1>登录</h1>
      <div class="sub">欢迎回来，请登录你的账号</div>

      <el-form label-position="top" @submit.prevent>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" placeholder="you@example.com" size="large" @keyup.enter="submit" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" size="large"
            @keyup.enter="submit" />
        </el-form-item>
      </el-form>

      <!-- 连续失败 3 次后显示提示，提交时自动弹滑块 -->
      <div v-if="needCaptcha" class="captcha-note-box">
        ⚠️ 连续登录失败次数过多，登录前需要完成一次滑块验证
      </div>

      <button class="btn btn-primary btn-block" :disabled="submitting" @click="submit">
        {{ submitting ? '登录中…' : '登录' }}
      </button>

      <div class="auth-footer">
        <router-link to="/forgot" class="muted">忘记密码？</router-link>
        &nbsp;·&nbsp; 还没有账号？<router-link to="/register">去注册</router-link>
      </div>
    </div>

    <!-- 滑块验证（登录被锁时自动弹出） -->
    <CaptchaVerify ref="captchaRef" title="请完成滑块验证后继续登录" />
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import CaptchaVerify from '@/components/CaptchaVerify.vue'
import { login } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const store = useUserStore()

const form = reactive({ email: '', password: '' })
const submitting = ref(false)
const needCaptcha = ref(false)
const captchaRef = ref(null)
let captchaAttempts = 0

function goHome() {
  const redirect = route.query.redirect
  const target = redirect ? String(redirect) : (store.isAdmin ? '/admin' : '/')
  router.push(target)
}

async function doLogin(param) {
  submitting.value = true
  try {
    const data = await login(form.email, form.password, param)
    store.setAuth(data.token, data.user)
    needCaptcha.value = false
    captchaAttempts = 0
    ElMessage.success('登录成功')
    goHome()
    return true
  } catch (e) {
    // login 接口错误已在拦截器静默，这里按业务码处理
    if (e.code === 1007) {
      needCaptcha.value = true
      ElMessage.warning('连续失败次数过多，请先完成滑块验证')
      await tryCaptchaLogin()
    } else if (e.code === 1008) {
      ElMessage.warning('滑块验证失败或已过期，请重新滑动')
      await tryCaptchaLogin()
    } else if (e.code === 429) {
      ElMessage.warning('尝试过于频繁，请稍后再试')
    } else {
      ElMessage.error(e.message || '登录失败')
    }
    return false
  } finally {
    submitting.value = false
  }
}

// 弹出滑块 → 成功后带 captcha_verify_param 重新登录（凭证 90s 内有效）
async function tryCaptchaLogin() {
  captchaAttempts += 1
  if (captchaAttempts > 3) {
    captchaAttempts = 0
    ElMessage.warning('验证次数过多，请刷新后重试')
    return
  }
  try {
    const param = await captchaRef.value.verify()
    await doLogin(param)
  } catch (e) {
    if (e?.message !== '已取消滑块验证') {
      ElMessage.error(e?.message || '滑块验证失败')
    }
    captchaAttempts = 0
  }
}

function submit() {
  if (!form.email || !form.password) {
    ElMessage.warning('请填写邮箱与密码')
    return
  }
  doLogin()
}
</script>

<style scoped>
.auth-wrap { max-width: 420px; margin: 60px auto; padding: 0 16px; }
.auth-card { padding: 32px 30px; }
.auth-card h1 { font-size: 22px; font-weight: 700; margin-bottom: 4px; }
.auth-card .sub { font-size: 14px; color: var(--text-muted); margin-bottom: 22px; }
.captcha-note-box {
  background: var(--warning-50); color: var(--warning);
  font-size: 13px; padding: 10px 14px; border-radius: 8px; margin-bottom: 14px;
}
.auth-footer { text-align: center; font-size: 13.5px; color: var(--text-muted); margin-top: 18px; }
.auth-footer a { color: var(--primary); font-weight: 500; }
</style>
