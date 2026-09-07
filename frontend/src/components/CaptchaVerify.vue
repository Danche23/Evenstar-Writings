<template>
  <!-- 阿里云验证码 2.0 封装：调用 verify() 弹出滑块，滑动成功 resolve 出 captchaVerifyParam -->
  <Teleport to="body">
    <div v-if="visible" class="cap-mask">
      <div class="cap-panel">
        <h3>{{ title }}</h3>
        <p class="cap-tip">拖动滑块完成验证，通过后请尽快提交（凭证 90 秒内有效）</p>
        <div v-if="!configured" class="cap-error">
          未配置滑块参数：请在 frontend/.env.development.local 填入
          VITE_CAPTCHA_PREFIX / VITE_CAPTCHA_SCENE_ID（取自后端 config.yaml 的 captcha 段）
        </div>
        <template v-else>
          <div v-if="loading" class="cap-loading">正在加载验证组件…</div>
          <div :id="mountId" v-show="!loading"></div>
          <div v-if="errorMsg" class="cap-error">{{ errorMsg }}<a class="cap-retry" @click="start">点击重试</a></div>
        </template>
        <div class="cap-actions">
          <button class="btn btn-ghost btn-sm" @click="cancel">取消</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { CAPTCHA_PREFIX, CAPTCHA_SCENE_ID, CAPTCHA_REGION } from '@/utils/captcha'

const SDK_URL = 'https://o.alicdn.com/captcha-frontend/aliyunCaptcha/AliyunCaptcha.js'
let uid = 0

const props = defineProps({
  title: { type: String, default: '请完成滑块验证' }
})
const emit = defineEmits(['success'])

const visible = ref(false)
const loading = ref(false)
const errorMsg = ref('')
const mountId = ref('')
const configured = ref(!!(CAPTCHA_PREFIX && CAPTCHA_SCENE_ID))

let resolver = null
let rejecter = null
let instance = null
let sdkPromise = null

// 官方 SDK 只需加载一次（region + prefix 需在加载前注入全局配置）
function ensureSdk() {
  if (window.initAliyunCaptcha) return Promise.resolve()
  if (!sdkPromise) {
    sdkPromise = new Promise((resolve, reject) => {
      window.AliyunCaptchaConfig = { region: CAPTCHA_REGION, prefix: CAPTCHA_PREFIX }
      const s = document.createElement('script')
      s.src = SDK_URL
      s.onload = () => resolve()
      s.onerror = () => reject(new Error('验证组件加载失败，请检查网络后重试'))
      document.head.appendChild(s)
    })
  }
  return sdkPromise
}

function destroyInstance() {
  if (instance && instance.destroy) {
    try { instance.destroy() } catch (e) { /* ignore */ }
  }
  instance = null
}

function start() {
  if (!configured.value) return
  loading.value = true
  errorMsg.value = ''
  destroyInstance()
  ensureSdk()
    .then(() => nextTick())
    .then(() => {
      window.initAliyunCaptcha({
        SceneId: CAPTCHA_SCENE_ID,
        mode: 'embed',
        element: `#${mountId.value}`,
        success: (data) => {
          const param = typeof data === 'string' ? data : data?.captchaVerifyParam
          if (!param) {
            errorMsg.value = '未获取到验证凭证，请重试'
            return
          }
          loading.value = false
          visible.value = false
          destroyInstance()
          emit('success', param)
          resolver && resolver(param)
        },
        fail: () => {
          loading.value = false
          errorMsg.value = '滑块验证未通过，请重试'
        },
        getInstance: (inst) => { instance = inst },
        onError: (info) => {
          loading.value = false
          errorMsg.value = '验证组件异常：' + (info?.msg || '请重试')
        },
        slideStyle: { width: 320, height: 40 },
        language: 'cn'
      })
      loading.value = false
    })
    .catch((e) => {
      loading.value = false
      errorMsg.value = e?.message || '验证组件加载失败'
    })
}

// 外部调用：const param = await captchaRef.verify()
function verify() {
  return new Promise((resolve, reject) => {
    resolver = resolve
    rejecter = reject
    visible.value = true
    mountId.value = `cap-el-${++uid}`
    nextTick(() => start())
  })
}

function cancel() {
  visible.value = false
  destroyInstance()
  rejecter && rejecter(new Error('已取消滑块验证'))
}

defineExpose({ verify })
</script>

<style scoped>
.cap-mask {
  position: fixed; inset: 0; z-index: 3000;
  background: rgba(17, 24, 39, 0.45);
  display: flex; align-items: center; justify-content: center;
}
.cap-panel {
  background: #fff; border-radius: 12px; width: 420px; max-width: 92vw;
  padding: 22px 24px; box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}
.cap-panel h3 { font-size: 16px; font-weight: 700; margin-bottom: 6px; }
.cap-tip { font-size: 12.5px; color: var(--text-muted); margin-bottom: 14px; }
.cap-loading { padding: 26px 0; text-align: center; color: var(--text-muted); font-size: 13px; }
.cap-error { margin: 8px 0; font-size: 13px; color: var(--danger); }
.cap-retry { color: var(--primary); cursor: pointer; margin-left: 8px; text-decoration: underline; }
.cap-actions { display: flex; justify-content: flex-end; margin-top: 12px; }
</style>
