package service

import (
	"errors"
	"rbac-admin/model"
	"rbac-admin/utils"

	"gorm.io/gorm"
)

// Register 用户注册
// 检查用户名唯一性后创建用户，密码通过 User.BeforeCreate 钩子自动哈希
func Register(username, password, email string) error {
	var existUser model.User
	result := model.DB.Where("username = ?", username).First(&existUser)
	if result.Error == nil {
		return errors.New("用户名已存在")
	}
	// 区分"记录不存在"与真正的数据库错误
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	user := model.User{
		Username: username,
		Password: password, // 明文密码，由 BeforeCreate 钩子自动哈希后写入数据库
		Email:    email,
	}
	return model.DB.Create(&user).Error
}

// Login 用户登录
// 验证用户名密码，成功返回 JWT Token 字符串
func Login(username, password string) (string, error) {
	var user model.User
	if err := model.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("用户名或密码错误")
	}
	// 使用 bcrypt 比对密码哈希
	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("用户名或密码错误")
	}
	return utils.GenerateToken(user.ID, user.Username)
}
