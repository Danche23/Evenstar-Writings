// 暮星絮语：按时段问候 + 按日轮换的短句（内置原创短句，文艺轻量）

export function greetNow(date = new Date()) {
  const h = date.getHours()
  if (h >= 5 && h < 9) return '清晨好，愿晨光落在你正翻开的书页上。'
  if (h >= 9 && h < 12) return '上午好，忙碌之间，也记得抬头看一眼云。'
  if (h >= 12 && h < 14) return '午安。饭后小憩片刻，星子还在路上。'
  if (h >= 14 && h < 18) return '午后安。不必急，黄昏还远，余温尚好。'
  if (h >= 18 && h < 22) return '暮安。暮色四合，星子初亮——今晚也请慢一点。'
  return '夜深了。早点休息，梦里也有星光。'
}

const QUOTES = [
  '暮色是写给白天的一封慢信。',
  '白昼把羊群放牧向远方，暮星升起时，带它们回家。',
  '把心事说给花听，花会替你开成春天。',
  '星子不问路，只管在夜里亮着。',
  '日子是慢慢写的，急不来。',
  '风经过花时，把香也分给了路过的人。',
  '月亮圆缺，都是它自己的圆满。',
  '写下来的念头，才算是认真活过的证据。',
  '黄昏把天幕调暗，好让星光登场。',
  '一茶一坐，时间忽然就温柔了。',
  '花开不是为了结果，是为了好看地开一次。',
  '昨夜的风，今晨的花，都值得被记下。',
  '读书如观星，越看越觉得自己渺小而辽阔。',
  '把喜欢的事做到安静，就是很浪漫的坚持。',
  '晚风不催人，星子不等人——去看想看的吧。',
  '墙角的野花，从不抱怨没人给它鼓掌。',
  '生活的诗意，常藏在没被安排好的缝隙里。'
]

export function dailyQuote(date = new Date()) {
  const start = new Date(date.getFullYear(), 0, 0)
  const day = Math.floor((date - start) / 86400000)
  return QUOTES[day % QUOTES.length]
}

// 文章心情题字：按标题稳定派生一句（同一篇文章永远同一句）
export function epigraphFor(text = '') {
  let h = 0
  const src = String(text || '')
  for (let i = 0; i < src.length; i++) h = (h * 31 + src.charCodeAt(i)) >>> 0
  return QUOTES[h % QUOTES.length]
}
