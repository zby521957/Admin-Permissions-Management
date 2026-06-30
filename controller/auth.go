package controller

import (
	"rbac-admin/service"
	"rbac-admin/utils"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"` // 用户名，3-64 字符
	Password string `json:"password" binding:"required,min=6"`        // 密码，至少 6 位
	Email    string `json:"email" binding:"omitempty,email"`          // 邮箱，可选但需符合邮箱格式
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// Register 用户注册接口
// POST /api/v1/register
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := service.Register(req.Username, req.Password, req.Email); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// Login 用户登录接口
// POST /api/v1/login
// 验证用户名密码，成功返回 JWT Token
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	token, err := service.Login(req.Username, req.Password)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}
	utils.Success(c, gin.H{"token": token})
}
