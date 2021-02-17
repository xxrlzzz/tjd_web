package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

// 用户名,获取token的时间, 是否登录
type Claims struct {
	Username string `json:"username"`
	Time string `json:"time"`
	Login bool `json:"login"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(username string, login bool) (string, error) {
	nowTime := time.Now()
	// 过期时间 3h
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		EncodeMD5(username),
		EncodeMD5(nowTime.String()),
		login,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "xxrl",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenClaims.SignedString(jwtSecret)
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
