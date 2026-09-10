// 后台溢出自动诊断：headless Edge(CDP) 登录 → 逐页扫描越界元素 + 全页截图
// 用法: node ovf_scan.mjs   (需先启动 msedge --remote-debugging-port=9223)
import fs from 'node:fs'

const CDP_HTTP = 'http://127.0.0.1:9223'
const BASE = 'http://localhost:5173'
const OUTDIR = 'D:/a_lezhijiaoyu/_ovf_shots'
fs.mkdirSync(OUTDIR, { recursive: true })

const ACCOUNT = { email: 'evenstar@evenstar.local', password: 'MuXing@2026' }

const sleep = (ms) => new Promise((r) => setTimeout(r, ms))

async function createTarget(url) {
  const res = await fetch(`${CDP_HTTP}/json/new?${encodeURIComponent(url)}`, { method: 'PUT' })
  return res.json()
}

function connect(wsUrl) {
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(wsUrl)
    let id = 0
    const pending = new Map()
    const listeners = []
    ws.onopen = () => resolve({
      send(method, params = {}) {
        return new Promise((res, rej) => {
          const mid = ++id
          pending.set(mid, { res, rej })
          ws.send(JSON.stringify({ id: mid, method, params }))
        })
      },
      onEvent(fn) { listeners.push(fn) },
      close() { ws.close() }
    })
    ws.onerror = (e) => reject(new Error('ws error ' + e.message))
    ws.onmessage = (ev) => {
      const msg = JSON.parse(ev.data)
      if (msg.id && pending.has(msg.id)) {
        const p = pending.get(msg.id)
        pending.delete(msg.id)
        msg.error ? p.rej(new Error(msg.error.message)) : p.res(msg.result)
      } else if (msg.method) {
        listeners.forEach((fn) => fn(msg))
      }
    }
  })
}

async function evalJs(cdp, expression) {
  const r = await cdp.send('Runtime.evaluate', {
    expression, awaitPromise: true, returnByValue: true
  })
  if (r.exceptionDetails) throw new Error('eval err: ' + JSON.stringify(r.exceptionDetails).slice(0, 300))
  return r.result.value
}

const SCAN = `(async () => {
  const vw = document.documentElement.clientWidth
  const hasScrollableAnc = (el) => {
    let n = el.parentElement
    while (n) {
      const o = getComputedStyle(n)
      if (o.overflowX === 'auto' || o.overflowX === 'scroll' || o.overflow === 'auto' || o.overflow === 'scroll') return true
      n = n.parentElement
    }
    return false
  }
  const out = []
  document.querySelectorAll('body *').forEach((el) => {
    const r = el.getBoundingClientRect()
    if (r.width > vw + 1 && r.right > vw + 1) {
      if (!hasScrollableAnc(el)) {
        el.style.outline = '3px dashed #e53935'
        const cls = String(el.className || '').trim().split(/\\s+/).slice(0, 2).join('.')
        out.push((el.tagName.toLowerCase()) + (cls ? '.' + cls : '') + ' w=' + Math.round(r.width))
      }
    }
  })
  return JSON.stringify({ vw, total: out.length, list: out.slice(0, 25),
    bodyScrollW: document.body ? document.body.scrollWidth : null })
})()`

async function capture(cdp, file) {
  const shot = await cdp.send('Page.captureScreenshot', { format: 'png', captureBeyondViewport: true })
  fs.writeFileSync(file, Buffer.from(shot.data, 'base64'))
}

async function main() {
  // 等待 CDP 可用
  for (let i = 0; i < 20; i++) {
    try { await fetch(`${CDP_HTTP}/json/version`); break } catch { await sleep(500) }
  }
  const t = await createTarget('about:blank')
  const cdp = await connect(t.webSocketDebuggerUrl)
  await cdp.send('Page.enable')
  await cdp.send('Runtime.enable')
  cdp.onEvent((m) => {
    if (m.method === 'Runtime.consoleAPICalled' && m.params.type === 'error') {
      const txt = m.params.args.map((a) => a.value ?? a.description ?? '').join(' ')
      console.log('[page-console-error]', txt.slice(0, 200))
    }
  })

  // 1) 打开站点并登录拿 token（经 vite 代理到后端）
  await cdp.send('Page.navigate', { url: BASE + '/login' })
  await sleep(2500)
  const loginRes = await evalJs(cdp, `(async () => {
    try {
      const r = await fetch('/api/auth/login', { method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: ${JSON.stringify(ACCOUNT.email)}, password: ${JSON.stringify(ACCOUNT.password)} }) })
      const j = await r.json()
      return JSON.stringify({ ok: r.ok, status: r.status, body: j })
    } catch (e) { return 'NET_ERR ' + e.message }
  })()`)
  console.log('[login]', loginRes)
  let token = '', user = null
  try {
    const parsed = JSON.parse(loginRes)
    if (parsed.body && parsed.body.code === 0 && parsed.body.data) {
      token = parsed.body.data.token
      user = parsed.body.data.user
    } else if (parsed.body && parsed.body.token) { token = parsed.body.token; user = parsed.body.user }
  } catch {}
  if (!token) { console.log('!! 登录失败，无法继续'); return }
  await evalJs(cdp, `localStorage.setItem('evenstar_token', ${JSON.stringify(token)});
    localStorage.setItem('evenstar_user', ${JSON.stringify(JSON.stringify(user))}); 'ok'`)

  const routes = ['/admin', '/admin/articles', '/admin/comments', '/admin/categories',
    '/admin/categories/sort', '/admin/tags', '/admin/tags/sort', '/admin/users',
    '/admin/uploads', '/admin/about']
  const sizes = [
    { name: 'mobile', w: 390, h: 844 },
    { name: 'desktop', w: 1360, h: 900 }
  ]
  for (const s of sizes) {
    await cdp.send('Emulation.setDeviceMetricsOverride', {
      width: s.w, height: s.h, deviceScaleFactor: 1, mobile: s.name === 'mobile'
    })
    for (const route of routes) {
      await cdp.send('Page.navigate', { url: BASE + route })
      await sleep(2600)
      const scan = await evalJs(cdp, SCAN)
      const info = JSON.parse(scan)
      const line = `[${s.name}] ${route} vw=${info.vw} 越界=${info.total}  bodyScrollW=${info.bodyScrollW}`
      console.log(line)
      if (info.total > 0) {
        console.log('   ' + info.list.join('\n   '))
      }
      await capture(cdp, `${OUTDIR}/${s.name}${route.replace(/\//g, '_') || '_home'}.png`)
    }
  }
  cdp.close()
  console.log('DONE → ' + OUTDIR)
}

main().catch((e) => { console.error('FATAL', e); process.exit(1) })
