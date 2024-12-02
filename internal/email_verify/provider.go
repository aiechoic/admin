package email_verify

import (
	"fmt"
	"github.com/aiechoic/admin/internal/ioc"
	"github.com/aiechoic/admin/internal/service"
)

var providers = ioc.NewProviders[*EmailVerification](func(name string) *ioc.Provider[*EmailVerification] {
	return ioc.NewProvider(func(c *ioc.Container) (*EmailVerification, error) {
		rds, err := service.GetDefaultRedisClient(c)
		if err != nil {
			return nil, err
		}
		sender, err := service.GetDefaultEmailSender(c)
		if err != nil {
			return nil, err
		}
		vp, err := service.GetViper(name, initConfig, c)
		if err != nil {
			return nil, err
		}
		var cfg Config
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}
		return cfg.NewEmailVerification(rds, sender), nil
	})
})

func GetEmailVerification(name string, c *ioc.Container) (*EmailVerification, error) {
	return providers.GetProvider(name).Get(c)
}

func GetDefaultVerification(c *ioc.Container) (*EmailVerification, error) {
	return GetEmailVerification(DefaultConfig, c)
}
