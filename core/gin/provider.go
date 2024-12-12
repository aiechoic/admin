package gin

import (
	"fmt"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
)

var Providers = ioc.NewProviders[*Server](func(name string, args ...any) *ioc.Provider[*Server] {
	return ioc.NewProvider(func(c *ioc.Container) (*Server, error) {
		vp, err := viper.GetViper(name, initConfig, c)
		if err != nil {
			return nil, err
		}
		var cfg Config
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}
		ginEngine, iRouter := cfg.NewGinEngine()
		return &Server{
			Engine:    ginEngine,
			ApiRouter: iRouter,
			HttpPort:  cfg.HttpPort,
		}, nil
	})
})

func GetServer(name string, c *ioc.Container) (*Server, error) {
	return Providers.GetProvider(name).Get(c)
}

func GetDefaultServer(c *ioc.Container) (*Server, error) {
	return GetServer(DefaultConfig, c)
}
