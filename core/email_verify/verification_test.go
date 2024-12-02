package email_verify_test

import (
	"github.com/aiechoic/admin/core/email_verify"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockEmailSender struct {
	sendFunc func(email string, data map[string]interface{}) error
}

func (m *MockEmailSender) Send(email string, data map[string]interface{}) error {
	return m.sendFunc(email, data)
}

var dataSetter = func(code string, expire time.Duration) map[string]interface{} {
	return map[string]interface{}{
		"code":           code,
		"expireInMinute": int64(expire / time.Minute),
	}
}

func setupEmailVerification(t *testing.T, expire time.Duration) (*email_verify.EmailVerification, *MockEmailSender) {
	cache := setupCache(t, expire)
	mockSender := &MockEmailSender{
		sendFunc: func(email string, data map[string]interface{}) error {
			return nil
		},
	}
	ev := email_verify.NewEmailVerification(cache, mockSender, "0123456789", 6)
	return ev, mockSender
}

func TestEmailVerification_SendVerificationCode(t *testing.T) {
	ev, mockSender := setupEmailVerification(t, 5*time.Minute)

	mockSender.sendFunc = func(email string, data map[string]interface{}) error {
		assert.Equal(t, "test@example.com", email)
		assert.NotEmpty(t, data["code"])
		assert.Equal(t, int64(5), data["expireInMinute"])
		return nil
	}

	err := ev.SendVerificationCode("test@example.com", dataSetter)
	assert.NoError(t, err)
}

func TestEmailVerification_GetWaitTime(t *testing.T) {
	expires := 5 * time.Minute
	ev, _ := setupEmailVerification(t, expires)

	err := ev.SendVerificationCode("test@example.com", dataSetter)
	assert.NoError(t, err)

	seconds, err := ev.GetWaitTime("test@example.com")
	assert.NoError(t, err)
	expireSeconds := int64(expires / time.Second)
	assert.True(t, expireSeconds-1 <= seconds && seconds <= expireSeconds)
}

func TestEmailVerification_VerifyCode(t *testing.T) {
	expires := 5 * time.Minute
	ev, mockSender := setupEmailVerification(t, expires)

	var code string
	mockSender.sendFunc = func(email string, data map[string]interface{}) error {
		code = data["code"].(string)
		return nil
	}
	err := ev.SendVerificationCode("test@example.com", dataSetter)
	assert.NoError(t, err)

	ok, err := ev.VerifyCode("test@example.com", code)
	assert.NoError(t, err)
	assert.True(t, ok)

	err = ev.SendVerificationCode("test@example.com", dataSetter)
	assert.NoError(t, err)

	var sleepSeconds int64 = 2
	time.Sleep(time.Duration(sleepSeconds) * time.Second)
	ok, err = ev.VerifyCode("test@example.com", "wrongcode")
	assert.NoError(t, err)
	assert.False(t, ok)

	seconds, err := ev.GetWaitTime("test@example.com")
	assert.NoError(t, err)
	expireSeconds := int64(expires / time.Second)
	assert.True(t, expireSeconds-1-sleepSeconds <= seconds && seconds <= expireSeconds-sleepSeconds)
}
