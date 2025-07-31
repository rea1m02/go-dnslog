# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN GOOS=linux GOARCH=amd64 go build -o go-dnslog .

# 生产阶段
FROM alpine:latest


# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app/

# 从构建阶段复制二进制文件
COPY --from=builder /app/go-dnslog .

# 复制配置文件
COPY --from=builder /app/config.yaml .

# 创建数据目录
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080 53

# 启动应用
CMD ["./go-dnslog"] 