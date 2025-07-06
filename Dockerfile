# 第一阶段: 构建
FROM golang:1.24-alpine AS builder

# 安装必要工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 拷贝进来
COPY api/go.mod api/go.sum ./

# 下载依赖
RUN go mod download

# 拷贝源代码
COPY api/. .

RUN go build -o app .

# 第二阶段: 生成最终镜像
FROM alpine:latest

# 时区可选
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从 builder 拷贝编译好的二进制文件
COPY --from=builder /app/app .

# 开放端口
EXPOSE 8080

# 启动
ENTRYPOINT ["./app"]