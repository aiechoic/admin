package verify

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type CodeCache struct {
	expires  time.Duration
	storeKey string
	rds      *redis.Client
}

func NewVerifyCodeCache(rds *redis.Client, storeKey string, expires time.Duration) *CodeCache {
	return &CodeCache{
		expires:  expires,
		storeKey: storeKey,
		rds:      rds,
	}
}

func (vcc *CodeCache) key(id string) string {
	return vcc.storeKey + ":" + id
}

func (vcc *CodeCache) SetVerifyCode(id, code string, checkExpire bool) error {
	expireIn, err := vcc.GetVerifyCodeExpire(id)
	if err != nil {
		return err
	}
	if checkExpire && expireIn > 0 {
		return errors.New("verify code is not expired")
	}
	if expireIn <= 0 {
		expireIn = vcc.expires
	}
	return vcc.rds.Set(context.Background(), vcc.key(id), code, expireIn).Err()
}

func (vcc *CodeCache) GetVerifyCode(id string) (string, error) {
	str, err := vcc.rds.Get(context.Background(), vcc.key(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return "", err
	}
	return str, nil
}

func (vcc *CodeCache) DelVerifyCode(id string) error {
	return vcc.rds.Del(context.Background(), vcc.key(id)).Err()
}

func (vcc *CodeCache) GetVerifyCodeExpire(id string) (time.Duration, error) {
	t, err := vcc.rds.TTL(context.Background(), vcc.key(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return 0, err
	}
	return t, nil
}
