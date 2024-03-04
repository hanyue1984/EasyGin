package tools

import (
	Config "EasyGin/app/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

var redisClients = make(map[string]*RedisClient)

type RedisClient struct {
	Client *redis.Client
	key    string
	pub    *redis.IntCmd //发布
	sub    *redis.PubSub //订阅
}

// 发布跟监听数据类型
type listenData struct {
	key  string
	cmd  string
	data interface{}
}
type publishMsg func(*redis.Message)

// Connect 创建连接,如果已经创建过那么直接获取连接
func (r RedisClient) Connect(key string, config Config.Redis) *RedisClient {
	if redisClients[key] != nil {
		return redisClients[key]
	} else {
		redisClients[key] = &RedisClient{
			Client: redis.NewClient(&redis.Options{
				Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
				Password:     config.Password,
				DB:           config.Db,
				PoolSize:     config.PoolSize,
				MinIdleConns: config.MinIdleConns,
			}),
			key: key,
		}
		return redisClients[key]
	}
}

func (r RedisClient) Set(ctx *gin.Context, key, value any, expiration time.Duration) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
	err = r.Client.Set(ctx, fmt.Sprintf("%s:%s", r.key, key), jsonData, expiration).Err()
	if err != nil {
		fmt.Println("Error Redis Set:", err)
	}
}

func (r RedisClient) Get(ctx *gin.Context, key, object interface{}) interface{} {
	jsonStr, err := r.Client.Get(ctx, fmt.Sprintf("%s:%s", r.key, key)).Result()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	// 解析 JSON 数据到结构体

	err = json.Unmarshal([]byte(jsonStr), &object)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}
	return object
}

// Subscribe 订阅
func (r RedisClient) Subscribe(channels string, funcMsg publishMsg) {
	r.sub = r.Client.Subscribe(context.Background(), channels)
	// 接收订阅的消息
	ch := r.sub.Channel()
	// 启动一个 goroutine 来处理接收到的消息
	go func() {
		for msg := range ch {
			funcMsg(msg)
		}
	}()
}

// Publish 发布
func (r RedisClient) Publish(channels string, message listenData) {
	// 将结构体转换为 JSON 字符串
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	r.pub = r.Client.Publish(context.Background(), channels, jsonData)
}
