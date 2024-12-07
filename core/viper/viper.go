package viper

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

const (
	Development Env = "development"
	Testing     Env = "testing"
	Production  Env = "production"
)

// adapterProvider is a provider for Adapter, which is a factory for creating Viper instances.
var adapterProvider = ioc.NewProvider(func(c *ioc.Container) (Adapter, error) {
	return NewLocalAdapter("configs", Testing), nil
})

// SetAdapter is a helper function to set the ViperAdapterProvider.
func SetAdapter(adapter Adapter) {
	adapterProvider = ioc.NewProvider(func(c *ioc.Container) (Adapter, error) {
		return adapter, nil
	})
}

// GetViper creates a new Viper instance with the given name and initial configuration.
func GetViper(name string, initConfig string, c *ioc.Container) (*viper.Viper, error) {
	viperAdapter, err := adapterProvider.Get(c)
	if err != nil {
		return nil, err
	}
	contentType := detectContentType([]byte(initConfig))
	if contentType == "" {
		return nil, fmt.Errorf("parameter 'initConfig' should be JSON, YAML, or TOML format")
	}
	return viperAdapter.NewViper(name, contentType, initConfig)
}

type Env string

// Adapter is a factory for creating Viper instances.
type Adapter interface {
	// NewViper creates a new Viper instance with the given name and initial configuration.
	NewViper(name, contentType string, initConfig string) (*viper.Viper, error)
}

// LocalAdapter is a Adapter implementation that stores configuration files locally.
type LocalAdapter struct {
	dir string
	env Env
}

func NewLocalAdapter(dir string, env Env) *LocalAdapter {
	envDir := filepath.Join(dir, string(env))
	_, err := os.Stat(envDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(envDir, 0755)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return &LocalAdapter{
		dir: dir,
		env: env,
	}
}

// NewViper implements the Adapter interface. If the configuration file does not exist, it creates a new
// file with the initial configuration. It automatically read environment variables that match the prefix and
// replace the separator with an underscore. For example, if the name is "app" and the separator is ".", the
// environment variable "app.port" will be read as "app_port".
func (l *LocalAdapter) NewViper(name, ext, initConfig string) (*viper.Viper, error) {
	initConfigData := []byte(initConfig)
	filename := filepath.Join(l.dir, string(l.env), fmt.Sprintf("%s.%s", name, ext))
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.WriteFile(filename, initConfigData, 0644)
			if err != nil {
				return nil, err
			}
			logrus.Infof("created default config file '%s'", filename)
		} else {
			return nil, err
		}
	}
	v := viper.New()
	// load environment variables
	v.SetConfigFile(filename)
	err = v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("reading config file \"%s\": %w", filename, err)
	}
	v.SetEnvPrefix(name)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	return v, nil
}

func detectContentType(content []byte) string {
	dt := map[string]interface{}{}
	err := json.Unmarshal(content, &dt)
	if err == nil {
		return "json"
	}
	err = yaml.Unmarshal(content, &dt)
	if err == nil {
		return "yaml"
	}
	err = toml.Unmarshal(content, &dt)
	if err == nil {
		return "toml"
	}
	return ""
}
