package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config 总配置结构体
// 通过 mapstructure 标签与 YAML 配置文件字段一一映射
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // HTTP 服务相关配置（端口等）
	Database DatabaseConfig `mapstructure:"database"` // MySQL 数据库连接配置
	JWT      JWTConfig      `mapstructure:"jwt"`      // JWT 签名密钥等认证配置
	Redis    RedisConfig    `mapstructure:"redis"`    // Redis 连接与连接池配置
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"` // 监听端口，如 ":8080"
}

// DatabaseConfig MySQL 数据库连接配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`     // 数据库主机地址
	Port     int    `mapstructure:"port"`     // 数据库端口
	User     string `mapstructure:"user"`     // 数据库用户名
	Password string `mapstructure:"password"` // 数据库密码
	DBName   string `mapstructure:"dbname"`   // 数据库名称
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"` // JWT 签名密钥，用于 Token 的签发与验证
}

// RedisConfig Redis 连接配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`      // Redis 服务器地址，如 "127.0.0.1:6379"
	Password string `mapstructure:"password"`  // Redis 认证密码
	DB       int    `mapstructure:"db"`        // Redis 数据库编号（0-15）
	PoolSize int    `mapstructure:"pool_size"` // 连接池大小
}

// AppConfig 全局配置实例，在 Load() 中初始化
// 各模块通过 config.AppConfig.Server / .Database 等方式读取配置
var AppConfig *Config

// Load 加载配置文件
// 默认读取 config.yaml，根据环境变量 ENV 的值切换配置文件：
//   - ENV 未设置或为空   → config.yaml（本地开发，连接本机 MySQL/Redis）
//   - ENV=docker         → config.docker.yaml（Docker Compose 环境，连接容器服务名）
//   - ENV=local          → config.local.yaml（本地调试用自定义配置）
func Load() error {
	env := os.Getenv("ENV")
	configName := "config"
	switch env {
	case "docker":
		configName = "config.docker"
	case "local":
		configName = "config.local"
	}
	viper.SetConfigName(configName) // 设置文件名（不含扩展名）
	viper.SetConfigType("yaml")     // 配置文件格式
	viper.AddConfigPath(".")        // 从当前目录搜索配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return err
	}
	log.Printf("配置加载成功: %s.yaml", configName)

	overrideFromEnv()
	return nil
}

func overrideFromEnv() {
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		AppConfig.Database.Password = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		AppConfig.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			AppConfig.Database.Port = port
		}
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		AppConfig.JWT.Secret = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		AppConfig.Redis.Addr = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		AppConfig.Redis.Password = v
	}
}

