package verify_test

import (
	"context"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/redis"
	"github.com/aiechoic/admin/core/verify"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 初始化 Cache
func setupCache(t *testing.T, expire time.Duration) *verify.CodeCache {
	c := ioc.NewContainer()
	rdb, err := redis.GetDefaultClient(c)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		rdb.FlushDB(context.Background())
		_ = rdb.Close()
	})
	return verify.NewVerifyCodeCache(rdb, "email_ver_code", expire)
}

func TestVerifyCodeCache(t *testing.T) {

	vcc := setupCache(t, 5*time.Second)

	t.Run("SetVerifyCode", func(t *testing.T) {
		err := vcc.SetVerifyCode("test@example.com", "123456", false)
		assert.NoError(t, err)

		err = vcc.SetVerifyCode("test@example.com", "123456", true)
		assert.Error(t, err)
	})

	t.Run("GetVerifyCode", func(t *testing.T) {
		err := vcc.SetVerifyCode("test@example.com", "123456", false)
		assert.NoError(t, err)
		code, err := vcc.GetVerifyCode("test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "123456", code)
	})

	t.Run("DelVerifyCode", func(t *testing.T) {
		err := vcc.SetVerifyCode("test@example.com", "123456", false)
		assert.NoError(t, err)
		err = vcc.DelVerifyCode("test@example.com")
		assert.NoError(t, err)

		code, err := vcc.GetVerifyCode("test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "", code)
	})

	t.Run("GetVerifyCodeExpire", func(t *testing.T) {
		err := vcc.SetVerifyCode("test@example.com", "123456", false)
		assert.NoError(t, err)
		expire, err := vcc.GetVerifyCodeExpire("test@example.com")
		assert.NoError(t, err)
		assert.True(t, 3*time.Second <= expire && expire <= 5*time.Second)
	})
}
