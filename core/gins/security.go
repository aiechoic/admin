package gins

import (
	"github.com/aiechoic/admin/core/openapi"
	"github.com/gin-gonic/gin"
)

var NoSecurity = &noSecurity{}

type noSecurity struct{}

func (n *noSecurity) Auth(*gin.Context) {}

func (n *noSecurity) SecuritySchemes() openapi.SecuritySchemes { return nil }

func (n *noSecurity) SecurityRequirement() map[string][]string { return nil }

type Security interface {
	Auth(c *gin.Context)
	SecuritySchemes() openapi.SecuritySchemes
	SecurityRequirement() map[string][]string
}
