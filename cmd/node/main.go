package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"raft/internal/raft"
	"raft/internal/rpcs"
	"raft/pkg/config"
	"syscall"
)

func main() {
	//依旧命令行参数
	path := flag.String("config", "config.json", "path")
	flag.Parse()

	//加载配置
	cfg, err := config.Load(*path)
	if err != nil {

		log.Fatal("fail to load config:%v", err)
	}
	node := raft.NewNode(cfg.NodeConfig(), cfg.FriendsAddr)

	//创建服务器
	server := rpcs.NewRPCServer(node)
	if err := server.Start(cfg.RPCAddr); err != nil {
		log.Fatal("fail to start server:%v", err)
	}
	node.Run()

	//监听SIGINT和SIGTERM
	zhongduan := make(chan os.Signal, 1)
	signal.Notify(zhongduan, syscall.SIGINT, syscall.SIGTERM)
	<-zhongduan

	node.Close()
}
