# Evenstar Writings

[![CI](https://github.com/Danche23/Evenstar-Writings/actions/workflows/ci.yml/badge.svg)](https://github.com/Danche23/Evenstar-Writings/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](backend/go.mod)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs&logoColor=white)](frontend/package.json)

> 一个前后端分离的个人博客 / 写作平台，自带完整后台管理。
> A full-stack personal blog platform with a built-in admin dashboard.

**English summary** — Evenstar Writings is a self-hosted blogging system with a public-facing reading site and an admin console's unified in one codebase. The backend is a Go (Gin + GORM) REST API backed by MySQL and Redis; the frontend is a Vue 3 + Element Plus SPA. It ships with Markdown authoring, hierarchical comments, email-code registration and password reset, image uploads (Aliyun OSS or local disk), slider captcha, JWT auth, rate limiting, and a Docker Compose deployment. Licensed under MIT.

---

## 目录

- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [目录结构](#目录结构)
- [快速开始](#快速开始)
- [部署](#部署)
- [配置说明](#配置说明)
- [API 文档](#api-文档)
- [安全说明](#安全说明)
- [License](#license)

## 功能特性

**前台（阅读端）**

- 门户首页：最新文章、热门文章、分类与标签导航
- 文章列表：按分类 / 标签筛选 + 关键词搜索 + 分页
- 文章详情：Markdown 渲染、代码高亮、浏览量统计
- 评论系统：两级嵌套评论、回复关系展示、评论置顶、软删除
- 归档页、分类页、标签云
- 留言板（游客留言）
- 用户中心：注册、登录、找回密码、修改资料、上传头像

**后台（管理端，需 `role=1`）**

- 数据看板：文章数 / 评论数 / 用户数 / 总浏览量
- 文章管理：Markdown 编辑器、草稿与发布、分类与标签绑定
- 分类 / 标签管理：增删改 + 拖拽排序
- 评论管理：分页浏览、置顶、删除
- 用户管理：搜索、禁用 / 解禁、删除
- 文件管理：上传记录与删除
- 站点设置（关于页内容）

**通用能力**

- JWT 认证，`token_version` 机制使改密后旧 token 立即失效
- 基于 Redis 的接口限流
- 邮箱验证码注册 / 找回密码（阿里云邮件推送）
- 阿里云验证码 2.0 滑块验证（登录连续失败后触发）
- 图片上传：阿里云 OSS，未配置时自动回退本地磁盘存储
- 统一响应结构、统一错误码、请求日志自动脱敏（密码不入日志）
- 健康检查接口 `/api/health`（含 MySQL / Redis 连通性探针）

## 技术栈

| 层 | 技术 |
|---|---|
| 后端语言 | Go 1.24 |
| Web 框架 | Gin |
| ORM | GORM (MySQL 8.0) |
| 缓存 / 限流 | Redis 7 (`go-redis/v9`) |
| 配置 | Viper（YAML + 环境变量覆盖） |
| 日志 | Zap + lumberjack（滚动切割） |
| 鉴权 | `golang-jwt/jwt/v5` + bcrypt |
| 前端 | Vue 3（`<script setup>`）+ Vite 5 |
| 状态 / 路由 | Pinia 2 + Vue Router 4 |
| UI 组件 | Element Plus |
| 富文本 | markdown-it + highlight.js + md-editor-v3 |
| XSS 防护 | DOMPurify |
| 部署 | Docker + Docker Compose（Nginx 托管前端） |

## 目录结构

```
.
├── backend/                    # Go 后端
│   ├── cmd/
│   │   ├── server/             # 服务入口
│   │   └── seed/               # 种子数据工具
│   ├── configs/
│   │   └── config.example.yaml # 配置模板（复制为 config.yaml 使用）
│   ├── internal/
│   │   ├── api/                # HTTP 处理器与路由（按业务模块划分）
│   │   ├── dto/                # 请求 / 响应数据结构
│   │   ├── model/              # GORM 数据模型
│   │   ├── repository/         # 数据访问层
│   │   ├── service/            # 业务逻辑层
│   │   ├── middleware/         # 认证 / CORS / 限流 / 恢复 / 请求日志
│   │   └── app/                # 应用装配与生命周期
│   ├── pkg/                    # 可复用基础包（config/database/jwt/logger/errors/storage/utils）
│   ├── scripts/init.sql        # 数据库建表脚本
│   ├── Dockerfile
│   └── Makefile
├── frontend/                   # Vue 3 前端
│   ├── src/
│   │   ├── api/                # 接口封装
│   │   ├── components/         # 通用组件
│   │   ├── layouts/            # 前台 / 后台布局
│   │   ├── router/             # 路由与守卫
│   │   ├── stores/             # Pinia 状态
│   │   ├── utils/              # 工具函数
│   │   └── views/              # 页面（含 admin/ 后台页面）
│   ├── Dockerfile
│   ├── nginx.conf
│   └── nginx.docker.conf
├── docs/
│   ├── openapi.yaml            # OpenAPI 3.0 接口文档
│   ├── 编码前准备.md            # 设计手记：页面结构 / 测试方案 / 工程收尾
│   └── prototype/              # 早期静态 UI 原型（设计参考）
├── docker-compose.yml          # 全容器化编排（web/backend/mysql/redis）
├── .env.example                # 部署密钥模板
└── LICENSE
```

## 快速开始

### 环境要求

| 依赖 | 版本 |
|---|---|
| Go | 1.24+ |
| Node.js | 20+ |
| MySQL | 8.0 |
| Redis | 7 |

### 1. 准备数据库

创建数据库并导入建表脚本：

```bash
mysql -u root -p -e "CREATE DATABASE evenstar DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p evenstar < backend/scripts/init.sql
```

### 2. 启动后端

```bash
cd backend
cp configs/config.example.yaml configs/config.yaml   # 按需修改 MySQL / Redis 连接
make deps
make run                                             # 默认监听 :8080
```

首次启动时，若 `users` 表中不存在管理员，会按 `EVENSTAR_ADMIN_*` 环境变量创建管理员；
未设置时生成随机密码并打印到后端日志。

验证服务是否就绪：

```bash
curl http://localhost:8080/api/health
```

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev                                          # http://localhost:5173
```

开发环境下 Vite 会把 `/api` 代理到 `http://localhost:8080`，无需额外配置跨域。

### 常用命令

```bash
cd backend
make build          # 编译到 bin/
make test           # 运行单元测试
make test-coverage  # 生成覆盖率报告
make fmt vet lint   # 格式化 / 静态检查
```

## 部署

项目附带 `docker-compose.yml`，一条编排拉起 `web`（Nginx + 前端静态产物）、
`backend`、`mysql`、`redis` 四个服务，只有 `web` 对外暴露端口，其余留在内部网络。

> **注意**：后端镜像采用「零编译」策略——服务器上不执行 `go build`（为适配小内存机型），
> 因此**必须先在本机交叉编译出 Linux 二进制**再构建镜像。

```bash
# 1) 本机交叉编译后端二进制（产物必须是 backend/server，且与 Dockerfile 同目录）
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o server ./cmd/server

# 2) 服务器上准备密钥文件（.env 已被 .gitignore 忽略，务必不要提交）
cp .env.example .env && chmod 600 .env
#    编辑 .env，填入真实的数据库密码、JWT_SECRET、邮件 / OSS / 验证码凭据

# 3) 构建镜像并启动
docker compose build
docker compose up -d mysql redis                     # 先起数据库
docker compose up -d backend web                     # 再起应用

docker compose ps                                    # 查看状态
docker compose logs -f backend                       # 查看后端日志
```

`backend/configs/config.yaml` 同样含密钥且不入镜像，运行时以只读方式挂载进容器。

## 配置说明

配置有两个来源，优先级从低到高：

1. `backend/configs/config.yaml` —— YAML 配置文件（本地开发常用）
2. 环境变量 —— 仅当变量非空时覆盖 YAML 中的对应项（容器化部署常用）

模板文件：

- [`backend/configs/config.example.yaml`](backend/configs/config.example.yaml) —— 后端完整配置模板
- [`backend/.env.example`](backend/.env.example) —— 后端环境变量模板
- [`frontend/.env.example`](frontend/.env.example) —— 前端环境变量模板
- [`.env.example`](.env.example) —— Docker Compose 部署密钥模板

主要配置项：

| 变量 | 说明 |
|---|---|
| `MYSQL_HOST` / `MYSQL_PASSWORD` / `MYSQL_DATABASE` | MySQL 连接（容器内 host 用 compose 服务名 `mysql`） |
| `REDIS_HOST` / `REDIS_PASSWORD` | Redis 连接（容器内 host 用 `redis`） |
| `JWT_SECRET` | JWT 签名密钥，**必须改为足够长的随机串** |
| `APP_MODE` / `LOG_LEVEL` | 运行模式（`debug` / `release`）与日志级别 |
| `CORS_ALLOW_ORIGINS` | 跨域白名单，填前端站点地址 |
| `MAIL_MOCK` | 设为 `true` 时验证码只打印到日志，不真发邮件（开发用） |
| `MAIL_HOST` / `MAIL_USERNAME` / `MAIL_PASSWORD` | SMTP 邮件推送配置 |
| `OSS_ACCESS_KEY_ID` / `OSS_ACCESS_KEY_SECRET` | 阿里云 OSS；未配置时自动回退本地磁盘存储 |
| `CAPTCHA_*` / `VITE_CAPTCHA_*` | 阿里云验证码 2.0 滑块，前后端需填同一场景 |
| `EVENSTAR_ADMIN_USERNAME` / `_EMAIL` / `_PASSWORD` | 首次启动时引导管理员账号 |

完整的键名与取值规则见各模板文件内的注释。

## API 文档

- OpenAPI 3.0 规范：[`docs/openapi.yaml`](docs/openapi.yaml)，可直接导入 Apifox / Postman / Swagger UI
- 健康检查：`GET /api/health`
- 所有业务接口以 `/api` 为前缀，后台接口位于 `/api/admin/*`，需要管理员身份

## 安全说明

- 仓库中**不含任何真实密钥**，仅保留 `*.example` 模板；请勿将 `config.yaml`、`.env`、证书等提交到版本库
- 密码使用 bcrypt 哈希存储；请求日志会自动脱敏 `password` 类字段
- 改密后通过 `token_version` 使旧 token 立即失效
- Markdown 渲染经 DOMPurify 过滤，上传文件按 MIME 白名单校验
- 发现安全漏洞请通过 [SECURITY.md](.github/SECURITY.md) 中说明的渠道反馈，请勿直接开公开 Issue

## 参与贡献

欢迎提交 Issue 与 Pull Request，请先阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。

## License

本项目基于 [MIT License](LICENSE) 开源。
