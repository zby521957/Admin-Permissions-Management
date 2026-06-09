package model

import "gorm.io/gorm"

// User 用户表
type User struct {
	//自动添加 ID、CreatedAt、UpdatedAt、DeletedAt
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string
	Roles    []Role `gorm:"many2many:user_roles;"` // many2many:user_roles; 表示通过中间表 user_roles 建立多对多关系
}

// Role 角色表
type Role struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"` // 角色与权限多对多
}

// Permission 权限表
type Permission struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string //权限描述
}
