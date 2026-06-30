package controller

import (
	"strconv"

	"rbac-admin/service"
	"rbac-admin/utils"

	"github.com/gin-gonic/gin"
)

// UpdateUserRequest 更新用户请求参数
// 字段均为可选，仅更新提供的非空字段
type UpdateUserRequest struct {
	Username string `json:"username"`                    // 新用户名
	Email    string `json:"email" binding:"omitempty,email"` // 新邮箱
}

// AssignRoleRequest 为用户分配角色请求参数
type AssignRoleRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"` // 角色 ID 列表，会替换用户当前所有角色
}

// GetUsers 获取用户列表
// GET /api/v1/users
func GetUsers(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		utils.Error(c, 500, "获取用户列表失败")
		return
	}
	utils.Success(c, users)
}

// GetUser 获取单个用户详情（含角色信息）
// GET /api/v1/users/:id
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的用户ID")
		return
	}
	user, err := service.GetUserByID(uint(id))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}
	utils.Success(c, user)
}

// UpdateUser 更新用户信息（用户名、邮箱）
// PUT /api/v1/users/:id
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的用户ID")
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateUser(uint(id), req.Username, req.Email); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// DeleteUser 删除用户（软删除）
// DELETE /api/v1/users/:id
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的用户ID")
		return
	}
	if err := service.DeleteUser(uint(id)); err != nil {
		utils.Error(c, 404, err.Error())
		return
	}
	utils.Success(c, nil)
}

// AssignRole 为用户分配角色（替换已有角色）
// POST /api/v1/users/:id/roles
func AssignRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的用户ID")
		return
	}
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := service.AssignRoleToUser(uint(id), req.RoleIDs); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}
