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
FROM nginx:1.31-alpine

# 站点配置以「模板」形式放入镜像：容器启动时官方 entrypoint 会执行 envsubst，
# 把 ${SERVER_NAME} 替换成环境变量的值，生成 /etc/nginx/conf.d/default.conf
COPY nginx.docker.conf.template /etc/nginx/templates/default.conf.template

# 兜底默认域名，仅为单独 docker run 本镜像时能通过 nginx 语法检查；
# 实际部署由 docker-compose.yml 通过 environment 覆盖为真实域名
ENV SERVER_NAME=localhost

# 上一阶段的构建产物
COPY --from=build /app/dist /usr/share/nginx/html

# 80/443 为 HTTPS 对外入口，8080 是过渡期明文入口
EXPOSE 80 443 8080
