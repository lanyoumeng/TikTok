package token

import (
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserClaims struct {
	UserId int64
	Name   string
	jwt.RegisteredClaims
}

// Sign 使用 jwtSecret 签发 token，token 的 claims 中会存放传入的 subject.
func Sign(userInfo *UserClaims, jwtKey string) (string, error) {
	// 设置 Token 的内容
	claims := &UserClaims{
		UserId: userInfo.UserId,
		Name:   userInfo.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			//发行者
			Issuer: "TikTok",
			//主题
			Subject: "TikTok_TOKEN",
			//生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
			//发布时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			//过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
		},
	}

	// 创建新的 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签发 Token
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 从 token 中解析出用户信息
func ParseToken(tokenString, jwtKey string) (*UserClaims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		log.Debug("解析Token失败", "err", err)
		return nil, err
	}

	// 检查 Token 有效性
	if !token.Valid {
		log.Debug("无效的 Token")
		return nil, errors.New("无效的 Token")
	}

	// 提取声明
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		log.Debug("无效的声明")
		return nil, errors.New("无效的声明")
	}

	return claims, nil
}
