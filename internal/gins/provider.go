package gins

import (
	"fmt"
	"github.com/aiechoic/admin/internal/ioc"
	"github.com/aiechoic/admin/internal/service"
)

var Providers = ioc.NewProviders[*Server](func(name string) *ioc.Provider[*Server] {
	return ioc.NewProvider(func(c *ioc.Container) (*Server, error) {
		vp, err := service.GetViper(name, InitConfig, c)
		if err != nil {
			return nil, err
		}

		var cfg Config
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}
		return cfg.NewServer(nil), nil
	})
})

func GetServer(name string, c *ioc.Container) (*Server, error) {
	return Providers.GetProvider(name).Get(c)
}

func GetDefaultServer(c *ioc.Container) (*Server, error) {
	return GetServer(DefaultConfig, c)
}
