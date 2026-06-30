package controller

import (
	"strconv"

	"rbac-admin/service"
	"rbac-admin/utils"

	"github.com/gin-gonic/gin"
)

// CreatePermissionRequest 创建权限请求参数
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"` // 权限标识名称，如 "user:list"
	Description string `json:"description"`             // 权限描述
}

// UpdatePermissionRequest 更新权限请求参数
type UpdatePermissionRequest struct {
	Name        string `json:"name"`        // 新权限名称
	Description string `json:"description"` // 新权限描述
}

// CreatePermission 创建权限
// POST /api/v1/permissions
// 若权限名已存在（含已软删除）则忽略创建
func CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := service.CreatePermission(req.Name, req.Description); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// GetPermissions 获取权限列表
// GET /api/v1/permissions
func GetPermissions(c *gin.Context) {
	perms, err := service.GetAllPermissions()
	if err != nil {
		utils.Error(c, 500, "获取权限列表失败")
		return
	}
	utils.Success(c, perms)
}

// GetPermission 获取单个权限详情
// GET /api/v1/permissions/:id
func GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的权限ID")
		return
	}
	perm, err := service.GetPermissionByID(uint(id))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}
	utils.Success(c, perm)
}

// UpdatePermission 更新权限信息（名称、描述）
// PUT /api/v1/permissions/:id
func UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的权限ID")
		return
	}
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdatePermission(uint(id), req.Name, req.Description); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// DeletePermission 删除权限（先清理角色关联，再软删除）
// DELETE /api/v1/permissions/:id
func DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的权限ID")
		return
	}
	if err := service.DeletePermission(uint(id)); err != nil {
		utils.Error(c, 404, err.Error())
		return
	}
	utils.Success(c, nil)
}
