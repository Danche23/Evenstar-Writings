# ===== backend/Dockerfile =====
# Evenstar Writings 后端镜像 —— 「零编译」版
#
# 为什么不用多阶段构建（2026-09-18 改）：
#   本项目服务器是 1.6G 小内存机型。在容器里跑 go build 会把内存吃满，
#   即使有 swap 兜底，也会把整机拖到近乎无响应（构建过程看起来像「服务器崩了」）。
#   本方案把编译搬到开发机，服务器只做 COPY：
#
#     本机交叉编译  →  上传二进制  →  服务器 build（数秒，内存零压力）
#
# 用法（两步）：
#   1) 本机编译（Windows / macOS / Linux 均可，需 Go 1.24+）：
#        cd backend
#        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o server ./cmd/server
#      ⚠️ 产物文件名必须是 server，且与本 Dockerfile 同目录（构建上下文）
#      ⚠️ Git Bash 下用 go 这类 Windows 原生程序时，输出路径别写 /d/...（会被当相对路径），
#         写 D:/... 或相对路径
#   2) 服务器上构建并重启：
#        cd /var/www/evenstar
#        docker compose build backend && docker compose up -d backend
#
# 想恢复「服务器内编译」：换回 git 历史里的多阶段 Dockerfile 即可（需 ≥4G 内存才顺畅）。

FROM alpine:3.24

# ca-certificates：SMTP/OSS 等 HTTPS 调用需要根证书
# tzdata：时区设为上海，保证日志时间正确
# mkdir：alpine 里没有这两个目录，预先建好（lumberjack 写日志 / 本地存储兜底都要写权限）
RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && mkdir -p /app/logs /app/storage

WORKDIR /app

# 直接拷入本机编译好的 Linux 二进制 —— 不执行任何编译
# chmod 是保险：文件从 Windows 打包/传输过来时可能丢掉执行位，
# 表现为容器启动报 "permission denied" 并反复重启
COPY server /app/server
RUN chmod +x /app/server

# 后端监听端口（与 configs/config.yaml 的 app.port 一致）
EXPOSE 8080

# -c 显式指定配置路径（运行时由 compose 挂载到 /app/configs/config.yaml）
CMD ["/app/server", "-c", "/app/configs/config.yaml"]
