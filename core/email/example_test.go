package email_test

import (
	"github.com/aiechoic/admin/core/email"
	"github.com/aiechoic/admin/core/ioc"
)

func ExampleGetSender() {
	c := ioc.NewContainer()

	sender, err := email.GetSender(email.DefaultSenderConfig, c)
	if err != nil {
		panic(err)
	}
	err = sender.Send("target@example.com", map[string]interface{}{
		"verifyCode": "123456",
		"expireIn":   "5min",
	})
	if err != nil {
		panic(err)
	}
}
