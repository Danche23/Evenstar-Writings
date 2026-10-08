# 贡献指南 / Contributing

感谢你有兴趣参与 Evenstar Writings！本文说明如何搭建开发环境、代码规范与提交流程。

## 开发环境

请先按 [README 的「快速开始」](README.md#快速开始) 跑通前后端。

| 依赖 | 版本 |
|---|---|
| Go | 1.24+ |
| Node.js | 20+ |
| MySQL | 8.0 |
| Redis | 7 |

## 分支约定

- `main` 是唯一稳定的开发主线，所有改动通过 Pull Request 合入
- 功能分支建议按类型命名：
  - `feat/xxx` 新功能
  - `fix/xxx` 缺陷修复
  - `docs/xxx` 文档
  - `refactor/xxx` 重构
  - `chore/xxx` 构建 / 依赖 / 杂项

## 提交信息

使用语义化的提交信息（Conventional Commits 风格），主题行用简洁的祈使句：

```
feat: 文章支持按标签筛选
fix: 修复登录失败 3 次后滑块不触发的问题
docs: 补充 OSS 配置说明
refactor: 抽出分页参数校验
chore(deps): 升级 gin 到 v1.11.0
```

一个提交只做一件事，避免把重构与功能改动混在一起。

## 代码规范

**后端（Go）**

- 遵循 `gofmt` 格式，提交前跑 `make fmt`
- 分层保持清晰：`api`（HTTP 层）→ `service`（业务）→ `repository`（数据访问），不要在 handler 里直接写 SQL
- 错误统一用 `pkg/errors` 中定义的错误码返回，不要 `panic`
- 新增业务逻辑请补充单元测试，运行 `make test`
- 建议提 PR 前跑 `make vet`（可选 `make lint`，需本地装 golangci-lint）

**前端（Vue）**

- 组件使用 `<script setup>` 组合式 API
- 接口调用统一放在 `src/api/`，不要在组件里直接写 axios
- 提交前确保 `npm run build` 能通过

## 安全红线

**绝对不要提交以下内容**，它们在 `.gitignore` 中已被忽略，请勿强行 `git add -f`：

- `backend/configs/config.yaml`（真实配置，含数据库密码）
- 根目录 `.env`（部署密钥）
- `ssl/` 证书目录（`.pem` / `.key`）
- 云服务 AccessKey、SMTP 密码、JWT 密钥
- `node_modules/`、`dist/`、编译产物、日志文件
- 任何含真实账号密码的本地调试脚本

需要新增配置项时，**只更新对应的 `*.example` 模板**，填中文占位符或 `${VAR}` 引用，不要把真实值写进去。
测试用例中出现的密码同样必须使用假值。

如果误提交了真实密钥，仅删除文件不够——密钥仍留在 Git 历史中，请**立即在云端轮换该密钥**。

## 提交 Pull Request

1. Fork 本仓库并从 `main` 切出功能分支
2. 完成改动，确保本地 `go vet ./... && go test ./...` 与 `npm run build` 均通过
3. 按 [PR 模板](.github/PULL_REQUEST_TEMPLATE.md) 填写说明，重点写清**动机**与**验证方式**
4. 创建 PR 并等待 CI 通过
5. 评审意见通过后由维护者合并

**涉及数据库结构变更**时，请同步更新 `backend/scripts/init.sql`；
**涉及接口变更**时，请同步更新 `docs/openapi.yaml`。

## 报告问题

- Bug 与功能建议请走 [Issue 模板](.github/ISSUE_TEMPLATE)
- 安全漏洞请勿开公开 Issue，见 [SECURITY.md](.github/SECURITY.md)
