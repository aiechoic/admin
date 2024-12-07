package src

import (
	"github.com/aiechoic/admin/core/gins"
	"github.com/aiechoic/admin/core/ioc"
)

func NewService(c *ioc.Container) *gins.Service {
	hs := NewHandlers()
	server, err := gins.GetDefaultAPIServer(c)
	if err != nil {
		panic(err)
	}
	baseUrl := server.API.Servers[0].Url
	return &gins.Service{
		Tag:  "User",
		Path: "/user",
		Routes: []gins.Route{
			{
				Method:  "POST",
				Path:    "upload",
				Handler: hs.UploadFile(baseUrl+"/user/file", "uploads"),
			},
			{
				Method:  "GET",
				Path:    "file/:filename",
				Handler: hs.ServeFile("filename", "uploads"),
			},
		},
	}
}
