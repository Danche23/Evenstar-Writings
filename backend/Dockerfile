# ===== backend/Dockerfile =====
# Evenstar Writings 后端镜像（Go 多阶段构建）
# 构建：docker build -t evenstar-backend:1.0 .
# 运行：docker run -d --name evenstar-backend -p 8081:8080 \
#         -v $(pwd)/configs:/app/configs \
#         --add-host host.docker.internal:host-gateway \
#         evenstar-backend:1.0
#
# 说明：
#   - 真实 config.yaml（含密钥）不入镜像，运行时通过 -v 挂载到 /app/configs
#   - config.yaml 里连宿主机 MySQL/Redis 时，host 写 host.docker.internal

# ===== 阶段1：构建 Go 二进制 =====
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 先只拷依赖清单，利用 docker 层缓存（改代码不用重新下载依赖）
COPY go.mod go.sum ./
RUN go mod download

# 拷全部源码并静态编译（项目为纯 Go 依赖，可关闭 CGO）
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o server ./cmd/server

# ===== 阶段2：极小的运行镜像 =====
FROM alpine:3.20

# ca-certificates：SMTP/OSS 等 HTTPS 调用需要根证书
# tzdata：设置容器时区为上海，保证时间日志正确
RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=builder /app/server /app/server

# 后端监听端口（与 configs/config.yaml 的 app.port 一致）
EXPOSE 8080

# -c 显式指定配置路径（挂载目录为 /app/configs/config.yaml）
CMD ["/app/server", "-c", "/app/configs/config.yaml"]
