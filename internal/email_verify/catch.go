package email_verify

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type VerifyCodeCache struct {
	expires         time.Duration
	emailVerCodeKey string
	rds             *redis.Client
}

func NewVerifyCodeCache(rds *redis.Client, emailVerCodeKey string, expires time.Duration) *VerifyCodeCache {
	return &VerifyCodeCache{
		expires:         expires,
		emailVerCodeKey: emailVerCodeKey,
		rds:             rds,
	}
}

func (vcc *VerifyCodeCache) key(email string) string {
	return vcc.emailVerCodeKey + ":" + email
}

func (vcc *VerifyCodeCache) SetEmailVerifyCode(email, code string, checkExpire bool) error {
	expireIn, err := vcc.GetEmailVerifyCodeExpire(email)
	if err != nil {
		return err
	}
	if checkExpire && expireIn > 0 {
		return errors.New("verify code is not expired")
	}
	if expireIn <= 0 {
		expireIn = vcc.expires
	}
	return vcc.rds.Set(context.Background(), vcc.key(email), code, expireIn).Err()
}

func (vcc *VerifyCodeCache) GetEmailVerifyCode(email string) (string, error) {
	str, err := vcc.rds.Get(context.Background(), vcc.key(email)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return "", err
	}
	return str, nil
}

func (vcc *VerifyCodeCache) DelEmailVerifyCode(email string) error {
	return vcc.rds.Del(context.Background(), vcc.key(email)).Err()
}

func (vcc *VerifyCodeCache) GetEmailVerifyCodeExpire(email string) (time.Duration, error) {
	t, err := vcc.rds.TTL(context.Background(), vcc.key(email)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return 0, err
	}
	return t, nil
}
