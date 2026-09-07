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
  return DOMPurify.sanitize(html)
}
