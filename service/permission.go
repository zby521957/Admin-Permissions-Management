package service

import (
	"errors"

	"rbac-admin/cache"
	"rbac-admin/model"

	"gorm.io/gorm"
)

// CreatePermission 创建权限
// 先检查权限名是否存在（含软删除记录），不存在则创建
func CreatePermission(name, description string) error {
	// 先检查是否存在同名权限（GORM 自动排除软删除记录，但 unique 约束包含软删除记录）
	var exist model.Permission
	err := model.DB.Unscoped().Where("name = ?", name).First(&exist).Error
	if err == nil {
		return errors.New("权限名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	perm := model.Permission{
		Name:        name,
		Description: description,
	}
	return model.DB.Create(&perm).Error
}

// GetAllPermissions 获取所有权限列表
func GetAllPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	if err := model.DB.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// GetPermissionByID 根据 ID 获取权限
func GetPermissionByID(id uint) (*model.Permission, error) {
	var perm model.Permission
	if err := model.DB.First(&perm, id).Error; err != nil {
		return nil, errors.New("权限不存在")
	}
	return &perm, nil
}

// UpdatePermission 更新权限信息（名称、描述）
// 仅更新提供的非空字段，更新名称时会检查唯一性
// 如果权限名称变更，会使所有引用该权限的用户缓存失效
func UpdatePermission(id uint, name, description string) error {
	var perm model.Permission
	if err := model.DB.First(&perm, id).Error; err != nil {
		return errors.New("权限不存在")
	}
	nameChanged := false
	if name != "" {
		// 检查新名称是否与其他权限冲突
		var exist model.Permission
		if err := model.DB.Unscoped().Where("name = ? AND id != ?", name, id).First(&exist).Error; err == nil {
			return errors.New("权限名已被占用")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		nameChanged = perm.Name != name
		perm.Name = name
	}
	if description != "" {
		perm.Description = description
	}
	err := model.DB.Save(&perm).Error
	if err != nil {
		return err
	}
	// 如果权限名称变更，使所有拥有该权限的用户缓存失效
	if nameChanged {
		invalidateUsersWithPermission(id)
	}
	return nil
}

// DeletePermission 删除权限（先清理角色关联，再软删除）
// 先手动删除中间表 role_permissions 中的关联记录，避免悬挂引用
// 删除后使所有拥有该权限的用户缓存失效
func DeletePermission(id uint) error {
	var perm model.Permission
	if err := model.DB.First(&perm, id).Error; err != nil {
		return errors.New("权限不存在")
	}
	// 使所有拥有该权限的用户缓存失效
	invalidateUsersWithPermission(id)
	// 清理 role_permissions 中间表中引用此权限的记录
	if err := model.DB.Exec("DELETE FROM role_permissions WHERE permission_id = ?", id).Error; err != nil {
		return err
	}
	return model.DB.Delete(&perm).Error
}

// invalidateUsersWithPermission 使所有通过角色间接拥有指定权限的用户缓存失效
func invalidateUsersWithPermission(permID uint) {
	var users []model.User
	err := model.DB.Raw(`
		SELECT DISTINCT u.* FROM users u
		INNER JOIN user_roles ur ON u.id = ur.user_id
		INNER JOIN role_permissions rp ON ur.role_id = rp.role_id
		WHERE rp.permission_id = ?
	`, permID).Scan(&users).Error
	if err == nil {
		for _, u := range users {
			_ = cache.InvalidateUserCache(u.ID)
		}
	}
}
