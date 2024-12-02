package email_verify_test

import (
	"context"
	"github.com/aiechoic/admin/internal/email_verify"
	"github.com/aiechoic/admin/internal/ioc"
	"github.com/aiechoic/admin/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 初始化 Cache
func setupCache(t *testing.T, expire time.Duration) *email_verify.VerifyCodeCache {
	c := ioc.NewContainer()
	rdb, err := service.GetDefaultRedisClient(c)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		rdb.FlushDB(context.Background())
		_ = rdb.Close()
	})
	return email_verify.NewVerifyCodeCache(rdb, "email_ver_code", expire)
}

func TestVerifyCodeCache(t *testing.T) {

	vcc := setupCache(t, 5*time.Second)

	t.Run("SetEmailVerifyCode", func(t *testing.T) {
		err := vcc.SetEmailVerifyCode("test@example.com", "123456", false)
		assert.NoError(t, err)

		err = vcc.SetEmailVerifyCode("test@example.com", "123456", true)
		assert.Error(t, err)
	})

	t.Run("GetEmailVerifyCode", func(t *testing.T) {
		err := vcc.SetEmailVerifyCode("test@example.com", "123456", false)
		assert.NoError(t, err)
		code, err := vcc.GetEmailVerifyCode("test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "123456", code)
	})

	t.Run("DelEmailVerifyCode", func(t *testing.T) {
		err := vcc.SetEmailVerifyCode("test@example.com", "123456", false)
		assert.NoError(t, err)
		err = vcc.DelEmailVerifyCode("test@example.com")
		assert.NoError(t, err)

		code, err := vcc.GetEmailVerifyCode("test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "", code)
	})

	t.Run("GetEmailVerifyCodeExpire", func(t *testing.T) {
		err := vcc.SetEmailVerifyCode("test@example.com", "123456", false)
		assert.NoError(t, err)
		expire, err := vcc.GetEmailVerifyCodeExpire("test@example.com")
		assert.NoError(t, err)
		assert.True(t, 3*time.Second <= expire && expire <= 5*time.Second)
	})
}
