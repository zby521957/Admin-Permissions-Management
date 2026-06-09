package main

import (
	"log"
	"rbac-admin/config"
	"rbac-admin/model"
	"rbac-admin/router"
)

func main() {
	// 1. 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	// 2. 连接数据库
	if err := model.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	// 3. 自动迁移表结构
	if err := model.DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	// 4. 初始化路由
	r := router.SetupRouter()
	// 5. 启动服务
	log.Printf("服务启动在 %s", config.AppConfig.Server.Port)
	if err := r.Run(config.AppConfig.Server.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
