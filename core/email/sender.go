package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
	"html/template"
	"net/smtp"
)

const DefaultSenderConfig = "email-sender"

var senderInitConfig = `
# email sender config

# title of the email
title: "Your Company"

# email address of the sender
from: "sender@example.com"

# smtp host
host: "smtp.example.com"

# smtp port
port: "465"

# smtp password
password: "password"

# email template file
template: "templates/email/verify.gohtml"

# email subject
subject: "Email Verification"
`

var SenderProviders = ioc.NewProviders[*Sender](func(name string, args ...any) *ioc.Provider[*Sender] {
	return ioc.NewProvider(func(c *ioc.Container) (*Sender, error) {

		vp, err := viper.GetViper(name, senderInitConfig, c)
		if err != nil {
			return nil, err
		}

		var cfg emailSenderConfig
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}

		var tpl *template.Template
		if cfg.Template != "" {
			tpl, err = template.ParseFiles(cfg.Template)
			if err != nil {
				return nil, fmt.Errorf("failed to parse email template: %w", err)
			}
		}

		return &Sender{
			opts: &cfg,
			tpl:  tpl,
		}, nil
	})
})

func GetSender(name string, c *ioc.Container) (*Sender, error) {
	return SenderProviders.GetProvider(name).Get(c)
}

func GetDefaultSender(c *ioc.Container) (*Sender, error) {
	return GetSender(DefaultSenderConfig, c)
}

type emailSenderConfig struct {
	Title    string `mapstructure:"title"`
	From     string `mapstructure:"from"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Template string `mapstructure:"template"`
	Subject  string `mapstructure:"subject"`
}

type Sender struct {
	opts *emailSenderConfig
	tpl  *template.Template
}

func (s *Sender) Send(to string, data map[string]interface{}) error {
	buf := bytes.NewBuffer(nil)
	err := s.tpl.Execute(buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	m := "From: " + s.opts.Title + " <" + s.opts.From + ">\n" +
		"To: " + to + "\n" +
		"Subject: " + s.opts.Subject + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		buf.String()

	auth := smtp.PlainAuth("", s.opts.From, s.opts.Password, s.opts.Host)

	// 设置 TLS 配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.opts.Host,
	}

	conn, err := tls.Dial("tcp", s.opts.Host+":"+s.opts.Port, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	client, err := smtp.NewClient(conn, s.opts.Host)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// 验证身份
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to auth: %w", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(s.opts.From); err != nil {
		return fmt.Errorf("failed to set mail: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set rcpt: %w", err)
	}

	// 获取写入邮件数据的写入器
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get writer: %w", err)
	}

	_, err = writer.Write([]byte(m))
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	err = client.Quit()
	if err != nil {
		return fmt.Errorf("failed to quit: %w", err)
	}
	return nil
}
