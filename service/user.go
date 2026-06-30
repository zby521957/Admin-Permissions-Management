package service

import (
	"errors"

	"rbac-admin/cache"
	"rbac-admin/model"

	"gorm.io/gorm"
)

// GetUserByID 根据 ID 获取用户（含角色信息）
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := model.DB.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// GetAllUsers 获取所有用户列表（含角色信息）
func GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := model.DB.Preload("Roles").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser 更新用户信息（用户名、邮箱）
// 仅更新提供的非空字段，更新用户名时会检查唯一性
func UpdateUser(id uint, username, email string) error {
	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}
	if username != "" {
		// 检查新用户名是否已被其他用户占用
		var exist model.User
		if err := model.DB.Where("username = ? AND id != ?", username, id).First(&exist).Error; err == nil {
			return errors.New("用户名已被占用")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	// Save 会触发 GORM 的 BeforeSave 钩子（如有），并更新所有字段
	return model.DB.Save(&user).Error
}

// DeleteUser 删除用户（GORM 软删除，设置 deleted_at 时间戳）
func DeleteUser(id uint) error {
	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}
	return model.DB.Delete(&user).Error
}

// AssignRoleToUser 为用户分配角色（替换已有角色）
// 会清空用户当前的角色关联，然后关联新的角色列表
// 分配成功后使该用户的权限缓存失效
func AssignRoleToUser(userID uint, roleIDs []uint) error {
	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}
	// 校验角色 ID 有效性
	var roles []model.Role
	if len(roleIDs) > 0 {
		if err := model.DB.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}
		if len(roles) != len(roleIDs) {
			return errors.New("角色ID无效")
		}
	}
	// Replace 会清空旧的角色关联并建立新的关联
	if err := model.DB.Model(&user).Association("Roles").Replace(roles); err != nil {
		return err
	}
	// 使该用户的权限缓存失效，下次请求重新加载
	_ = cache.InvalidateUserCache(userID)
	return nil
}
