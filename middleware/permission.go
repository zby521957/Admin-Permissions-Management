package middleware

import (
	"rbac-admin/cache"
	"rbac-admin/model"
	"rbac-admin/utils"

	"github.com/gin-gonic/gin"
)

// RequirePermission 权限校验中间件
// 接收一个权限标识（如 "user:list"），校验当前登录用户是否拥有该权限
// 优先从 Redis 缓存读取权限列表，缓存未命中时查询数据库并回写缓存
// 无权限时返回 403
func RequirePermission(permName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取 JWT 中间件注入的用户 ID
		userID, exists := c.Get("userID")
		if !exists {
			utils.Error(c, 401, "未登录")
			c.Abort()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			utils.Error(c, 401, "认证信息无效")
			c.Abort()
			return
		}

		// 优先从 Redis 缓存读取用户权限
		perms, err := cache.GetUserPermissionsFromCache(uid)
		if err != nil {
			// 缓存未命中，从数据库加载用户及其角色、权限
			var user model.User
			if err := model.DB.Preload("Roles.Permissions").First(&user, uid).Error; err != nil {
				utils.Error(c, 401, "用户不存在")
				c.Abort()
				return
			}
			// 遍历角色 → 权限，收集所有权限标识
			permSet := make([]string, 0)
			for _, role := range user.Roles {
				for _, perm := range role.Permissions {
					permSet = append(permSet, perm.Name)
				}
			}
			// 回写缓存，忽略写入错误（非关键路径）
			_ = cache.CacheUserPermissions(uid, permSet)
			perms = permSet
		}

		// 检查是否拥有目标权限
		hasPerm := false
		for _, p := range perms {
			if p == permName {
				hasPerm = true
				break
			}
		}
		if !hasPerm {
			utils.Error(c, 403, "无操作权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
