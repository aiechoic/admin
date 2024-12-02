package env

import (
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/service"
)

const DefaultEnvConfig = "env"

var ServerInitConfig = `
# Environment configuration

# version
version: "1.0.0"

# secret key, should always set in environment variable
# for example: export ENV_SECRET_KEY=your_secret_key 
secret_key: ""
`

type Env struct {
	Version   string `mapstructure:"version"`
	SecretKey string `mapstructure:"secret_key"`
}

var envProvider = ioc.NewProvider(func(c *ioc.Container) (*Env, error) {
	vp, err := service.GetViper(DefaultEnvConfig, ServerInitConfig, c)
	if err != nil {
		return nil, err
	}
	var server Env
	err = vp.Unmarshal(&server)
	if err != nil {
		return nil, err
	}

	return &server, nil
})

func GetEnv(c *ioc.Container) (*Env, error) {
	return envProvider.Get(c)
}
