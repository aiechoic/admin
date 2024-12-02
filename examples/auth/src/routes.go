package src

import (
	"github.com/aiechoic/admin/core/gins"
)

func NewService(auth *JWTAuth) *gins.Service {
	hs := NewHandlers(auth)
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
