package gin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

const DefaultConfig = "gin-server"

var initConfig = `
# Gin Server Config

# http port to listen on
http_port: 8080

# api root path
api_root: "/api/v1"

# gin mode, can be "debug", "release", "test" or "" for default 
gin_mode: "debug"

# enable default logger middleware
enable_logger: true

# enable default recovery middleware
enable_recovery: true

# enable cross-origin resource sharing
enable_cors: true

# cors configurations, only used when enable_cors is true
cors_allow_methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"]
cors_allow_headers: ["Origin", "Content-Length", "Content-Type", "Authorization"]
cors_allow_credentials: false
cors_allow_origins: ["*"]
cors_max_age: "12h"

`

type Config struct {
	ApiRoot              string        `mapstructure:"api_root"`
	HttpPort             int           `mapstructure:"http_port"`
	GinMode              string        `mapstructure:"gin_mode"`
	EnableLogger         bool          `mapstructure:"enable_logger"`
	EnableRecovery       bool          `mapstructure:"enable_recovery"`
	EnableCORS           bool          `mapstructure:"enable_cors"`
	CorsAllowMethods     []string      `mapstructure:"cors_allow_methods"`
	CorsAllowHeaders     []string      `mapstructure:"cors_allow_headers"`
	CorsAllowCredentials bool          `mapstructure:"cors_allow_credentials"`
	CorsAllowOrigins     []string      `mapstructure:"cors_allow_origins"`
	CorsMaxAge           time.Duration `mapstructure:"cors_max_age"`
}

func (c *Config) NewGinEngine() (*gin.Engine, gin.IRouter) {
	if c.GinMode != "" {
		gin.SetMode(c.GinMode)
	}
	r := gin.New()
	if c.EnableLogger {
		r.Use(gin.Logger())
	}
	if c.EnableRecovery {
		r.Use(gin.Recovery())
	}
	if c.EnableCORS {
		r.Use(cors.New(cors.Config{
			AllowMethods:     c.CorsAllowMethods,
			AllowHeaders:     c.CorsAllowHeaders,
			AllowCredentials: c.CorsAllowCredentials,
			AllowOrigins:     c.CorsAllowOrigins,
			MaxAge:           c.CorsMaxAge,
		}))
	}
	var i gin.IRouter = r
	if c.ApiRoot != "" {
		i = r.Group(c.ApiRoot)
	}
	return r, i
}
