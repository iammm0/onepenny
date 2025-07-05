package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

// jwtSecret 是用于签名JWT的密钥。为了安全起见，应该从配置或环境变量中获取。
var jwtSecret = []byte("ToGoOrNotToGo,ItIsAQuestion.") // 替换为你自己的密钥

// GenerateJWT 生成一个新的JWT，用于用户认证
func GenerateJWT(userID uuid.UUID) (string, error) {
	// 创建一个新的JWT令牌，使用HS256签名方法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),                       // 在令牌中存储用户ID为字符串
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 设置令牌的过期时间为72小时后
	})

	// 使用密钥对令牌进行签名并生成字符串形式的令牌
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT 验证JWT，并返回令牌中的用户ID
func ValidateJWT(tokenString string) (uuid.UUID, error) {
	// 解析并验证令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确认签名方法是否为HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	// 提取令牌中的用户ID
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, errors.New("invalid token claims")
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, err
		}
		return userID, nil
	}

	return uuid.Nil, errors.New("invalid token")
}
