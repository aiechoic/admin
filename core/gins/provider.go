package gins

import (
	"github.com/aiechoic/admin/core/gin"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/openapi"
)

var Providers = ioc.NewProviders(func(name string, args ...string) *ioc.Provider[*APIServer] {
	return ioc.NewProvider(func(c *ioc.Container) (*APIServer, error) {
		ginConfig := args[0]
		openapiConfig := args[1]
		api, err := openapi.GetOpenapi(openapiConfig, c)
		if err != nil {
			return nil, err
		}
		engine, err := gin.GetServer(ginConfig, c)
		if err != nil {
			return nil, err
		}
		return &APIServer{
			API:   api,
			Engin: engine,
		}, nil
	})
})

func GetAPIServer(ginConfig, openapiConfig string, c *ioc.Container) (*APIServer, error) {
	return Providers.GetProvider("", ginConfig, openapiConfig).Get(c)
}

func GetDefaultAPIServer(c *ioc.Container) (*APIServer, error) {
	return GetAPIServer(gin.DefaultConfig, openapi.DefaultConfig, c)
}
