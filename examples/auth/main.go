package main

import (
	"context"
	"github.com/aiechoic/admin/examples/auth/src"
	"github.com/aiechoic/admin/internal/gins"
	"github.com/aiechoic/admin/internal/ioc"
	"github.com/aiechoic/admin/src/doc"
)

func main() {
	c := ioc.NewContainer()

	server, err := gins.GetDefaultServer(c)
	if err != nil {
		panic(err)
	}

	auth := src.NewJWTAuth("secret_key")

	server.Register(
		doc.NewDocService(server.API),
		src.NewService(auth),
	)

	server.SetSecuritySchemes(src.SecuritySchemes)

	server.Run(context.Background())
}
