package util

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// 从环境变量读取签名密钥，推荐在程序启动时就加载好
var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 为了示例，这里提供了默认值。生产环境请务必通过 env 注入强随机密钥！
		secret = "ToGoOrNotToGo,ItIsAQuestion."
	}
	jwtSecret = []byte(secret)
}

// CustomClaims 定义了我们要在 JWT 中携带的自定义字段
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT 为指定的 userID 生成一个签名好的 tokenString
func GenerateJWT(userID uuid.UUID) (string, error) {
	claims := CustomClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间：72 小时后
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// 你还可以设置 Issuer、Subject、NotBefore 等字段
			//Issuer:    "my-app",
			//Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT 解析并验证 tokenString，返回其中的 userID
func ValidateJWT(tokenString string) (uuid.UUID, error) {
	// 解析 token 时，指定要使用的 CustomClaims 类型
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 强校验签名算法
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	// 校验通过后，类型断言出 CustomClaims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token claims")
	}

	// 解析出 UUID
	uid, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, err
	}
	return uid, nil
}
