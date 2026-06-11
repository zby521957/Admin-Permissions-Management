package controller

import (
	"strconv"

	"rbac-admin/service"
	"rbac-admin/utils"

	"github.com/gin-gonic/gin"
)

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

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

func GetPermissions(c *gin.Context) {
	perms, err := service.GetAllPermissions()
	if err != nil {
		utils.Error(c, 500, "获取权限列表失败")
		return
	}
	utils.Success(c, perms)
}

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

func DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "无效的权限ID")
		return
	}
	if err := service.DeletePermission(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, nil)
}
