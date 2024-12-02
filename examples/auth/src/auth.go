package src

import (
	"fmt"
	"github.com/aiechoic/admin/core/openapi"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

var SecurityScheme = "userAuth"

var SecuritySchemes = openapi.SecuritySchemes{
	SecurityScheme: &openapi.SecurityScheme{
		Type: "apiKey",
		In:   "header",
		Name: "Authorization",
	},
}

type customClaims struct {
	*User
	jwt.RegisteredClaims
}

type JWTAuth struct {
	secretKey []byte
}

func NewJWTAuth(secretKey string) *JWTAuth {
	return &JWTAuth{
		secretKey: []byte(secretKey),
	}
}

func (j *JWTAuth) GetToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func (j *JWTAuth) Auth(c *gin.Context) {
	token := j.GetToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	user, err := j.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}

func (j *JWTAuth) GenerateToken(user *User) (string, error) {
	// 创建自定义 Claims
	claims := customClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 设置过期时间为 24 小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 签发时间
		},
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名
	return token.SignedString(j.secretKey)
}

func (j *JWTAuth) ParseToken(tokenString string) (*User, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 提取自定义 Claims
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (j *JWTAuth) GetUser(c *gin.Context) *User {
	user, ok := c.Get("user")
	if !ok {
		token := j.GetToken(c)
		u, err := j.ParseToken(token)
		if err != nil {
			return nil
		}
		c.Set("user", u)
		return u
	}
	return user.(*User)
}

func (j *JWTAuth) SecurityScheme() map[string][]string {
	return map[string][]string{SecurityScheme: {}}
}
