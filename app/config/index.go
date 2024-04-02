package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// Config 基础配置定义
type Config struct {
	Name                  string   `json:"name"`
	Port                  int      `json:"port"`
	GatewayIdentification string   `json:"gatewayIdentification"`
	Database              Database `json:"database"`
	RedisCommon           Redis    `json:"redis_common"`
}

// Database 数据库定义
type Database struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	Port     int    `json:"port"`
	SSLMode  string `json:"sslmode"`
}

// Redis 数据库定义 redis 如果使用tools下的redisClient 支持多个redis连接 需要Config配置多个该类型的redis配置在初始化的时候传入配置每个redis都是唯一不会重新创建一个连接可以理解成单例模式
type Redis struct {
	Host         string `json:"host"`
	Password     string `json:"password"`
	Db           int    `json:"db"`
	Port         int    `json:"port"`
	MinIdleConns int    `json:"MinIdleConns"`
	PoolSize     int    `json:"poolSize"`
}

var AppConfig *Config

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func LoadConfig(env string) *Config {
	// 根据环境加载不同的配置
	// 这里以简单示例为主，实际项目中可能需要更复杂的配置加载逻辑
	configPath := filepath.Join(getCurrentPath(), env, "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	decoder := json.NewDecoder(file)
	var config Config
	// 解码 JSON 文件
	if err := decodeJSON(decoder, &config); err != nil {
		panic(err)
	}
	return &config
}
func decodeJSON(decoder *json.Decoder, v interface{}) error {
	for {
		if err := decoder.Decode(v); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}
