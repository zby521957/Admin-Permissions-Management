package middleware

import (
	"rbac-admin/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Error(c, 401, "未提供认证令牌")
			c.Abort()
			return
		}
		tokenString := authHeader[7:] // 去掉 "Bearer "
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Error(c, 401, "无效的令牌")
			c.Abort()
			return
		}
		// 将用户信息注入上下文，后续处理器可直接使用
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
