# 使用 Golang 官方镜像
FROM golang:1.22-alpine as builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 到工作目录
COPY go.mod main/go.sum ./

# 下载依赖
RUN go mod download

# 复制其他源代码文件到工作目录
COPY main .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o onepenny ./main.go

# 使用 Alpine Linux 作为最终运行镜像
FROM alpine:latest
WORKDIR /root/

# 添加证书（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates

# 从构建阶段复制构建好的二进制文件
COPY --from=builder /app/onepenny .

# 暴露端口
EXPOSE 8080

# 运行程序
CMD ["./onepenny-server"]
