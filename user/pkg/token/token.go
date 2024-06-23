package token

import (
	"errors"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserInfo struct {
	UserId int64  `json:"id" `
	Name   string `json:"name" valid:"alphanum,required,stringlength(6|32)"`
}

// Sign 使用 jwtSecret 签发 token，token 的 claims 中会存放传入的 subject.
func Sign(userInfo *UserInfo, jwtKey string) (tokenString string, err error) {
	// 将 UserInfo 转换为 map[string]interface{}
	userInfoMap := map[string]interface{}{
		"UserId": userInfo.UserId,
		"Name":   userInfo.Name,
	}

	// Token 的内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userInfo": userInfoMap,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(48 * time.Hour).Unix(),
	})
	// 签发 token
	tokenString, err = token.SignedString([]byte(jwtKey))

	return
}

// ParseToken 从 token 中解析出用户信息
func ParseToken(tokenString, jwtKey string) (*UserInfo, error) {
	// 解析 Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	// 检查 Token 有效性
	if !token.Valid {
		return nil, errors.New("无效的 Token")
	}

	// 提取声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("无效的声明")
	}

	// 从声明中提取用户信息
	userInfoMap, ok := claims["userInfo"].(map[string]interface{})
	if !ok {
		return nil, errors.New("用户信息格式无效")
	}

	// 构建 UserInfo 结构体
	user := &UserInfo{}

	userIDValue, ok := userInfoMap["UserId"]
	if !ok {
		return nil, errors.New("缺少用户ID")
	}

	// 使用反射进行类型断言和转换
	userIDType := reflect.TypeOf(userIDValue).Kind()
	switch userIDType {
	case reflect.Float64:
		user.UserId = int64(userIDValue.(float64))
	case reflect.Int64:
		user.UserId = userIDValue.(int64)
	default:
		return nil, errors.New("用户ID类型无效")
	}

	user.Name, ok = userInfoMap["Name"].(string)
	if !ok {
		return nil, errors.New("用户名类型无效")
	}

	return user, nil
}
