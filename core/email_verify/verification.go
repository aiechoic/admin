package email_verify

import (
	"github.com/aiechoic/admin/core/random"
	"time"
)

type EmailSender interface {
	Send(email string, data map[string]interface{}) error
}

type EmailVerification struct {
	cache        *VerifyCodeCache
	sender       EmailSender
	randomCharts []rune
	codeLength   int
}

func NewEmailVerification(
	cache *VerifyCodeCache,
	sender EmailSender,
	randomCharts string,
	codeLength int,
) *EmailVerification {
	return &EmailVerification{
		cache:        cache,
		sender:       sender,
		randomCharts: []rune(randomCharts),
		codeLength:   codeLength,
	}
}

type DataSetter func(code string, expire time.Duration) map[string]interface{}

func (ev *EmailVerification) SendVerificationCode(email string, ds DataSetter) error {

	code := random.StringWithCharset(ev.randomCharts, ev.codeLength)

	err := ev.cache.SetEmailVerifyCode(email, code, true)
	if err != nil {
		return err
	}

	// Send email
	data := ds(code, ev.cache.expires)
	if err = ev.sender.Send(email, data); err != nil {
		return err
	}
	return nil
}

func (ev *EmailVerification) GetWaitTime(email string) (seconds int64, err error) {
	ex, err := ev.cache.GetEmailVerifyCodeExpire(email)
	if err != nil {
		return 0, err
	}
	return int64(ex / time.Second), nil
}

func (ev *EmailVerification) VerifyCode(email, code string) (ok bool, err error) {
	sCode, err := ev.cache.GetEmailVerifyCode(email)
	if err != nil {
		return false, err
	}
	if sCode != "" && sCode == code {
		// delete the code after verification
		err = ev.cache.DelEmailVerifyCode(email)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		// set empty code to prevent brute force attack
		err = ev.cache.SetEmailVerifyCode(email, "", false)
		return false, err
	}
}
