package model

import (
	"rbac-admin/utils"

	"gorm.io/gorm"
)

// User 用户表
// 通过 gorm.Model 自动注入 ID、CreatedAt、UpdatedAt、DeletedAt 字段
// 与 Role 通过中间表 user_roles 建立多对多关联
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"` // 用户名，唯一且非空
	Password string `gorm:"not null"`        // 密码哈希，由 BeforeCreate 钩子自动生成
	Email    string                          // 邮箱地址
	Roles    []Role `gorm:"many2many:user_roles;"` // 用户拥有的角色列表
}

// Role 角色表
// 与 Permission 通过中间表 role_permissions 建立多对多关联
type Role struct {
	gorm.Model
	Name        string       `gorm:"unique;not null"`            // 角色名称，唯一且非空
	Description string                                        // 角色描述
	Permissions []Permission `gorm:"many2many:role_permissions;"` // 角色拥有的权限列表
}

// BeforeCreate GORM 钩子：创建用户前自动哈希密码
// 确保密码永远不会以明文形式存入数据库
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashed, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashed
	return nil
}

// Permission 权限表
// 权限是最小粒度的访问控制单元，如 user:list、role:create 等
type Permission struct {
	gorm.Model
	Name        string `gorm:"unique;not null"` // 权限标识名称，如 "user:list"
	Description string                          // 权限描述
}
