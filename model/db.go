package model

import (
	"fmt"
	"log"
	"rbac-admin/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局数据库连接实例，通过 InitDB 初始化后可在各处使用
var DB *gorm.DB

// InitDB 初始化 MySQL 数据库连接
// 从配置中读取数据库参数，拼接 DSN 后建立连接
func InitDB() error {
	cfg := config.AppConfig.Database
	// 拼接 DSN (Data Source Name)，启用 utf8mb4 字符集、自动解析时间、本地时区
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Println("数据库连接成功")
	return nil
}
