# ===== frontend/Dockerfile =====
# Evenstar Writings 前端镜像：Vite 构建静态产物 → 塞进 nginx 镜像托管
# 由 docker-compose.yml 的 web 服务使用：build: ./frontend
#
# 为什么要多阶段：服务器只需装 Docker，不必另外装 Node。
# 构建和运行环境都封在镜像里，本机没 Docker 也不影响。

# ===== 阶段1：构建静态产物 =====
FROM node:20-alpine AS build
WORKDIR /app

# 先只拷依赖清单，利用层缓存：只改业务代码时不会重装依赖
COPY package.json package-lock.json ./

# 大陆服务器装 npm 包走官方 registry 会卡死，改用国内镜像
RUN npm config set registry https://registry.npmmirror.com

RUN npm ci --registry=https://registry.npmmirror.com

# 拷全部源码（node_modules / dist 等已由 .dockerignore 排除）
COPY . .
RUN npm run build

# ===== 阶段2：nginx 托管 =====
FROM nginx:1.27-alpine

# 容器内站点配置（listen 8080 + 反代 backend:8080）
# 覆盖官方默认的 default.conf
COPY nginx.docker.conf /etc/nginx/conf.d/default.conf

# 上一阶段的构建产物
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 8080
