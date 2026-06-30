package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rbac-admin/config"

	"github.com/redis/go-redis/v9"
)

// Rdb 全局 Redis 客户端实例，通过 InitRedis 初始化后可在各处使用
var Rdb *redis.Client

// InitRedis 初始化 Redis 连接
// 从配置中读取参数并建立连接池，支持连接池大小配置
func InitRedis() error {
	cfg := config.AppConfig.Redis
	Rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     // Redis 服务器地址
		Password: cfg.Password, // Redis 认证密码
		DB:       cfg.DB,       // 使用的数据库编号
		PoolSize: cfg.PoolSize, // 连接池大小
	})
	ctx := context.Background()
	if err := Rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis 连接失败: %w", err)
	}
	return nil
}

// userPermKey 生成用户权限缓存的 Redis key
// 格式：user:perms:<userID>
func userPermKey(userID uint) string {
	return fmt.Sprintf("user:perms:%d", userID)
}

// CacheUserPermissions 将用户权限列表缓存到 Redis
// 缓存有效期为 1 小时，过期后需重新从数据库加载
func CacheUserPermissions(userID uint, permissions []string) error {
	ctx := context.Background()
	data, _ := json.Marshal(permissions)
	return Rdb.Set(ctx, userPermKey(userID), data, 1*time.Hour).Err()
}

// GetUserPermissionsFromCache 从 Redis 缓存中获取用户权限列表
// 返回 error 表示缓存未命中或已过期，调用方应回源数据库查询
func GetUserPermissionsFromCache(userID uint) ([]string, error) {
	ctx := context.Background()
	val, err := Rdb.Get(ctx, userPermKey(userID)).Result()
	if err != nil {
		return nil, err
	}
	var perms []string
	if err := json.Unmarshal([]byte(val), &perms); err != nil {
		return nil, err
	}
	return perms, nil
}

// InvalidateUserCache 删除指定用户的权限缓存
// 当用户角色或权限发生变更时调用，确保下次请求从数据库重新加载
func InvalidateUserCache(userID uint) error {
	ctx := context.Background()
	return Rdb.Del(ctx, userPermKey(userID)).Err()
}
