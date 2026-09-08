import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'

import App from './App.vue'
import router from './router'
import './styles/global.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus, { locale: zhCn })

app.mount('#app')

/* ============ 开发期溢出诊断（不影响生产） ============
   用法：访问 http://localhost:5173/?ovf=1（各页面分别开），
   超宽元素会被红色描边，控制台输出 [OVF] 前缀的精确元素清单。 */
if (import.meta.env.DEV && new URLSearchParams(window.location.search).has('ovf')) {
  const scanOverflow = () => {
    const vw = document.documentElement.clientWidth
    const offenders = []
    document.querySelectorAll('body *').forEach((el) => {
      const w = el.getBoundingClientRect().width
      if (w > vw + 1) {
        el.style.outline = '2px solid #e53935'
        const cls = String(el.className || '')
          .split(/\s+/).slice(0, 3).join('.')
        offenders.push(`${el.tagName.toLowerCase()}.${cls} → ${Math.round(w)}px (视口 ${vw}px)`)
      }
    })
    console.log(`[OVF] 视口=${vw}px  溢出元素数=${offenders.length}`)
    offenders.slice(0, 60).forEach((s) => console.warn('[OVF]', s))
  }
  window.addEventListener('load', () => setTimeout(scanOverflow, 600))
  window.addEventListener('load', () => setTimeout(scanOverflow, 2000))
}
