package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("your-secret-key") // 生产环境应从环境变量获取

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // 15分钟有效期

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ParseToken 解析和验证JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GenerateRefreshToken 生成刷新令牌（简单实现，实际应更复杂）
func GenerateRefreshToken() string {
	// 实际应用中应使用更安全的方法生成
	return jwt.NewString() // 这是一个伪代码，实际应使用安全的随机生成方法
}
