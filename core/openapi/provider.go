package openapi

import (
	"fmt"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
)

const DefaultConfig = "openapi"

var initConfig = `
# OpenAPI Initial Configuration
# see more info: https://swagger.io/specification/#openapi-object.

# openapi version
openapi: "3.0.0"

# info
info:
  title: "API Documentation"
  description: ""
  version: "1.0.0"

# api servers
servers:
  - url: "http://localhost:8080/api/v1"
    description: "Local Server"

# no need to set other fields, it will be generated automatically.
`

var Providers = ioc.NewProviders[*Openapi](func(name string, a ...any) *ioc.Provider[*Openapi] {
	return ioc.NewProvider(func(c *ioc.Container) (*Openapi, error) {
		vp, err := viper.GetViper(name, initConfig, c)
		if err != nil {
			return nil, err
		}

		var api Openapi
		err = vp.Unmarshal(&api)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}
		if api.Paths == nil {
			api.Paths = make(map[string]PathItem)
			api.Components = Components{
				Schemas:         make(map[string]*Schema),
				SecuritySchemes: make(SecuritySchemes),
			}
		}
		return &api, nil
	})
})

func GetOpenapi(name string, c *ioc.Container) (*Openapi, error) {
	return Providers.GetProvider(name).Get(c)
}

func GetDefaultOpenapi(c *ioc.Container) (*Openapi, error) {
	return GetOpenapi(DefaultConfig, c)
}
