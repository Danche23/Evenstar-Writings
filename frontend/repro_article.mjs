// 复现：1) 表格操作列是否可达  2) 草稿文章打开时的控制台报错
import fs from 'node:fs'

const CDP_HTTP = 'http://127.0.0.1:9223'
const BASE = 'http://localhost:5173'
const OUT = 'D:/a_lezhijiaoyu/_ovf_shots/repro'
fs.mkdirSync(OUT, { recursive: true })
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
    const evts = []
    ws.onopen = () => resolve({
      send(method, params = {}) {
        return new Promise((res, rej) => {
          const mid = ++id
          pending.set(mid, { res, rej })
          ws.send(JSON.stringify({ id: mid, method, params }))
        })
      },
      events: evts,
      close() { ws.close() }
    })
    ws.onerror = (e) => reject(new Error('ws error ' + e.message))
    ws.onmessage = (ev) => {
      const msg = JSON.parse(ev.data)
      if (msg.id && pending.has(msg.id)) {
        const p = pending.get(msg.id)
        pending.delete(msg.id)
        msg.error ? p.rej(new Error(msg.error.message)) : p.res(msg.result)
      } else if (msg.method) evts.push(msg)
    }
  })
}
async function evalJs(cdp, expression) {
  const r = await cdp.send('Runtime.evaluate', { expression, awaitPromise: true, returnByValue: true })
  if (r.exceptionDetails) throw new Error('eval err: ' + JSON.stringify(r.exceptionDetails).slice(0, 400))
  return r.result.value
}
function consoleLines(events, since = 0) {
  const out = []
  events.forEach((m) => {
    if (m.method === 'Runtime.consoleAPICalled' && m.params.type === 'error') {
      out.push('[err] ' + m.params.args.map((a) => a.value ?? a.description ?? '').join(' ').slice(0, 300))
    }
    if (m.method === 'Runtime.exceptionThrown') {
      const d = m.params.exceptionDetails
      out.push('[exc] ' + JSON.stringify(d.text + ' ' + (d.exception?.description || '')).slice(0, 400))
    }
  })
  return out
}
async function cap(cdp, file) {
  const s = await cdp.send('Page.captureScreenshot', { format: 'png', captureBeyondViewport: true })
  fs.writeFileSync(file, Buffer.from(s.data, 'base64'))
}

