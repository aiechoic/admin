package src

import (
	"github.com/aiechoic/admin/core/gins"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/jwt"
)

func NewService(c *ioc.Container) *gins.Service {
	auth, err := jwt.GetDefaultAuth[User](c)
	if err != nil {
		panic(err)
	}
	hs := NewHandlers(c)
	return &gins.Service{
		Tag:  "User",
		Path: "/user",
		Routes: []gins.Route{
			{
				Method:  "POST",
				Path:    "login",
				Handler: hs.Login(),
			},
			{
				Method:   "GET",
				Path:     "info",
				Security: auth,
				Handler:  hs.GetInfo(),
			},
		},
	}
}
