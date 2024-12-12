package env

import (
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
)

const DefaultEnvConfig = "env"

var initConfig = `
# Environment configuration

# secret key, should always set in environment variable
# for example: export ENV_SECRET_KEY=your_secret_key 
secret_key: ""
`

type Env struct {
	SecretKey string `mapstructure:"secret_key"`
}

var envProvider = ioc.NewProvider(func(c *ioc.Container) (*Env, error) {
	vp, err := viper.GetViper(DefaultEnvConfig, initConfig, c)
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
