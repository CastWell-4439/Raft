**简介**   
本项目是基于 Go 语言实现的 Raft 共识算法，用于解决分布式系统中的一致性问题。为MIT6.824课程的Lab 2，Lab 1 MapReduce部分也可在本人的主页找到

**项目结构**  

··· plaintext
Raft/
├── cmd/                  # 程序入口
│   ├── node/             # 节点服务入口
│   │   └── main.go       # 节点启动逻辑
│   └── client/           # 客户端工具
│       └── main.go       # 客户端请求逻辑
├── internal/             
│   ├── raft/             
│   │   ├── state.go      # 节点状态管理
│   │   ├── log.go        # 日志处理
│   │   ├── xuanju.go     # 选举机制实现
│   │   ├── fuzhi.go      # 日志复制实现
│   │   └── node.go       # 节点核心逻辑
│   ├── rpcs/             
│   │   ├── client.go     # RPC客户端
│   │   └── server.go     # RPC服务端
│   └── storage/          # 存储
│       └── memory.go     # 内存存储实现（这里没有实现持久化，如果有空的话我会回来补这部分的）
├── pkg/                  
│   ├── config/           
│   │   └── config.go     # 配置
│   └── types/            # 公共类型定义
│       └── common.go     # 状态与RPC类型定义
├── config.json           # 节点配置文件
└── go.mod                
···

**功能**  
核心功能
1. 节点状态管理
    实现了 Raft 节点的三种状态：Follower、Candidate、Leader
    通过 State 结构体维护节点状态信息
2. leader选举
    基于任期的选举机制，任期单调递增
    fllower超时未收到心跳时转为candidate并发起选举
    candidate向其他节点请求投票，获得多数票则成为leader
3. 日志复制
    leader负责将日志条目复制到所有跟随者节点
    通过AppendEntries RPC实现日志同步,包含前置日志检查机制
    leader维护每个follower的NextIndex和MatchIndex以跟踪复制进度
    当日志被多数节点复制后，leader更新提交索引
4. 心跳机制
    leader定期发送心跳（空的 AppendEntries 请求）维持领导地位
    follower收到心跳后重置选举超时计时器  
    可配置的心跳间隔和选举超时范围

**配置说明**
以配置文件 config.json 示例：  
···json
{
  "NodeID": 1,                  // 当前节点ID
  "RPCAddr": "0.0.0.0:8080",    // 当前节点RPC地址
  "FriendsAddr": {              // 其他节点地址映射（ID->地址）
    "2": "node2:8081",
    "3": "node3:8082"
  },
  "HeartBeat": "50ms",          // 心跳间隔
  "TimeMax": "200ms",           // 选举超时最大值
  "TimeMin": "100ms"            // 选举超时最小值
}
···   

**使用方法**  
1.复制配置文件并修改为不同节点的配置：
  ···bash
  cp config.json config-node1.json
  cp config.json config-node2.json
  cp config.json config-node3.json
···  
2.分别启动
···bash
go run cmd/node/main.go -config config-node1.json

go run cmd/node/main.go -config config-node2.json

go run cmd/node/main.go -config config-node3.json
···   
3.使用客户端向领导人节点发送命令  
···bash
go run cmd/client/main.go -serveraddr "0.0.0.0:8080" -command " "
····   

更新：  
1.修改了client中call的和server中register的名字不一样导致无法使用的问题  
2.修改了log.go中check（）的逻辑问题
