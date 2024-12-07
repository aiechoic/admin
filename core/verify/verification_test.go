package verify_test

import (
	"github.com/aiechoic/admin/core/verify"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockSender struct {
	sendFunc func(email string, data map[string]interface{}) error
}

func (m *MockSender) Send(email string, data map[string]interface{}) error {
	return m.sendFunc(email, data)
}

var dataSetter = func(code string, expire time.Duration) map[string]interface{} {
	return map[string]interface{}{
		"code":           code,
		"expireInMinute": int64(expire / time.Minute),
	}
}

func setupVerification(t *testing.T, expire time.Duration) (*verify.Verification, *MockSender) {
	cache := setupCache(t, expire)
	mockSender := &MockSender{
		sendFunc: func(email string, data map[string]interface{}) error {
			return nil
		},
	}
	ev := verify.NewVerification(cache, mockSender, "0123456789", 6)
	return ev, mockSender
}

func TestVerification_SendCode(t *testing.T) {
	ev, mockSender := setupVerification(t, 5*time.Minute)

	mockSender.sendFunc = func(email string, data map[string]interface{}) error {
		assert.Equal(t, "test@example.com", email)
		assert.NotEmpty(t, data["code"])
		assert.Equal(t, int64(5), data["expireInMinute"])
		return nil
	}

	err := ev.SendCode("test@example.com", dataSetter)
	assert.NoError(t, err)
}

func TestVerification_GetWaitTime(t *testing.T) {
	expires := 5 * time.Minute
	ev, _ := setupVerification(t, expires)

	err := ev.SendCode("test@example.com", dataSetter)
	assert.NoError(t, err)

	seconds, err := ev.GetWaitTime("test@example.com")
	assert.NoError(t, err)
	expireSeconds := int64(expires / time.Second)
	assert.True(t, expireSeconds-1 <= seconds && seconds <= expireSeconds)
}

func TestVerification_VerifyCode(t *testing.T) {
	expires := 5 * time.Minute
	ev, mockSender := setupVerification(t, expires)

	var code string
	mockSender.sendFunc = func(email string, data map[string]interface{}) error {
		code = data["code"].(string)
		return nil
	}
	err := ev.SendCode("test@example.com", dataSetter)
	assert.NoError(t, err)

	ok, err := ev.VerifyCode("test@example.com", code)
	assert.NoError(t, err)
	assert.True(t, ok)

	err = ev.SendCode("test@example.com", dataSetter)
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
