package redis

import (
	"context"
	"fmt"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
)

const DefaultConfig = "redis"

var initConfig = `
# Redis configuration file
#
# configure "redis.conf" for add user and password:
# user yourusername on +@all ~* >somepassword

# use debug mode
debug: true

# Redis server
host: "localhost"

# Redis port
port: 6379

# Redis database
db: 9

# Username for Redis server
username: ""

# Password for Redis server
password: ""
`

type Config struct {
	Debug    bool   `mapstructure:"debug"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Providers defined providers for redis client, it can be redefined
var Providers = ioc.NewProviders(func(name string, args ...any) *ioc.Provider[*redis.Client] {
	return ioc.NewProvider(func(c *ioc.Container) (*redis.Client, error) {

		vp, err := viper.GetViper(name, initConfig, c)
		if err != nil {
			return nil, err
		}

		var cfg Config
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}

		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			Username: cfg.Username,
			DB:       cfg.DB,
		})

		err = client.Ping(context.Background()).Err()
		if err != nil {
			return nil, fmt.Errorf("failed to ping redis: %w", err)
		}

		if cfg.Debug {
			client.AddHook(redisHook{
				logger: log.New(log.Writer(), "redis: ", log.LstdFlags),
			})
		}
		return client, nil
	})
})

func GetClient(name string, c *ioc.Container) (*redis.Client, error) {
	return Providers.GetProvider(name).Get(c)
}

func GetDefaultClient(c *ioc.Container) (*redis.Client, error) {
	return GetClient(DefaultConfig, c)
}

type redisHook struct {
	logger *log.Logger
}

func (h redisHook) printf(format string, v ...interface{}) {
	h.logger.Printf(format, v...)
}

func (h redisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		h.logger.Printf("Dialing to %s:%s\n", network, addr)
		return next(ctx, network, addr)
	}
}

func (h redisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		err := next(ctx, cmd)
		h.printf("command: %s, err: %v\n", cmd, err)
		return err
	}
}

func (h redisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		err := next(ctx, cmds)
		h.printf("pipeline: %v, err: %v\n", cmds, err)
		return err
	}
}
