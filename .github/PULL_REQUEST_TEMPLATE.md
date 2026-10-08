## 这个 PR 做了什么

<!-- 一句话说明改动目的；如有对应 Issue，请写 Closes #123 -->

## 改动类型

- [ ] Bug 修复
- [ ] 新功能
- [ ] 重构 / 优化（不改变外部行为）
- [ ] 文档
- [ ] 构建 / 部署 / CI
- [ ] 其他：

## 主要改动点

<!-- 列出关键文件与实现思路，便于 reviewer 快速进入上下文 -->

-
-

## 测试与验证

<!-- 说明你如何验证这次改动；贴出关键命令与输出 -->

- [ ] `cd backend && go vet ./... && go test ./...` 通过
- [ ] `cd frontend && npm run build` 通过
- [ ] 已在真实运行环境中手动验证

复现 / 验证步骤：

```
# 例如：
# go test ./internal/service/ -run TestXxx -v
```

## 自查清单

- [ ] 未提交任何真实密钥、密码、证书或本地配置（只保留 `*.example` 模板）
- [ ] 未提交 `node_modules`、`dist`、编译产物、日志等生成物
- [ ] 涉及数据库结构变更时，已同步更新 `backend/scripts/init.sql`
- [ ] 涉及接口变更时，已同步更新 `docs/openapi.yaml`
- [ ] 代码风格与现有代码保持一致
