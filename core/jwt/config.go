package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

const DefaultConfig = "jwt"

var initConfig = `
# JWT configuration file

# secret key, should always set in environment variable
# for example: export JWT_SECRET=your_secret
secret: "your_secret"

# sign method, support HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512
sign_method: "HS256"

# scheme name, this name will show in openapi document auth list
scheme: "user_auth"

# token expires time
expires: "168h"
`

type Config struct {
	Secret     string        `mapstructure:"secret"`
	SignMethod string        `mapstructure:"sign_method"`
	Scheme     string        `mapstructure:"scheme"`
	Expires    time.Duration `mapstructure:"expires"`
}

var auths = sync.Map{}

func newAuth[T any](c *Config) *Auth[T] {
	if v, ok := auths.Load(c); ok {
		return v.(*Auth[T])
	}
	var sm jwt.SigningMethod
	switch c.SignMethod {
	case "HS256":
		sm = jwt.SigningMethodHS256
	case "HS384":
		sm = jwt.SigningMethodHS384
	case "HS512":
		sm = jwt.SigningMethodHS512
	case "RS256":
		sm = jwt.SigningMethodRS256
	case "RS384":
		sm = jwt.SigningMethodRS384
	case "RS512":
		sm = jwt.SigningMethodRS512
	case "ES256":
		sm = jwt.SigningMethodES256
	case "ES384":
		sm = jwt.SigningMethodES384
	case "ES512":
		sm = jwt.SigningMethodES512
	default:
		sm = jwt.SigningMethodHS256
	}
	auth := NewAuth[T](c.Secret, c.Scheme, sm, c.Expires)
	auths.Store(c, auth)
	return auth
}
