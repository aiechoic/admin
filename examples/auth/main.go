package main

import (
	"context"
	"github.com/aiechoic/admin/core/gins"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/examples/auth/src"
	"github.com/aiechoic/admin/src/doc"
)

func main() {
	c := ioc.NewContainer()

	server, err := gins.GetDefaultAPIServer(c)
	if err != nil {
		panic(err)
	}

	server.Register(
		src.NewService(c),
		doc.NewService(server.API),
	)

	server.Run(context.Background())
}
