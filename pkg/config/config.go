package config

import (
	"encoding/json"
	"math/rand"
	"os"
	"raft/pkg/types"
	"time"
)

type Config struct {
	NodeID      int
	RPCAddr     string
	FriendsAddr map[int]string
	HeartBeat   time.Duration
	TimeMax     time.Duration
	TimeMin     time.Duration
}

func Load(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	er := json.Unmarshal(data, &config)
	if er != nil {
		return nil, err
	}
	return &config, nil
}

//func (c *Config)leixingzhuanhuan() time.Duration {
//	middle:=c.TimeMax-c.TimeMin
//	randN := rand.Intn(middle)
//}

func (c *Config) NodeConfig() types.NodeConfig {
	friends := make([]int, 0, len(c.FriendsAddr))
	for id := range c.FriendsAddr {
		if id != c.NodeID {
			friends = append(friends, id)
		}
	}

	return types.NodeConfig{
		ID:        c.NodeID,
		RPCAddr:   c.RPCAddr,
		Friends:   friends,
		HeartBeat: c.HeartBeat,
		TimeOut: func() time.Duration {
			temp := int64(c.TimeMax - c.TimeMin)
			random := rand.Int63n(temp)
			return c.TimeMin + time.Duration(random)
		}(), //谁整出来的类型真几把麻烦。服了
	}
}
