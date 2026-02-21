# Build stage
# 核心修改：Go版本从1.21升级到1.25，对齐go.mod要求
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache gcc musl-dev

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译为静态二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# 运行时必要库
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口 (Render 会自动设置 PORT 环境变量)
EXPOSE 8080

# 运行应用
CMD ["./main"]
