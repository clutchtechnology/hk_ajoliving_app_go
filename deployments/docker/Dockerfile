# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 复制 go.mod 和 go.sum 先下载依赖（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /api ./cmd/api

# 运行阶段
FROM alpine:3.19

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Hong_Kong

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /api .

# 复制配置文件
COPY configs/config.prod.yaml ./configs/config.yaml

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./api"]
