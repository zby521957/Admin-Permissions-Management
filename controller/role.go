package controller

import (
	"strconv"

	"rbac-admin/service"
	"rbac-admin/utils"

	"github.com/gin-gonic/gin"
)

// CreateRoleRequest 创建角色请求参数
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"` // 角色名称，必填
	Description string `json:"description"`             // 角色描述
}

// UpdateRoleRequest 更新角色请求参数
type UpdateRoleRequest struct {
	Name        string `json:"name"`        // 新角色名称
	Description string `json:"description"` // 新角色描述
}

// AssignPermissionsRequest 为角色分配权限请求参数
type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"` // 权限 ID 列表
}

// CreateRole 创建角色
// POST /api/v1/roles
func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := service.CreateRole(req.Name, req.Description); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// GetRoles 获取角色列表（含权限信息）
// GET /api/v1/roles
func GetRoles(c *gin.Context) {
	roles, err := service.GetAllRoles()
	if err != nil {
		utils.Error(c, 500, "获取角色列表失败")
		return
	}
	utils.Success(c, roles)
}

// GetRole 获取单个角色详情（含权限信息）
// GET /api/v1/roles/:id
func GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的角色ID")
		return
	}
	role, err := service.GetRoleByID(uint(id))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}
	utils.Success(c, role)
}

// UpdateRole 更新角色信息（名称、描述）
// PUT /api/v1/roles/:id
func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的角色ID")
		return
	}
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateRole(uint(id), req.Name, req.Description); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// DeleteRole 删除角色（先清除关联权限，再软删除）
// DELETE /api/v1/roles/:id
func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的角色ID")
		return
	}
	if err := service.DeleteRole(uint(id)); err != nil {
		utils.Error(c, 404, err.Error())
		return
	}
	utils.Success(c, nil)
}

// AssignPermissions 为角色分配权限（替换已有权限）
// POST /api/v1/roles/:id/permissions
func AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的角色ID")
		return
	}
	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := service.AssignPermissionsToRole(uint(id), req.PermissionIDs); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}
