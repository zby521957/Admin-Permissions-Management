package main

import (
	"log"
	"rbac-admin/cache"
	"rbac-admin/config"
	"rbac-admin/model"
	"rbac-admin/router"
)

func main() {
	// 1. 加载配置文件（config.yaml 或 config.local.yaml）
	if err := config.Load(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 2. 连接 MySQL 数据库
	if err := model.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 2.1 连接 Redis（用于权限缓存）
	if err := cache.InitRedis(); err != nil {
		log.Fatalf("Redis 初始化失败: %v", err)
	}
	log.Println("Redis 连接成功")

	// 3. 自动迁移表结构（GORM AutoMigrate 会根据模型定义自动建表/加字段）
	if err := model.DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 4. 初始化路由（注册所有 API 端点及中间件）
	r := router.SetupRouter()

	// 5. 启动 HTTP 服务
	log.Printf("服务启动在 %s", config.AppConfig.Server.Port)
	if err := r.Run(config.AppConfig.Server.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
