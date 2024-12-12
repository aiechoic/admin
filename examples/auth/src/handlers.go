package src

import (
	"github.com/aiechoic/admin/core/gins"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/jwt"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	auth *jwt.Auth[User]
}

func NewHandlers(c *ioc.Container) *Handlers {
	auth, err := jwt.GetDefaultAuth[User](c)
	if err != nil {
		panic(err)
	}
	return &Handlers{
		auth: auth,
	}
}

func (hs *Handlers) Login() gins.Handler {
	type params struct {
		Username string `json:"username" binding:"required" description:"The username"`
		Password string `json:"password" binding:"required" description:"The password"`
	}
	type response struct {
		Token string `json:"token" description:"The token"`
	}
	return gins.Handler{
		Request: gins.Request{
			Form: params{},
		},
		Response: gins.Response{
			Json: response{},
		},
		Handle: func(c *gin.Context) {
			var p params
			if err := c.ShouldBind(&p); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			token, err := hs.auth.GenerateToken(&User{
				ID:       1,
				Username: p.Username,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, response{Token: token})
		},
	}
}

func (hs *Handlers) GetInfo() gins.Handler {
	return gins.Handler{
		Response: gins.Response{
			Json: User{},
		},
		Handle: func(c *gin.Context) {
			user := hs.auth.GetUser(c)
			c.JSON(200, user)
		},
	}
}
