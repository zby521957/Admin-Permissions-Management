package middleware

import (
	"rbac-admin/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
// 从请求头 Authorization 中提取 Bearer Token，解析后将用户信息注入上下文
// 认证失败时返回 401 并中断请求链
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取 Authorization 头，校验 Bearer 前缀
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Error(c, 401, "未提供认证令牌")
			c.Abort()
			return
		}
		// 解析 Token，去掉 "Bearer " 前缀（7 个字符）
		tokenString := authHeader[7:]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Error(c, 401, "无效的令牌")
			c.Abort()
			return
		}
		// 将用户信息注入上下文，后续处理器可直接使用 c.Get("userID") / c.Get("username")
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
