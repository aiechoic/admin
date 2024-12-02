package gins

import (
	"github.com/aiechoic/admin/core/openapi"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

const DefaultConfig = "gin-server"

var InitConfig = `
# Gin Server Config

# title for openapi
api_title: "API Documentation"

# version of the api
api_version: "1.0.0"

# http port to listen on
http_port: 8080

# api root path
api_root: "/api/v1"

# api servers
api_servers:
  - url: "http://localhost:8080/api/v1"
    description: "Local Server"

# gin mode, can be "debug", "release", "test"
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

type openAPIServer struct {
	// on production, this should be the actual url, e.g. https://api.example.com/api/v1
	Url         string `mapstructure:"url"`
	Description string `mapstructure:"description"`
}

type Config struct {
	ApiTitle             string          `mapstructure:"api_title"`   // for openapi
	ApiVersion           string          `mapstructure:"api_version"` // for openapi
	ApiServers           []openAPIServer `mapstructure:"api_servers"` // for openapi
	HttpPort             int             `mapstructure:"http_port"`
	APIRoot              string          `mapstructure:"api_root"`
	GinMode              string          `mapstructure:"gin_mode"`
	EnableLogger         bool            `mapstructure:"enable_logger"`
	EnableRecovery       bool            `mapstructure:"enable_recovery"`
	EnableCORS           bool            `mapstructure:"enable_cors"`
	CorsAllowMethods     []string        `mapstructure:"cors_allow_methods"`
	CorsAllowHeaders     []string        `mapstructure:"cors_allow_headers"`
	CorsAllowCredentials bool            `mapstructure:"cors_allow_credentials"`
	CorsAllowOrigins     []string        `mapstructure:"cors_allow_origins"`
	CorsMaxAge           time.Duration   `mapstructure:"cors_max_age"`
}

func (g *Config) NewServer(defaultEngine *gin.Engine) *Server {
	engine, router := g.NewEngine(defaultEngine)
	api := g.NewOpenAPI()
	s := &Server{
		API:       api,
		Port:      g.HttpPort,
		Engine:    engine,
		APIRouter: router,
	}
	return s
}

func (g *Config) NewEngine(defaultEngine *gin.Engine) (*gin.Engine, gin.IRouter) {
	if g.GinMode != "" {
		gin.SetMode(g.GinMode)
	}
	var r *gin.Engine
	if defaultEngine != nil {
		r = defaultEngine
	} else {
		r = gin.New()
	}
	if g.EnableLogger {
		r.Use(gin.Logger())
	}
	if g.EnableRecovery {
		r.Use(gin.Recovery())
	}
	if g.EnableCORS {
		r.Use(cors.New(cors.Config{
			AllowMethods:     g.CorsAllowMethods,
			AllowHeaders:     g.CorsAllowHeaders,
			AllowCredentials: g.CorsAllowCredentials,
			AllowOrigins:     g.CorsAllowOrigins,
			MaxAge:           g.CorsMaxAge,
		}))
	}
	var i gin.IRouter = r
	if g.APIRoot != "" {
		i = r.Group(g.APIRoot)
	}
	return r, i
}

func (g *Config) NewOpenAPI() *openapi.Openapi {
	var servers []*openapi.Server
	for _, s := range g.ApiServers {
		servers = append(servers, &openapi.Server{
			Url:         s.Url,
			Description: s.Description,
		})
	}
	info := &openapi.Info{
		Title:   g.ApiTitle,
		Version: g.ApiVersion,
	}
	return &openapi.Openapi{
		Openapi: "3.0.0",
		Info:    info,
		Servers: servers,
		Components: openapi.Components{
			SecuritySchemes: map[string]*openapi.SecurityScheme{},
			Schemas:         map[string]*openapi.Schema{},
		},
		Paths: map[string]openapi.PathItem{},
	}
}
