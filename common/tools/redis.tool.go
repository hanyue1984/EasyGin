package tools

import (
	Config "EasyGin/app/config"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

var redisClients = make(map[string]RedisClient)

type RedisClient struct {
	client *redis.Client
	key    string
}

var config *Config.Config

// Connect 创建连接,如果已经创建过那么直接获取连接
func (r RedisClient) Connect(key string, config Config.Redis) RedisClient {
	client, exists := redisClients[key]
	if exists {
		return client
	} else {
		redisClients[key] = RedisClient{
			client: redis.NewClient(&redis.Options{
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
	err = r.client.Set(ctx, fmt.Sprintf("%s:%s", r.key, key), jsonData, expiration).Err()
	if err != nil {
		fmt.Println("Error Redis Set:", err)
	}
}

func (r RedisClient) Get(ctx *gin.Context, key, object interface{}) interface{} {
	jsonStr, err := r.client.Get(ctx, fmt.Sprintf("%s:%s", r.key, key)).Result()
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
