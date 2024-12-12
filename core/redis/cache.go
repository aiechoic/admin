package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
	"github.com/redis/go-redis/v9"
	"strings"
	"sync"
	"time"
)

var initCacheConfig = `
# Cache configuration

# redis key prefix
key: "$name"

# redis key expiration
expiration: "5m"
`

var CacheProviders = ioc.NewProviders(func(name string, args ...string) *ioc.Provider[*CacheConfig] {
	_initCacheConfig := strings.Replace(initCacheConfig, "$name", name, -1)
	return ioc.NewProvider(func(c *ioc.Container) (*CacheConfig, error) {
		vp, err := viper.GetViper(name, _initCacheConfig, c)
		if err != nil {
			return nil, err
		}
		var cfg CacheConfig
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}
		if cfg.Key == "" {
			return nil, fmt.Errorf("config '%s': cache key can not be empty", name)
		}
		return &cfg, nil
	})
})

func GetCache[T any](cacheConfig, redisConfig string, c *ioc.Container) (*Cache[T], error) {
	cfg, err := CacheProviders.GetProvider(cacheConfig).Get(c)
	if err != nil {
		return nil, err
	}
	return newCache[T](cfg, redisConfig, c)
}

func GetDefaultCache[T any](cacheConfig string, c *ioc.Container) (*Cache[T], error) {
	return GetCache[T](cacheConfig, DefaultConfig, c)
}

var caches = sync.Map{}

func newCache[T any](c *CacheConfig, redisConfig string, container *ioc.Container) (*Cache[T], error) {
	cache, ok := caches.Load(c)
	if ok {
		return cache.(*Cache[T]), nil
	}
	rds, err := GetClient(redisConfig, container)
	if err != nil {
		return nil, err
	}
	cache = NewCache[T](rds, c.Key, c.Expiration)
	caches.Store(c, cache)
	return cache.(*Cache[T]), nil
}

type CacheConfig struct {
	Key        string
	Expiration time.Duration
}

type Cache[T any] struct {
	rds        *redis.Client
	key        string
	expiration time.Duration
}

func NewCache[T any](rds *redis.Client, key string, expiration time.Duration) *Cache[T] {
	return &Cache[T]{
		rds:        rds,
		key:        key,
		expiration: expiration,
	}
}

func (c *Cache[T]) Set(id string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rds.Set(context.Background(), c.key+id, data, c.expiration).Err()
}

// Get returns the value of the key. If the key does not exist, nil is returned.
func (c *Cache[T]) Get(id string) (*T, error) {
	key := c.key + id
	data, err := c.rds.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var value T
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (c *Cache[T]) Del(id string) error {
	return c.rds.Del(context.Background(), c.key+id).Err()
}
