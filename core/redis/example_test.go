package redis_test

import (
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/redis"
)

func ExampleGetClient() {
	c := ioc.NewContainer()

	client, err := redis.GetClient("redis", c)
	if err != nil {
		panic(err)
	}
	_ = client
}
