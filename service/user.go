package service

import (
	"errors"
	"rbac-admin/model"
)

// GetUserByID 根据 ID 获取用户
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := model.DB.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// GetAllUsers 获取所有用户列表
func GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := model.DB.Preload("Roles").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser 更新用户信息（用户名、邮箱）
func UpdateUser(id uint, username, email string) error {
	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}
	// 更新字段（仅当新值非空时更新）
	if username != "" {
		// 检查新用户名是否已被占用（排除自身）
		var exist model.User
		if err := model.DB.Where("username = ? AND id != ?", username, id).First(&exist).Error; err == nil {
			return errors.New("用户名已被占用")
		}
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	return model.DB.Save(&user).Error
}

// DeleteUser 删除用户
func DeleteUser(id uint) error {
	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}
	// GORM 默认软删除（如果使用了 gorm.Model）
	return model.DB.Delete(&user).Error
}

// AssignRoleToUser 为用户分配角色（替换已有角色）
func AssignRoleToUser(userID uint, roleIDs []uint) error {
	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}
	var roles []model.Role
	if len(roleIDs) > 0 {
		if err := model.DB.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return errors.New("角色ID无效")
		}
	}
	// 替换关联（清空旧角色，添加新角色）
	if err := model.DB.Model(&user).Association("Roles").Replace(roles); err != nil {
		return err
	}
	return nil
}
