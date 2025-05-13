# 使用官方 golang 镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o dify-auto-update cmd/main.go


# 创建用于挂载的目录
RUN mkdir -p /app/watch

# 设置环境变量默认值
ENV DIFY_API_KEY=""
ENV DIFY_BASE_URL="http://192.168.101.236:48060"


# 启动命令
CMD ./dify-auto-update watch --api-key=$DIFY_API_KEY --base-url=$DIFY_BASE_URL --folder=/app/watch