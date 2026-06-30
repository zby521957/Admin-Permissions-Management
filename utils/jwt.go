package utils

import (
	"errors"
	"rbac-admin/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷结构体
// 包含用户基本信息，嵌入 jwt.RegisteredClaims 提供标准字段
type Claims struct {
	UserID               uint   `json:"user_id"`  // 用户 ID
	Username             string `json:"username"` // 用户名
	jwt.RegisteredClaims                          // 标准 JWT 注册声明（过期时间、签发时间等）
}

// GenerateToken 生成 JWT Token
// 使用 HS256 算法签名，有效期为 24 小时
func GenerateToken(userID uint, username string) (string, error) {
	now := time.Now() // <-- 只取一次时间

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// ParseToken 解析并验证 JWT Token
// 验证签名和有效期，成功返回 Claims，失败返回 error
// 保证返回 (nil, error) 或 (*Claims, nil)，绝不会返回 (nil, nil)
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	// 防范以下场景：
	//   1. token.Claims 类型断言失败（非 *Claims）
	//   2. token.Valid 为 false
	// 此时 err 为 nil，若返回 (nil, nil) 会导致调用方解引用 panic
	return nil, errors.New("无效的令牌")
}
