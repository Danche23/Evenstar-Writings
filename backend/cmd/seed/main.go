package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/pkg/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/redis/go-redis/v9"
)

// md 将多行拼接为正文
func md(parts ...string) string { return strings.Join(parts, "\n") }

type seedArticle struct {
	title, summary string
	body           []string
	cats, tags     []string
	views          uint
	date           time.Time
	status         int8
}

func main() {
	cfg := config.MustLoad("configs/config.yaml")
	m := cfg.Database.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal("connect:", err)
	}

	// ========== 1. 清空业务表并重置自增 ==========
	for _, t := range []string{"comments", "article_tags", "article_categories", "articles", "tags", "categories", "messages", "uploads", "users"} {
		if t == "comments" {
			// comments 有自引用外键(parent_id)，先删子行再删整表
			_ = db.Exec("DELETE FROM comments WHERE parent_id IS NOT NULL").Error
		}
		if err := db.Exec("DELETE FROM `" + t + "`").Error; err != nil {
			log.Fatalf("clean %s: %v", t, err)
		}
		_ = db.Exec("ALTER TABLE `" + t + "` AUTO_INCREMENT = 1")
	}

	// ========== 2. 用户 ==========
	hash := func(pw string) string {
		b, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		return string(b)
	}
	adminPw := hash("MuXing@2026") // 管理员初始密码，请在后台修改
	admin := model.User{
		Username: "evenstar", Nickname: "暮星", Email: "evenstar@evenstar.local",
		Password: adminPw, Role: 1, Status: 1,
		Bio: "暮星，意为把白昼散落出去的一切带回家。写代码，也写随笔。",
	}
	db.Create(&admin)

	type gu struct{ username, nickname, email, bio string }
	users := []gu{
		{"linmu", "林暮", "linmu@evenstar.local", "前端工程师，喜欢把日子过得像代码一样整洁。"},
		{"yaya", "芽芽", "yaya@evenstar.local", "大四学生，正在学习 Go，养了一盆文竹。"},
		{"qingshan", "青山", "qingshan@evenstar.local", "工作五年的后端，夜里读书白天搬砖。"},
		{"moli", "茉莉", "moli@evenstar.local", "设计师，相信好看的东西值得慢慢来。"},
		{"yefeng", "夜风", "yefeng@evenstar.local", "折腾过很多技术，最后发现跑步最治失眠。"},
		{"xiaobei", "小北", "xiaobei@evenstar.local", "在南方的小城写文档，周末去旧书店。"},
		{"chenxing", "晨星", "chenxing@evenstar.local", "数据相关的从业者，业余拍照。"},
		{"nanfeng", "南风", "nanfeng@evenstar.local", "考研党，间歇性写博客，持续性地想躺平。"},
		{"wanqing", "晚晴", "wanqing@evenstar.local", "爱做饭的测试工程师，评论永远带表情。"},
		{"shuimu", "水木", "shuimu@evenstar.local", "运维一枚，热衷把流程画成好看的图。"},
		{"anran", "安然", "anran@evenstar.local", "刚入行的新人，觉得写笔记是快乐的复利。"},
	}
	userIDs := []uint{admin.ID}
	upw := hash("1234abcd") // 普通用户统一测试密码
	for _, u := range users {
		row := model.User{Username: u.username, Nickname: u.nickname, Email: u.email,
			Password: upw, Role: 2, Status: 1, Bio: u.bio}
		db.Create(&row)
		userIDs = append(userIDs, row.ID)
	}
	nickOf := map[uint]string{admin.ID: "暮星"}
	for _, u := range users {
		var r model.User
		db.Select("id", "nickname").Where("username = ?", u.username).First(&r)
		nickOf[r.ID] = u.nickname
	}

	// ========== 3. 分类 / 标签 ==========
	catNames := []string{"编程技术", "前端开发", "数据库与缓存", "部署与上线", "读书笔记", "生活随笔", "工具与效率"}
	catID := map[string]uint{}
	for i, n := range catNames {
		c := model.Category{Name: n, SortOrder: i + 1}
		db.Create(&c)
		catID[n] = c.ID
	}
	tagNames := []string{"Go", "Gin", "GORM", "Redis", "MySQL", "Vue3", "TypeScript", "并发", "Docker", "Linux", "性能优化", "单元测试", "读书", "效率工具", "摄影", "随笔"}
	tagID := map[string]uint{}
	for i, n := range tagNames {
		t := model.Tag{Name: n, SortOrder: i + 1}
		db.Create(&t)
		tagID[n] = t.ID
	}

	// ========== 4. 文章（技术 7 篇 + 生活 7 篇） ==========
	D := func(day int) time.Time { return time.Date(2026, 7, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, day) }
	articles := []seedArticle{
		{"用 Go 泛型写一个可复用的分页查询", "把日常里重复的『查列表 + 分页』抽成一个泛型函数，是我第一次觉得泛型真的能减少代码。", []string{
			"三个月的 Go 项目做下来，我发现自己写得最多的代码不是业务逻辑，而是各种列表查询：`ListByXxx(page, size)`、`CountByXxx()`……它们长得几乎一样，却因为实体不同反复复制粘贴。",
			"## 先看看重复的部分", "以 GORM 为例，大部分列表查询都是固定三板斧：组装查询条件、统计总数、分页取数据。区别只在『条件怎么拼』和『查到什么实体』。", "于是我把它们抽成一个泛型函数，签名长这样：", "```go\nfunc Paginate[T any](q *gorm.DB, page, size int) ([]T, int64, error)\n```", "内部先 `Count` 拿总数，再 `Offset().Limit().Find(&items)`，调用方只需传入带好条件的 `q`。", "## 泛型在这里解决了什么", "没有泛型时，要么写 N 份几乎相同的函数，要么用 `interface{}` 丢掉类型。泛型让『同一份分页逻辑、每个调用点都有自己的类型』成为可能。", "- 好处：样板代码明显变少，逻辑统一收口", "- 代价：报错信息偶尔更绕，初学者理解成本略高", "## 一点体会", "> 泛型不是银弹，但当重复的模式出现第三次时，它值得被考虑。", "把它用在自己项目的列表接口后，新增一个模块只需要写『条件』，分页骨架再也不用碰了。", "```go\nitems, total, err := Paginate[model.Article](db.Where(\"status = ?\", 2), page, size)\n```", "这就是编程里很解压的一刻：删掉重复，留下意图。"},
			[]string{"编程技术", "数据库与缓存"}, []string{"Go", "GORM", "性能优化"}, 866, D(58), 2},
		{"Redis 缓存三问：穿透、击穿、雪崩的『简单够用』解法", "面试高频的三兄弟，被讲得天花乱坠。回到小站场景，哪种防护才值得写进代码？", []string{
			"聊到 Redis 缓存，穿透、击穿、雪崩这三个词几乎是必考题。面试时可以讲得很深，但回到自己维护的小博客，很多方案其实是过度设计。",
			"## 缓存穿透：查一个不存在的东西", "请求一个数据库里也没有的 key，缓存永远miss，压力全部打给 DB。最简单的办法是『空值也缓存一分钟』，再加一层布隆过滤器属于锦上添花。",
			"```go\n// 命中空缓存也返回，避免穿透\nif v == \"\" {\n\treturn emptyResult, nil\n}\n```",
			"## 缓存击穿：热点 key 过期的一瞬间", "同一时刻大量请求打到刚好过期的热门 key。小方案：加互斥锁，只放一个请求去回源，其余等待。",
			"## 缓存雪崩：一大批 key 同时失效", "给过期时间加随机扰动（基础时间 + rand 1~5 分钟），让失效时间散开，比上 Sentinel 之类的组件实在得多。",
			"## 选择判断", "- 个人站/小项目：空值缓存 + 随机过期，十几行代码就够", "- 大流量系统：才需要布隆过滤器、多级缓存、限流降级那一整套", "> 先确认问题真的存在，再决定要不要背上复杂度。", "这也是我给自己项目写防刷逻辑时的原则：满足功能，不堆砌。"},
			[]string{"数据库与缓存"}, []string{"Redis", "性能优化", "并发"}, 912, D(62), 2},
		{"Gin 中间件手记：从日志到鉴权的三层组织", "中间件是最容易写乱的部分。这篇记录我整理后的分层思路：全局、路由组、单路由。", []string{
			"Gin 的中间件看起来很简单：一个 `c.Next()` 前后各做各的事。但项目变大后，中间件放在哪一层、按什么顺序执行，很快会变成一团乱麻。",
			"## 三层组织法", "- 全局层：`engine.Use()`，放 Recovery、请求日志、CORS", "- 路由组层：比如 `/api/admin` 统一挂 `Auth + AdminOnly`", "- 单路由层：个别接口再挂限流之类的细粒度中间件", "```go\nadmin := r.Group(\"/api/admin\")\nadmin.Use(middleware.Auth(), middleware.AdminOnly())\n{\n\tadmin.GET(\"/users\", h.AdminListUsers)\n}\n```",
			"## 关于执行顺序", "中间件按挂载顺序先进先执行；`c.Next()` 之后的部分是『后置处理』，适合写耗时统计、统一错误兜底。", "把 `defer` 和 `c.Next()` 配合好，可以做到：进来先记开始时间，出去再算耗时写日志。",
			"## 容易踩的坑", "1. 在中间件里直接 `c.Abort()` 却忘了写响应，客户端会一直转圈", "2. 把业务状态塞进 `c.Set()` 的 key 起得随意，后面取错", "3. 鉴权中间件里查了两次用户：一次验证、一次用——应该只查一次存进上下文",
			"> 中间件是横切面，不是业务兜底。", "整理成三层之后，新增一个受保护接口的成本，从『到处复制鉴权代码』变成了『放进正确的组』。"},
			[]string{"编程技术"}, []string{"Go", "Gin", "性能优化"}, 743, D(53), 2},
		{"MySQL 索引失效的五个日常场景", "明明是普通查询，怎么还是全表扫描？把最容易踩的五个坑记下来。", []string{
			"优化慢查询，第一步永远是看执行计划。`EXPLAIN` 里出现 `type: ALL` 时，多半是索引没有按预期生效。",
			"## 场景一：对索引列用了函数", "```sql\nWHERE DATE(created_at) = '2026-07-01'\n```", "函数让索引失去意义，应改成范围查询：`created_at >= ? AND created_at < ?`。",
			"## 场景二：隐式类型转换", "字符串列用数字去查（或反过来），MySQL 可能放弃索引。",
			"## 场景三：左模糊", "`LIKE '%keyword'` 无法用前缀索引；搜索引擎的需求请交给专门的方案，别硬拗 SQL。",
			"## 场景四：联合索引没走最左前缀", "`(a, b, c)` 的联合索引，跳过 `a` 直接用 `b` 过滤，索引基本白建。",
			"## 场景五：OR 连接了非索引列", "`a = 1 OR other = 2`，其中一列无索引，优化器常直接选全表扫描；改成 `UNION` 或两段查询。",
			"## 排查顺序", "- 先 `EXPLAIN` 确认访问类型", "- 再确认索引列有没有被函数/转换包裹", "- 最后检查查询顺序是否满足最左前缀", "> 索引优化不是背口诀，而是看执行计划说话。", "数据库的『面试题』大多这样：原理简单，验证起来却要认真。"},
			[]string{"数据库与缓存"}, []string{"MySQL", "性能优化"}, 655, D(47), 2},
		{"Vue3 组合式开发半年后的十个体会", "从 Options API 迁到 setup 半年，记录那些『早知道就好了』的体会。", []string{
			"半年前把博客后台从零搭起时直接选了组合式 API。半年下来，有些体会值得写下来，给正在入门的你参考。",
			"## 体会一：按『逻辑』组织，而不是按『选项』", "Options API 把同一个功能的数据、方法拆到不同区块；组合式让我们把一段功能聚在一起，可读性对中等复杂页面提升明显。",
			"## 体会二：ref 与 reactive 别混用成习惯", "统一用 `ref` 最简单：取值的 `ref.value` 写多了就习惯了，模板里会自动解包。",
			"## 体会三：watch 别写太多", "能通过计算属性解决的问题，不要用 watch 去『同步状态』，后者容易造成隐式循环。",
			"## 体会四：组件别太胖", "超过两百行的单文件组件，先想想能不能拆出子组件或 composable。",
			"```js\n// 一个把列表页公共逻辑抽出来的 composable 雏形\nfunction usePagedList(fetcher) {\n\tconst list = ref([])\n\tconst page = ref(1)\n\tconst load = async () => { list.value = (await fetcher(page.value)) || [] }\n\treturn { list, page, load }\n}\n```",
			"## 体会五到十（简版）", "- 5. `v-model` 的组件封装要早点抽象", "- 6. 路由懒加载：后台页面全量打包有 700KB，拆开立竿见影", "- 7. 别在模板里写复杂表达式，难看且难测", "- 8. 给状态管理留白：能不用 Pinia 就不上", "- 9. 深链参数从 `route.query` 读，别存进 store", "- 10. 先跑起来，再谈规范", "> 技术选型没有对错，只有适不适合当下的自己。", "半年的体会就这些，欢迎在评论区交流你的版本。"},
			[]string{"前端开发"}, []string{"Vue3", "TypeScript", "性能优化"}, 588, D(40), 2},
		{"给后端小白的 Docker 部署笔记：一个 Go 服务的一生", "从『在我电脑上是好的』到『容器里也能跑』，记录部署一个 Go 博客的完整过程。", []string{
			"部署是每个后端都绕不开的一课。这篇以博客 API 为例，记录从 Dockerfile 到 docker-compose 的完整链路。",
			"## 第一步：多阶段构建镜像", "Go 是静态编译语言，非常适合做成轻量镜像：", "```dockerfile\nFROM golang:1.23 AS build\nWORKDIR /app\nCOPY go.mod go.sum ./\nRUN go mod download\nCOPY . .\nRUN CGO_ENABLED=0 go build -o server ./cmd/server\n\nFROM alpine:3.20\nCOPY --from=build /app/server /server\nEXPOSE 8080\nCMD [\"/server\"]\n```",
			"## 第二步：把配置和程序分开", "镜像里不放任何真实配置，用环境变量注入；本地开发用 `config.yaml`，容器里用 `APP_ENV=prod` 覆盖。",
			"## 第三步：compose 一把梭", "MySQL + Redis + API 三个服务用 docker-compose 编排，数据库数据挂卷，避免容器一删全没。",
			"## 常见翻车现场", "- 容器里连不上宿主 MySQL：`host.docker.internal` 或干脆把 DB 也放进 compose 网络", "- 时区不对：镜像里补 `TZ=Asia/Shanghai`", "- 健康检查没配，负载均衡把没就绪的容器也派发流量",
			"## 收尾检查清单", "1. 构建产物不包含源码里的密钥", "2. 数据卷已经挂载", "3. 日志有输出途径（stdout / 文件卷）", "4. 回滚方式想好了：镜像带版本号", "> 部署的核心不是把服务跑起来，而是让它能被安全地重启。", "祝你的服务也能『一次构建，到处运行』。"},
			[]string{"部署与上线"}, []string{"Docker", "Linux", "Go"}, 477, D(33), 2},
		{"单元测试从零开始：给函数安全感", "测试不是 KPI，是给未来改代码的自己一份安心。", []string{
			"写测试这件事，起初我是抗拒的：感觉代码写两遍，产出慢一倍。直到一次重构差点把登录逻辑改坏，才开始认真对待。",
			"## 从纯函数测起", "最容易上手的测试对象是『没有外部依赖』的纯函数：参数进，结果出。比如格式化时间、解析 token。",
			"```go\nfunc TestFormatDate(t *testing.T) {\n\tcases := []struct{ in, want string }{\n\t\t{\"2026-07-01T10:00:00+08:00\", \"2026-07-01\"},\n\t}\n\tfor _, c := range cases {\n\t\tif got := FormatDate(c.in); got != c.want {\n\t\t\tt.Errorf(\"FormatDate(%q) = %q, want %q\", c.in, got, c.want)\n\t\t}\n\t}\n}\n```",
			"## 表驱动测试是 Go 的温柔", "一个 `cases` 数组走完所有分支，新增用例只是加一行，读起来像文档。",
			"## 涉及数据库怎么办", "真实依赖太重，就抽象接口、注入替身；项目里我习惯把 repo 层做成可替换的，测试时给个内存版。",
			"## 别追求覆盖率数字", "覆盖率 60% 覆盖到核心逻辑，远好过 90% 都是无效断言。先把最重要的登录、支付、权限这些路径测住。",
			"> 测试的价值，在于你深夜改代码时敢不敢按回车。", "从那以后我养成了习惯：新函数写完，先花十分钟给它一点安全感。"},
			[]string{"编程技术"}, []string{"Go", "单元测试"}, 391, D(26), 2},
		{"当我谈读书时，我谈些什么——七月书单", "村上春树式的标题，七月读的三本书：一本谈写作，一本谈技术史，一本是诗集。", []string{
			"七月读得不算多，三本。但每一本都留下了值得写下来的东西。",
			"## 一、《写作这回事》——斯蒂芬·金", "这不是教你成为作家，而是告诉你写作是一件需要『每天开工』的手艺。他反复强调的一件事是：删掉形容词，让句子自己走路。",
			"> 我写作的秘诀很简单：关上门写，打开门改。",
			"## 二、《代码之外》——几位老程序员的访谈录", "技术书读多了会以为世界是逻辑的。这本书里全是反例：有人在车库写出了改变行业的产品，也有人在国企维护二十年的 COBOL。技术史的一半是人的故事。",
			"## 三、《我喜爱一切不彻底的事物》——张定浩", "诗不负责解决问题，它负责把问题摆得好看一点。这本书适合睡前读两页，梦里都是桂花的味道。",
			"## 一点摘录式的感想", "读书最大的好处，是让孤独变成了可以交流的东西——你不需要认识作者，却能在他人的句子里认出自己。",
			"立个 flag：八月读三本技术 + 一本小说，让输入和输出保持平衡。"},
			[]string{"读书笔记"}, []string{"读书", "随笔"}, 534, D(16), 2},
		{"夏末的夜跑，与一颗叫暮星的星", "晚风、路灯和跑步时想明白的小事。", []string{
			"九月的晚风开始有凉意了。八点半出门夜跑，路灯把影子拉得很长，蝉鸣比七月收敛了许多。",
			"## 跑到第三圈，想明白的事", "白天写代码的时候，大脑被各种 deadline 占满；只有跑步时，思绪才会自己浮上来。今晚想明白的是：不必把每件事都做到最好，做完本身就很了不起。",
			"## 暮色与暮星", "跑完坐在河边的长椅上拉伸，天边还剩最后一抹蓝。第一颗星亮起来的时候，我突然想起萨福写暮星的那句诗——它带回白昼散出去的一切。", "> 白昼把我们散出去，入夜，总该有谁把我们带回来。",
			"## 关于跑步的碎碎念", "- 三公里是底线，五公里是舒适区", "- 跑步不减肥，但跑完那一小时你会原谅全世界", "- 夜跑记得穿亮色，比配速重要",
			"如果你也夜跑，欢迎留言分享你的路线和心事。"},
			[]string{"生活随笔"}, []string{"随笔", "摄影"}, 468, D(11), 2},
		{"旧书店半日：在南方小城的下午捡到宝", "纸质书的香气、老板娘养的猫，和一个下午的漫游。", []string{
			"周末本来只想去超市，路过街角那家开了十几年的旧书店，还是没忍住推门进去。",
			"## 旧书店的规矩", "老板娘在柜台后织毛衣，头也不抬地说：自己找，别弄乱。她养的那只橘猫正睡在《辞海》上，对每个进店的人都不屑一顾。",
			"书架按自己的逻辑排着：门口是小说，往里走是历史，楼梯下藏着整套《古拉格群岛》，角落有一摞九十年代的《计算机世界》——铅字年代的编程杂志，读起来像考古。",
			"## 战利品", "最后带走三本：一本一九九七年印次的《苏菲的世界》，一本缺了封面的博尔赫斯选集，还有一本不知谁在扉页写下『读于北京到上海的火车上』的散文集。",
			"> 旧书最好的部分，是上一任主人的痕迹。",
			"结账时老板娘终于抬头：一共十八块。猫醒了，伸了个懒腰，又睡回去。",
			"这大概就是小城的浪漫：时间很慢，书很便宜，而黄昏很快就来。"},
			[]string{"生活随笔"}, []string{"读书", "随笔"}, 305, D(7), 2},
		{"桌面整理术：给工作区一个『家』", "数字和物理的桌面一起整理，效率会以想象不到的方式回来。", []string{
			"这篇想写的不是鸡汤，而是一套可以照着做的整理流程。起因是某天找一份合同花了二十分钟，终于决定动手。",
			"## 物理桌面：只留三种东西", "1. 正在用的：键盘、鼠标、笔记本", "2. 每天要看的：一杯水、一盆植物", "3. 一件让你高兴的小摆件——我放的是在旧书店淘的月球灯", "其余全部进抽屉或垃圾桶。",
			"## 数字桌面：三文件夹原则", "```\n~/inbox    # 一切新东西先进来，每周清空一次\n~/archive  # 已归档的项目与资料\n~/tmp      # 三天内确定要删的\n```",
			"规则很简单：新文件一律先进 `inbox`；每周五下午用半小时清 inbox，要么归类，要么进 tmp。",
			"## 效果", "整理完的第一个感受不是『变干净了』，而是『找东西不用再想』。大脑省下来的那点注意力，足以支撑一整个下午的专注。",
			"> 秩序不是束缚，是把注意力还给重要的事。", "如果你也有一个乱糟糟的桌面，不妨从清空桌面图标开始，五分钟就够。"},
			[]string{"工具与效率"}, []string{"效率工具", "随笔"}, 412, D(20), 2},
		{"写作习惯养成记：从月更到周更的秘密", "不靠意志力，靠一套低门槛的流程。", []string{
			"三年前我一年写不出一篇文章，总觉得『没准备好』。现在我保持周更，靠的不是更自律，而是把写作的门槛拆到几乎为零。",
			"## 秘密一：先建『素材箱』", "平时想到的任何题目、金句、截图，三秒钟内丢进素材箱，绝不在当时展开。素材够了，文章自然有得写。",
			"## 秘密二：写草稿不写文章", "第一稿只允许自己写『烂句子』，目标不是质量，而是把骨架搭出来。写完搁一天，第二天再改，质量会翻倍。",
			"## 秘密三：固定仪式感", "每周日晚九点，热水、台灯、同一支笔（或同一个编辑器主题）。身体记住这个仪式后，进入状态会越来越快。",
			"## 关于平台与反馈", "发表在博客上有一个隐藏好处：评论区的每一条认真回复，都是一次小的正反馈，支撑你写下一篇。",
			"> 写作是思考的显影液。", "把写作当成健身：不追求一次写完美，追求长期在场。八月，要不要一起开始？"},
			[]string{"生活随笔", "工具与效率"}, []string{"效率工具", "随笔"}, 356, D(4), 2},
		{"给一年后的自己写一封信", "没有主题的随笔，写给未来某个周二晚上的自己。", []string{
			"今天是九月第一个周日，窗外在下雨。我想趁天气刚好，给一年后的自己写封信，不讲道理，只讲近况。",
			"## 近况", "博客终于像样了：有了稳定的文章、认真回复的读者，以及一颗被无数次提到却依然喜欢的星。工作上还在和 Redis、并发和无穷无尽的 bug 搏斗，偶尔赢。",
			"## 几个问题想问问一年后的你", "- 还在坚持周更吗，还是又被忙碌偷走了？", "- 那篇写了三个月的『GORM 深入浅出』，发布了吗？", "- 阳台的文竹，还活着吗？",
			"## 一些不改的约定", "1. 不追热点，不写自己不信的东西", "2. 每篇至少有一处对读者有用的细节", "3. 遇到好文章就留言，把感谢说出来", "> 愿你还是那个会在黄昏看星星的人。",
			"这封信没有结尾，因为故事还在写。如果你读到这篇，也算和一年后的我打了个照面——那时候记得提醒我回信。"},
			[]string{"生活随笔"}, []string{"随笔"}, 302, D(1), 2},
		{"（草稿）从零实现一个极简任务队列", "还没写完的草稿，先占个位——想聊聊 goroutine、channel 与背压。", []string{
			"## 想写的内容", "1. 为什么需要任务队列：削峰与解耦", "2. 用 channel 实现 worker pool 的最小版本", "3. 任务失败重试与优雅退出", "4. 背压：队列满了怎么办",
			"## 目前的半成品代码", "```go\nfunc Worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {\n\tdefer wg.Done()\n\tfor j := range jobs {\n\t\tprocess(j)\n\t}\n}\n```",
			"## 待补充", "- 并发安全与指标收集", "- 与 Redis 队列的取舍对比", "- 实验数据：不同 worker 数下的吞吐",
			"> 未完待续……",
			"如果你也写过类似的东西，欢迎先来评论区聊聊你的方案。"},
			[]string{"编程技术"}, []string{"Go", "并发"}, 0, D(68), 1},
	}

	artID := []uint{}
	for _, a := range articles {
		pub := a.date
		art := model.Article{AuthorID: admin.ID, Title: a.title, Summary: a.summary,
			Content: md(a.body...), Status: a.status, Views: a.views, PublishedAt: &pub}
		db.Create(&art)
		artID = append(artID, art.ID)
		for _, cn := range a.cats {
			db.Create(&model.ArticleCategory{ArticleID: art.ID, CategoryID: catID[cn]})
		}
		for _, tn := range a.tags {
			db.Create(&model.ArticleTag{ArticleID: art.ID, TagID: tagID[tn]})
		}
	}
	fmt.Printf("seeded %d articles\n", len(artID))

	// ========== 5. 评论（含置顶） ==========
	// 顺序：先一级（记录 id），后二级
	type pair struct{ a, u int }
	topSpecs := []struct {
		art, user int
		top       bool
		content   string
	}{
		{0, 1, true, "这篇写得真好，泛型那段豁然开朗，收藏了！"},
		{0, 5, false, "博主用的 GORM 版本是多少？Paginate 的泛型约束具体怎么写的？"},
		{0, 7, false, "同一个模式我项目里也出现过，明天就重构试试。"},
		{1, 4, false, "面试被问雪崩，终于看到能落地的版本了。随机过期时间这个思路赞。"},
		{1, 9, true, "作为小站维护者深有体会：很多高级方案真的用不上，够用就好 😄"},
		{2, 2, false, "中间件执行顺序那块终于理清楚了，感谢！"},
		{2, 3, false, "c.Abort() 忘记写响应这个坑我踩过，差点排查一晚上。"},
		{3, 6, false, "EXPLAIN 真是好老师，第五个场景让我想起上周的慢查询。"},
		{4, 2, false, "体会四深有同感，组件超 200 行就该拆了。"},
		{7, 9, false, "张定浩那本我也在读！桂花那段写得好。"},
		{8, 10, false, "夜跑党路过，三公里底线 +1。暮星那段看哭了。"},
		{9, 5, false, "旧书店坐标可以透露一下吗？想去淘书。"},
		{11, 8, false, "素材箱的方法好用，已经坚持两周了。"},
		{12, 3, false, "读到这封信感觉被温柔地叫醒了一下，谢谢博主。"},
	}
	topIDs := []uint{}
	topArt := []uint{} // topIDs[i] 所属文章
	for _, s := range topSpecs {
		uid := userIDs[s.user]
		cm := model.Comment{ArticleID: artID[s.art], UserID: &uid, Content: s.content}
		if s.top {
			now := time.Now().Add(-24 * time.Hour)
			cm.IsTop, cm.TopTime = 1, &now
		}
		db.Create(&cm)
		topIDs = append(topIDs, cm.ID)
		topArt = append(topArt, artID[s.art])
	}
	replySpecs := []struct {
		top, u  int
		content string
	}{
		{0, 2, "同感，特别是把条件封装进去之后，新增模块真的只剩写 where 了。"},
		{1, 1, "谢谢喜欢！小站规模下确实够用就好，等真到了需要布隆过滤器那天再说。"},
		{1, 4, "哈哈博主求生欲拉满，先立免责声明。"},
		{2, 6, "第三个坑太真实了，Set 的 key 建议统一放个常量文件。"},
		{2, 1, "@青山 对，我现在都是定义 CtxKey 常量，避免拼错。"},
		{5, 8, "同求 repo 内存版写法，能展开讲讲吗？"},
		{7, 1, "@茉莉 桂花那段我也最舍不得删哈哈。"},
		{10, 1, "@芽芽 收到，素材箱也是评论区给我的灵感！"},
		{11, 10, "等回复！我也想去，最好是有猫的那家。"},
	}
	for _, s := range replySpecs {
		uid := userIDs[s.u]
		cm := model.Comment{ArticleID: topArt[s.top], UserID: &uid,
			ParentID: &topIDs[s.top], ReplyToID: &topIDs[s.top], Content: s.content}
		db.Create(&cm)
	}
	fmt.Println("seeded comments:", len(topSpecs)+len(replySpecs))

	// ========== 6. 留言板 ==========
	msgs := []struct {
		u       int
		content string
	}{
		{1, "第一次来，被首页那句『白昼散落的句子，该回家了』击中了。"},
		{3, "关注很久了，从技术文章到生活随笔都爱看，加油！"},
		{4, "可以求一个 RSS 吗？想订阅～"},
		{5, "夜跑路线那篇让我今晚也出门跑了三公里。"},
		{7, "留言板这个功能本身就很浪漫。"},
		{8, "刚注册的新人，报到！顺便表白那盆文竹哈哈。"},
		{9, "水木说的对，留言板真浪漫。祝博主每天都有一颗很亮的星。"},
	}
	for _, s := range msgs {
		uid := userIDs[s.u]
		db.Create(&model.Message{UserID: &uid, Content: s.content})
	}
	fmt.Println("seeded messages:", len(msgs))

	// ========== 7. Redis 热门榜回灌 ==========
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.DB,
	})
	defer rdb.Close()
	_ = rdb.Del(ctx, "hot:articles").Err()
	for i, a := range articles {
		if a.status != 2 {
			continue
		}
		_ = rdb.ZAdd(ctx, "hot:articles", redis.Z{Score: float64(a.views), Member: strconv.FormatUint(uint64(artID[i]), 10)}).Err()
	}
	fmt.Println("redis hot list refreshed")

	fmt.Println("=== SEED DONE ===")
	fmt.Println("管理员: evenstar@evenstar.local / MuXing@2026 （昵称 暮星）")
	fmt.Println("普通用户统一密码: 1234abcd")
}
