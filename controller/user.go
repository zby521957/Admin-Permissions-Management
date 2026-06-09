package controller

import (
	"strconv"

	"rbac-admin/service"
	"rbac-admin/utils"

	"github.com/gin-gonic/gin"
)

// 定义请求结构体
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type AssignRoleRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		utils.Error(c, 500, "获取用户列表失败")
		return
	}
	utils.Success(c, users)
}

// GetUser 获取单个用户详情
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

// UpdateUser 更新用户
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

// DeleteUser 删除用户
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

// AssignRole 为用户分配角色
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