async function main() {
  for (let i = 0; i < 20; i++) { try { await fetch(`${CDP_HTTP}/json/version`); break } catch { await sleep(500) } }
  const t = await createTarget('about:blank')
  const cdp = await connect(t.webSocketDebuggerUrl)
  await cdp.send('Page.enable')
  await cdp.send('Runtime.enable')

  await cdp.send('Page.navigate', { url: BASE + '/login' }); await sleep(2000)
  const loginRes = await evalJs(cdp, `(async () => {
    const r = await fetch('/api/auth/login', { method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: 'evenstar@evenstar.local', password: 'MuXing@2026' }) })
    const j = await r.json()
    return JSON.stringify(j)
  })()`)
  const loginParsed = JSON.parse(loginRes)
  if (!loginParsed || loginParsed.code !== 0 || !loginParsed.data) { console.log('登录失败', loginRes); cdp.close(); process.exit(1) }
  const token = loginParsed.data.token
  const user = loginParsed.data.user
  await evalJs(cdp, `localStorage.setItem('evenstar_token', ${JSON.stringify(token)});
    localStorage.setItem('evenstar_user', ${JSON.stringify(JSON.stringify(user))}); 'ok'`)

  // 找到标题含「从零实现」的文章（草稿/发布都可）
  const arts = JSON.parse(await evalJs(cdp, `(async () => {
    const r = await fetch('/api/admin/articles?page=1&page_size=50', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('evenstar_token') } })
    const j = await r.json()
    return JSON.stringify(j)
  })()`))
  let target = null
  const data = arts.code === 0 ? arts.data : null
  if (data && Array.isArray(data.list)) {
    target = data.list.find((a) => (a.title || '').includes('从零实现')) || data.list.find((a) => a.status === 1)
    console.log('文章总数', data.list.length, '| 目标:', target ? `${target.id}「${target.title}」status=${target.status}` : '未找到')
  } else {
    console.log('列表接口异常', JSON.stringify(arts).slice(0, 300))
  }

  // 1) 表格操作列可达性（mobile 390）
  await cdp.send('Emulation.setDeviceMetricsOverride', { width: 390, height: 844, deviceScaleFactor: 1, mobile: true })
  const evStart = cdp.events.length
  await cdp.send('Page.navigate', { url: BASE + '/admin/articles' }); await sleep(3000)
  const tableInfo = await evalJs(cdp, `(() => {
    const wrap = document.querySelector('.el-table__body-wrapper .el-scrollbar__wrap, .el-table__body-wrapper')
    const ops = [...document.querySelectorAll('.el-table button, .el-table a')].filter(b => /编辑|删除/.test(b.textContent))
    const rect = ops[0] ? ops[0].getBoundingClientRect() : null
    const info = {
      wrapClient: wrap ? wrap.clientWidth : null,
      wrapScroll: wrap ? wrap.scrollWidth : null,
      opBefore: rect ? { left: Math.round(rect.left), right: Math.round(rect.right) } : null,
      vw: document.documentElement.clientWidth
    }
    return JSON.stringify(info)
  })()`)
  console.log('[表体内滚-初始]', tableInfo)
  const after = await evalJs(cdp, `(() => {
    const card = document.querySelector('.admin-card')
    if (card) card.scrollLeft = 99999
    return card ? { sl: card.scrollLeft, sw: card.scrollWidth } : null
  })()`)
  await sleep(400)
  const tableInfo2 = await evalJs(cdp, `(() => {
    const ops = [...document.querySelectorAll('.el-table button, .el-table a')].filter(b => /编辑|删除/.test(b.textContent))
    const rect = ops[0] ? ops[0].getBoundingClientRect() : null
    return JSON.stringify({ scrolled: ${JSON.stringify(after)}, opAfter: rect ? { left: Math.round(rect.left), right: Math.round(rect.right) } : null, vw: document.documentElement.clientWidth })
  })()`)
  console.log('[表体内滚-滑到底]', tableInfo2)
  await cap(cdp, OUT + '/articles_table_scrolled.png')
  console.log('[表格页console错误]', consoleLines(cdp.events, evStart).slice(0, 10))

  if (target) {
    // 2) 编辑页报错复现
    const ev2 = cdp.events.length
    await cdp.send('Page.navigate', { url: `${BASE}/admin/articles/${target.id}/edit` }); await sleep(3500)
    const editPage = await evalJs(cdp, `({ t: document.title, bodyLen: document.body.innerText.length,
      bodyHead: document.body.innerText.slice(0, 120) })`)
    console.log('[编辑页]', JSON.stringify(editPage))
    const errs = consoleLines(cdp.events, ev2)
    console.log('[编辑页console错误 数量]', errs.length)
    errs.slice(0, 15).forEach((e) => console.log('   ' + e))
    await cap(cdp, OUT + '/edit_page.png')

    // 3) 前台详情页打开草稿
    const ev3 = cdp.events.length
    await cdp.send('Page.navigate', { url: `${BASE}/articles/${target.id}` }); await sleep(3500)
    const detail = await evalJs(cdp, `({ t: document.title, bodyHead: document.body.innerText.slice(0, 150) })`)
    console.log('[前台详情]', JSON.stringify(detail))
    const errs3 = consoleLines(cdp.events, ev3)
    console.log('[前台详情console错误 数量]', errs3.length)
    errs3.slice(0, 15).forEach((e) => console.log('   ' + e))
    await cap(cdp, OUT + '/detail_page.png')
  }
  cdp.close()
  console.log('DONE')
  process.exit(0)
}
main().catch((e) => { console.error('FATAL', e); process.exit(1) })
