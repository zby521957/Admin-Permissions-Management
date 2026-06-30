package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 算法哈希密码
// cost 使用默认值（10），在安全性和性能之间取得平衡
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证明文密码是否与 bcrypt 哈希匹配
// 返回 true 表示密码正确
func CheckPassword(plain, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
