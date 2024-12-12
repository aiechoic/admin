package verify

import (
	"github.com/aiechoic/admin/core/random"
	"time"
)

type Sender interface {
	Send(id string, data map[string]interface{}) error
}

type Verification struct {
	cache        *CodeCache
	sender       Sender
	randomCharts []rune
	codeLength   int
}

func NewVerification(
	cache *CodeCache,
	sender Sender,
	randomCharts string,
	codeLength int,
) *Verification {
	return &Verification{
		cache:        cache,
		sender:       sender,
		randomCharts: []rune(randomCharts),
		codeLength:   codeLength,
	}
}

type DataSetter func(code string, expire time.Duration) map[string]interface{}

func (ev *Verification) SendCode(id string, ds DataSetter) error {

	code := random.StringWithCharset(ev.randomCharts, ev.codeLength)

	err := ev.cache.SetVerifyCode(id, code, true)
	if err != nil {
		return err
	}

	// Send email
	data := ds(code, ev.cache.expires)
	if err = ev.sender.Send(id, data); err != nil {
		return err
	}
	return nil
}

func (ev *Verification) GetWaitTime(id string) (seconds int64, err error) {
	ex, err := ev.cache.GetVerifyCodeExpire(id)
	if err != nil {
		return 0, err
	}
	return int64(ex / time.Second), nil
}

func (ev *Verification) VerifyCode(id, code string) (ok bool, err error) {
	sCode, err := ev.cache.GetVerifyCode(id)
	if err != nil {
		return false, err
	}
	if sCode != "" && sCode == code {
		// delete the code after verification
		err = ev.cache.DelVerifyCode(id)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		// set empty code to prevent brute force attack
		err = ev.cache.SetVerifyCode(id, "", false)
		return false, err
	}
}
