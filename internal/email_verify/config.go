package email_verify

import (
	"github.com/redis/go-redis/v9"
	"time"
)

const DefaultConfig = "email-verify"

const initConfig = `
# email verify config

# redis key for store email code
redis_key: "email_ver_code"

# code expire time in seconds
expire: 300

# random charts for generate code
random_charts: "0123456789"

# code length
length: 6
`

type Config struct {
	RedisKey     string `mapstructure:"redis_key"`
	Expire       int    `mapstructure:"expire"`
	RandomCharts string `mapstructure:"random_charts"`
	Length       int    `mapstructure:"length"`
}

func (c *Config) NewEmailVerification(rds *redis.Client, sender EmailSender) *EmailVerification {
	cache := NewVerifyCodeCache(rds, c.RedisKey, time.Duration(c.Expire)*time.Second)
	return NewEmailVerification(cache, sender, c.RandomCharts, c.Length)
}
