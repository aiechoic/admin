package jwt

import (
	"fmt"
	"github.com/aiechoic/admin/core/openapi"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type customClaims[T any] struct {
	T *T
	jwt.RegisteredClaims
}

type Auth[T any] struct {
	secret  []byte
	method  jwt.SigningMethod
	ctxKey  string
	scheme  string
	expires time.Duration
}

func NewAuth[T any](secret, scheme string, signMethod jwt.SigningMethod, expires time.Duration) *Auth[T] {
	return &Auth[T]{
		secret:  []byte(secret),
		method:  signMethod,
		scheme:  scheme,
		ctxKey:  uuid.NewString(),
		expires: expires,
	}
}

func (j *Auth[T]) GetToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func (j *Auth[T]) Auth(c *gin.Context) {
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

	c.Set(j.ctxKey, user)
	c.Next()
}

func (j *Auth[T]) GenerateToken(user *T) (string, error) {
	// 创建自定义 Claims
	claims := customClaims[T]{
		T: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expires)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                // 签发时间
		},
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名
	return token.SignedString(j.secret)
}

func (j *Auth[T]) ParseToken(tokenString string) (*T, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &customClaims[T]{}, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	// 提取自定义 Claims
	if claims, ok := token.Claims.(*customClaims[T]); ok && token.Valid {
		return claims.T, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (j *Auth[T]) GetUser(c *gin.Context) *T {
	// try to get user from context
	user, ok := c.Get(j.ctxKey)
	if !ok {
		token := j.GetToken(c)
		u, err := j.ParseToken(token)
		if err != nil {
			return nil
		}
		c.Set(j.ctxKey, u)
		return u
	}
	return user.(*T)
}

// SecuritySchemes defines the security schemes for the auth, used for openapi generation
// see details in https://swagger.io/specification/#components-object
func (j *Auth[T]) SecuritySchemes() openapi.SecuritySchemes {
	return openapi.SecuritySchemes{
		j.scheme: &openapi.SecurityScheme{
			Type: "apiKey",
			In:   "header",
			Name: "Authorization",
		},
	}
}

// SecurityRequirement defines the security requirement for the auth, used for openapi generation
// see details in https://swagger.io/specification/#security-requirement-object
func (j *Auth[T]) SecurityRequirement() map[string][]string {
	return map[string][]string{
		j.scheme: {},
	}
}
