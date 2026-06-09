package service

import (
	"errors"
	"rbac-admin/model"
)

// CreateRole 创建新角色
func CreateRole(name, description string) error {
	var exist model.Role
	if err := model.DB.Where("name = ?", name).First(&exist).Error; err == nil {
		return errors.New("角色名已存在")
	}
	role := model.Role{
		Name:        name,
		Description: description,
	}
	return model.DB.Create(&role).Error
}

// GetAllRoles 获取所有角色（预加载权限）
func GetAllRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := model.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID 根据 ID 获取角色
func GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := model.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, errors.New("角色不存在")
	}
	return &role, nil
}

// UpdateRole 更新角色信息（名称、描述）
func UpdateRole(id uint, name, description string) error {
	var role model.Role
	if err := model.DB.First(&role, id).Error; err != nil {
		return errors.New("角色不存在")
	}
	if name != "" {
		// 检查新名称是否与其他角色冲突
		var exist model.Role
		if err := model.DB.Where("name = ? AND id != ?", name, id).First(&exist).Error; err == nil {
			return errors.New("角色名已被占用")
		}
		role.Name = name
	}
	if description != "" {
		role.Description = description
	}
	return model.DB.Save(&role).Error
}

// DeleteRole 删除角色（软删除）
func DeleteRole(id uint) error {
	var role model.Role
	if err := model.DB.First(&role, id).Error; err != nil {
		return errors.New("角色不存在")
	}
	// 清除关联权限（多对多中间表记录会被自动删除吗？需手动删除以避免悬挂引用）
	if err := model.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		return err
	}
	return model.DB.Delete(&role).Error
}

// AssignPermissionsToRole 为角色分配权限（替换已有权限）
func AssignPermissionsToRole(roleID uint, permIDs []uint) error {
	var role model.Role
	if err := model.DB.First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}
	var perms []model.Permission
	if len(permIDs) > 0 {
		if err := model.DB.Where("id IN ?", permIDs).Find(&perms).Error; err != nil {
			return errors.New("权限ID无效")
		}
	}
	return model.DB.Model(&role).Association("Permissions").Replace(perms)
}
