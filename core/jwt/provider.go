package jwt

import (
	"fmt"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
)

var Providers = ioc.NewProviders(func(name string, args ...string) *ioc.Provider[*Config] {
	return ioc.NewProvider(func(c *ioc.Container) (*Config, error) {
		vp, err := viper.GetViper(name, initConfig, c)
		if err != nil {
			return nil, err
		}
		var cfg Config
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}
		return &cfg, nil
	})
})

func GetAuth[T any](name string, c *ioc.Container) (*Auth[T], error) {
	cfg, err := Providers.GetProvider(name).Get(c)
	if err != nil {
		return nil, err
	}
	return newAuth[T](cfg), nil
}

func GetDefaultAuth[T any](c *ioc.Container) (*Auth[T], error) {
	return GetAuth[T](DefaultConfig, c)
}
