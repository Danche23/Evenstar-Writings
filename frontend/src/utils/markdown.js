import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js/lib/common'
import DOMPurify from 'dompurify'
import 'highlight.js/styles/atom-one-dark.css'

// 统一 Markdown 渲染器：原文 → HTML → DOMPurify 消毒 → 上屏
const md = new MarkdownIt({
  html: true,          // 允许内嵌 HTML（依赖 DOMPurify 消毒防 XSS）
  linkify: true,       // 自动识别链接
  breaks: true,        // 换行即 <br>
  highlight(str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre><code class="hljs language-${lang}">` +
          hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
          '</code></pre>'
      } catch (e) { /* fallthrough */ }
    }
    return `<pre><code class="hljs">${md.utils.escapeHtml(str)}</code></pre>`
  }
})

export function renderMarkdown(source = '') {
  const html = md.render(source)
  return addHeadingIds(DOMPurify.sanitize(html))
}

// 从 markdown 原文提取 h2/h3 目录（与 addHeadingIds 使用同一序号规则，id 一一对应）
export function extractToc(source = '') {
  const toc = []
  let n = 0
  let inCode = false
  const lines = String(source || '').split('\n')
  for (const line of lines) {
    if (/^\s*(```|~~~)/.test(line)) { inCode = !inCode; continue }
    if (inCode) continue
    const m = line.match(/^(\#{2,3})\s+(.+)$/)
    if (m) {
      n++
      const level = m[1].length // 2 | 3
      const text = m[2].trim().replace(/[*_`~]/g, '').replace(/\[([^\]]+)\]\([^)]*\)/g, '$1').trim()
      if (text) toc.push({ level, id: 'sec-' + n, text })
    }
  }
  return toc
}

// 给渲染出的 h2/h3 依次补 id（与 extractToc 的 sec-N 对应）
let tocCounter = 0
function addHeadingIds(html) {
  tocCounter = 0
  return html.replace(/<h([23])(?=[ >])/g, (all, lv) => {
    tocCounter++
    return `<h${lv} id="sec-${tocCounter}"`
  })
}
