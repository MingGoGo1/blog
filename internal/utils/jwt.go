package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("default-jwt-secret-please-change-in-production")
var jwtExpireHour = 24 // 默认24小时

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpireHour) * time.Hour * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// SetJWTSecret 设置JWT密钥
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// SetJWTExpireHour 设置JWT过期时间（小时）
func SetJWTExpireHour(hours int) {
	jwtExpireHour = hours
}

// GetJWTExpireDuration 获取JWT过期时间（Duration）
func GetJWTExpireDuration() time.Duration {
	return time.Duration(jwtExpireHour) * time.Hour * 30
}
