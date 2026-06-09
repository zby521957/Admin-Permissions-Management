package service

import (
	"errors"
	"rbac-admin/model"
	"rbac-admin/utils"
)

func Register(username, password, email string) error {
	var existUser model.User
	result := model.DB.Where("username = ?", username).First(&existUser)
	if result.Error == nil {
		return errors.New("用户名已存在")
	}
	user := model.User{
		Username: username,
		Password: password, // 密码在 BeforeCreate 钩子中自动哈希
		Email:    email,
	}
	return model.DB.Create(&user).Error
}

func Login(username, password string) (string, error) {
	var user model.User
	if err := model.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("用户名或密码错误")
	}
	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("用户名或密码错误")
	}
	return utils.GenerateToken(user.ID, user.Username)
}
