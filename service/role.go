package service

import (
	"errors"

	"rbac-admin/cache"
	"rbac-admin/model"

	"gorm.io/gorm"
)

// CreateRole 创建新角色
// 检查角色名唯一性后创建
func CreateRole(name, description string) error {
	var exist model.Role
	if err := model.DB.Where("name = ?", name).First(&exist).Error; err == nil {
		return errors.New("角色名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	role := model.Role{
		Name:        name,
		Description: description,
	}
	return model.DB.Create(&role).Error
}

// GetAllRoles 获取所有角色（预加载关联的权限列表）
func GetAllRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := model.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID 根据 ID 获取角色（含权限信息）
func GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := model.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, errors.New("角色不存在")
	}
	return &role, nil
}

// UpdateRole 更新角色信息（名称、描述）
// 仅更新提供的非空字段，更新名称时会检查唯一性
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
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		role.Name = name
	}
	if description != "" {
		role.Description = description
	}
	return model.DB.Save(&role).Error
}

// DeleteRole 删除角色（先清除权限关联，再软删除）
// 先清除中间表 role_permissions 中的关联记录，避免悬挂引用
func DeleteRole(id uint) error {
	var role model.Role
	if err := model.DB.First(&role, id).Error; err != nil {
		return errors.New("角色不存在")
	}
	// 清除该角色与所有权限的关联关系
	if err := model.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		return err
	}
	// 使拥有该角色的所有用户权限缓存失效
	invalidateUsersWithRole(id) // <-- 加这一行
	return model.DB.Delete(&role).Error
}

// AssignPermissionsToRole 为角色分配权限（替换已有权限）
// 会清空角色当前的权限关联，然后关联新的权限列表
// 分配成功后使所有拥有该角色的用户缓存失效
func AssignPermissionsToRole(roleID uint, permIDs []uint) error {
	var role model.Role
	if err := model.DB.First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}
	// 校验权限 ID 有效性
	var perms []model.Permission
	if len(permIDs) > 0 {
		if err := model.DB.Where("id IN ?", permIDs).Find(&perms).Error; err != nil {
			return err
		}
		if len(perms) != len(permIDs) {
			return errors.New("权限ID无效")
		}
	}
	err := model.DB.Model(&role).Association("Permissions").Replace(perms)
	if err != nil {
		return err
	}
	// 使所有拥有该角色的用户缓存失效
	invalidateUsersWithRole(roleID)
	return nil
}

// invalidateUsersWithRole 使所有拥有指定角色的用户权限缓存失效
func invalidateUsersWithRole(roleID uint) {
	var users []model.User
	if err := model.DB.Where("id IN (SELECT user_id FROM user_roles WHERE role_id = ?)", roleID).Find(&users).Error; err == nil {
		for _, u := range users {
			_ = cache.InvalidateUserCache(u.ID)
		}
	}
}
