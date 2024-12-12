package viper_test

import (
	"github.com/aiechoic/admin/core/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalViperAdapter_NewViper(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "test_configs")
	env := viper.Testing
	adapter := viper.NewLocalAdapter(dir, env)

	// Clean up test directory
	defer os.RemoveAll(dir)

	// Test creating a new Viper instance with initial configuration
	name := "test"
	initConfig := `key: "value"`
	v, err := adapter.NewViper(name, "yaml", initConfig)
	assert.NoError(t, err)

	// Verify the configuration value
	assert.Equal(t, "value", v.GetString("key"))

	// Test reading from existing configuration file
	v, err = adapter.NewViper(name, "yaml", `key: "new_value"`)
	assert.NoError(t, err)

	// Verify the configuration value remains the same
	assert.Equal(t, "value", v.GetString("key"))

	// Test reading from system environment variables
	os.Setenv("TEST_APP_PORT", "env_value")
	defer os.Unsetenv("APP_PORT")

	v, err = adapter.NewViper("test_app", "yaml", `port: "value"`)
	assert.NoError(t, err)

	type config struct {
		Port string `mapstructure:"port"`
	}
	var c config
	err = v.Unmarshal(&c)
	assert.NoError(t, err)
	assert.Equal(t, "env_value", c.Port)
}
