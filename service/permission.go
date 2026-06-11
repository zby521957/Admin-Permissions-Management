package service

import (
	"errors"
	"rbac-admin/model"

	"gorm.io/gorm/clause"
)

// CreatePermission 创建权限，若名称已存在（包括软删除的记录）则忽略
func CreatePermission(name, description string) error {
	perm := model.Permission{
		Name:        name,
		Description: description,
	}
	result := model.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&perm)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("权限名已存在")
	}
	return nil
}

// GetAllPermissions 获取所有权限
func GetAllPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	if err := model.DB.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// GetPermissionByID 根据ID获取权限
func GetPermissionByID(id uint) (*model.Permission, error) {
	var perm model.Permission
	if err := model.DB.First(&perm, id).Error; err != nil {
		return nil, errors.New("权限不存在")
	}
	return &perm, nil
}

// UpdatePermission 更新权限
func UpdatePermission(id uint, name, description string) error {
	var perm model.Permission
	if err := model.DB.First(&perm, id).Error; err != nil {
		return errors.New("权限不存在")
	}
	if name != "" {
		var exist model.Permission
		if err := model.DB.Where("name = ? AND id != ?", name, id).First(&exist).Error; err == nil {
			return errors.New("权限名已被占用")
		}
		perm.Name = name
	}
	if description != "" {
		perm.Description = description
	}
	return model.DB.Save(&perm).Error
}

// DeletePermission 删除权限（清理关联后软删除）
func DeletePermission(id uint) error {
	var perm model.Permission
	if err := model.DB.First(&perm, id).Error; err != nil {
		return errors.New("权限不存在")
	}
	if err := model.DB.Exec("DELETE FROM role_permissions WHERE permission_id = ?", id).Error; err != nil {
		return err
	}
	return model.DB.Delete(&perm).Error
}
