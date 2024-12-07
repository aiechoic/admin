package verify

import (
	"github.com/redis/go-redis/v9"
	"time"
)

const DefaultConfig = "email-verify"

const initConfig = `
# Verify configuration file

# redis key for store verify code
redis_key: "ver_code"

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

func (c *Config) NewVerification(rds *redis.Client, sender Sender) *Verification {
	cache := NewVerifyCodeCache(rds, c.RedisKey, time.Duration(c.Expire)*time.Second)
	return NewVerification(cache, sender, c.RandomCharts, c.Length)
}
