package tools

import (
	Config "EasyGin/app/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var (
	redisClients = make(map[string]*RedisClient)
	mu           sync.Mutex // 用于并发安全的互斥锁
)

type RedisClient struct {
	Client *redis.Client
	Key    string
	pub    *redis.IntCmd //发布
	sub    *redis.PubSub //订阅
}

// ListenData 发布跟监听数据类型
type ListenData struct {
	Key  string
	Cmd  string
	Data interface{}
}
type publishMsg func(*redis.Message)

// Connect 创建连接,如果已经创建过那么直接获取连接
func (r RedisClient) Connect(key string, config Config.Redis) *RedisClient {
	mu.Lock()
	defer mu.Unlock() // 确保并发安全

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
			Key: key,
		}
		return redisClients[key]
	}
}

func (r RedisClient) Set(ctx *gin.Context, key, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	err = r.Client.Set(ctx, fmt.Sprintf("%s:%s", r.Key, key), jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("error Redis Set: %v", err)
	}

	return nil
}

func (r RedisClient) Get(ctx *gin.Context, key string, object interface{}) error {
	jsonStr, err := r.Client.Get(ctx, fmt.Sprintf("%s:%s", r.Key, key)).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonStr), object)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	return nil
}
func (r RedisClient) Expire(ctx *gin.Context, key string, expiration time.Duration) error {
	err := r.Client.Expire(ctx, fmt.Sprintf("%s:%s", r.Key, key), expiration).Err()
	if err != nil {
		return fmt.Errorf("error resetting expiration time: %v", err)
	}
	return nil
}

func (r RedisClient) Destroy(ctx *gin.Context, key string) error {
	err := r.Client.Del(ctx, fmt.Sprintf("%s:%s", r.Key, key)).Err()
	if err != nil {
		return fmt.Errorf("error deleting key from Redis: %v", err)
	}

	return nil
}

// Subscribe 订阅
func (r RedisClient) Subscribe(channels string, funcMsg publishMsg) {
	r.sub = r.Client.Subscribe(context.Background(), channels)
	// 接收订阅的消息
	ch := r.sub.Channel()
	// 启动一个goroutine来处理接收到的消息
	go func() {
		defer func() {
			if r := recover(); r != nil {

			}
		}() // 使用defer recover()来捕获并处理panic
		for msg := range ch {
			funcMsg(msg)
		}
	}()
}

// Publish 发布
func (r RedisClient) Publish(channels string, message ListenData) error {
	// 将结构体转换为 JSON 字符串
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}
	r.pub = r.Client.Publish(context.Background(), channels, jsonData)
	return nil
}
