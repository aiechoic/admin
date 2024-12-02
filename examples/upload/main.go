package main

import (
	"context"
	"github.com/aiechoic/admin/core/gins"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/examples/upload/src"
	"github.com/aiechoic/admin/src/doc"
)

func main() {
	c := ioc.NewContainer()

	server, err := gins.GetDefaultServer(c)
	if err != nil {
		panic(err)
	}

	server.Register(
		doc.NewDocService(server.API),
		src.NewService(c),
	)

	server.Run(context.Background())
}
