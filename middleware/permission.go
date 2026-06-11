package middleware

import (
	"rbac-admin/model"
	"rbac-admin/utils"

	"github.com/gin-gonic/gin"
)

// RequirePermission 检查当前用户是否拥有指定权限
func RequirePermission(permName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 JWT 中间件注入的上下文中获取 userID
		userID, exists := c.Get("userID")
		if !exists {
			utils.Error(c, 401, "未登录")
			c.Abort()
			return
		}

		// 查询用户 → 角色 → 权限
		var user model.User
		if err := model.DB.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
			utils.Error(c, 401, "用户不存在")
			c.Abort()
			return
		}

		// 检查用户的所有角色中是否包含所需权限
		hasPermission := false
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				if perm.Name == permName {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			utils.Error(c, 403, "无操作权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
