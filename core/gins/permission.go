package gins

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"runtime"
)

var authHandlerPermissions = map[uintptr]*Permission{}

var sortHandlerPermissions = map[string][]*Permission{}

type Permission struct {
	Tag    string
	Method string
	Path   string
	Code   string // hash of the method and path
}

func setHandlerPermission(fn gin.HandlerFunc, p *Permission) {
	pointer := reflect.ValueOf(fn).Pointer()
	if _, ok := authHandlerPermissions[pointer]; ok {
		panic(fmt.Sprintf("handler %s already exist", runtime.FuncForPC(pointer).Name()))
	}
	authHandlerPermissions[pointer] = p
	sortHandlerPermissions[p.Tag] = append(sortHandlerPermissions[p.Tag], p)
}

func GetHandlerPermission(c *gin.Context) *Permission {
	return authHandlerPermissions[reflect.ValueOf(c.Handler()).Pointer()]
}

func GetAllPermissions() map[string][]*Permission {
	return sortHandlerPermissions
}

func getStringHash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
