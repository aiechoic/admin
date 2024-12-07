package verify

import (
	"fmt"
	"github.com/aiechoic/admin/core/email"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/redis"
	"github.com/aiechoic/admin/core/viper"
)

var Providers = ioc.NewProviders[*Verification](func(name string, args ...string) *ioc.Provider[*Verification] {
	return ioc.NewProvider(func(c *ioc.Container) (*Verification, error) {
		redisConfig := args[0]
		emailConfig := args[1]
		rds, err := redis.GetClient(redisConfig, c)
		if err != nil {
			return nil, err
		}
		var sender Sender
		sender, err = email.GetSender(emailConfig, c) // 可改为电话验证, 只需要实现一个电话验证的 Sender 接口
		if err != nil {
			return nil, err
		}
		vp, err := viper.GetViper(name, initConfig, c)
		if err != nil {
			return nil, err
		}
		var cfg Config
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}
		return cfg.NewVerification(rds, sender), nil
	})
})

func GetEmailVerification(name, redisConfig, emailConfig string, c *ioc.Container) (*Verification, error) {
	return Providers.GetProvider(name, redisConfig, emailConfig).Get(c)
}

func GetDefaultVerification(c *ioc.Container) (*Verification, error) {
	return GetEmailVerification(DefaultConfig, redis.DefaultConfig, email.DefaultSenderConfig, c)
}
