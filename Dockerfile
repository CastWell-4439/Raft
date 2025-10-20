#构建阶段
FROM golang:1.25.1-alpine AS builder


WORKDIR /raft

#复制go模块文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

#复制
COPY . .

#编译节点服务
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /raft-node ./cmd/node

#运行阶段
FROM alpine:latest


WORKDIR /raft

#从构建阶段复制编译好的可执行文件
COPY --from=builder /raft-node .

COPY config.json .

#RPC服务端口
EXPOSE 8080

#启动！！！
CMD ["./config.json", "-config", "config.json"]