package gins

import (
	"context"
	"fmt"
	engin "github.com/aiechoic/admin/core/gin"
	"github.com/aiechoic/admin/core/openapi"
	"github.com/gin-gonic/gin"
	"strings"
)

type APIServer struct {
	API   *openapi.Openapi
	Engin *engin.Server
}

func (s *APIServer) Register(services ...*Service) {
	for _, service := range services {
		s.register(service)
	}
}

// getSwaggerPath returns the swagger path for a given gin path
// e.g. /api/v1/user/:id -> /api/v1/user/{id}
func getSwaggerPath(ginPath string) string {
	pathParts := strings.Split(ginPath, "/")
	var hasPathParams bool
	for i, part := range pathParts {
		if strings.HasPrefix(part, ":") {
			paramName := part[1:]
			part = fmt.Sprintf("{%s}", paramName)
			hasPathParams = true
		}
		pathParts[i] = part
	}
	if hasPathParams {
		return strings.Join(pathParts, "/")
	} else {
		return ginPath
	}
}

// getRouteFullPath returns the full path for a route, if the route path is relative(not starting with "/")
// it will be appended to the base path, otherwise it will return the route path as is
// for example:
//
//	/foo, /bar -> /bar
//	/foo, bar -> /foo/bar
func getRouteFullPath(basePath, routePath string) string {
	if strings.HasPrefix(routePath, "/") {
		return routePath
	} else {
		if routePath == "" {
			return basePath
		} else {
			return basePath + "/" + routePath
		}
	}
}

func (s *APIServer) register(service *Service) {
	o := s.API
	r := s.Engin.ApiRouter
	o.Tags = append(o.Tags, &openapi.Tag{
		Name:        service.Tag,
		Description: service.Description,
	})
	for _, route := range service.Routes {
		// openapi spec requires method to be lowercase
		route.Method = strings.ToLower(route.Method)

		// get full path
		fullPath := getRouteFullPath(service.Path, route.Path)

		// get swagger path
		swaggerPath := getSwaggerPath(fullPath)
		pathItem, ok := o.Paths[swaggerPath]
		if !ok {
			pathItem = make(openapi.PathItem)
			o.Paths[swaggerPath] = pathItem
		} else {
			if _, ok = pathItem[route.Method]; ok {
				panic(fmt.Sprintf(
					"service %s route %s: method %s already registered",
					service.Tag, route.Path, route.Method,
				))
			}
		}

		uniqueKey := fmt.Sprintf("%s-%s", strings.Replace(strings.TrimPrefix(fullPath, "/"), "/", "-", -1), route.Method)

		// get request parameters and request body contents
		parameters, refs := route.Handler.Request.getParameters(uniqueKey)
		o.AddComponentsSchemas(refs)
		requestContent, refs := route.Handler.Request.getBodyContents(uniqueKey)
		o.AddComponentsSchemas(refs)

		// get response contents
		responseContent, refs := route.Handler.Response.getBodyContents(uniqueKey)
		o.AddComponentsSchemas(refs)

		// create operation
		op := &openapi.Operation{
			Tags:        []string{service.Tag},
			Summary:     route.Summary,
			Description: route.Description,
			Deprecated:  route.Deprecated,
			Parameters:  parameters,
			RequestBody: &openapi.RequestBody{
				Content:     requestContent,
				Description: route.Handler.Request.Description,
			},
			Responses: map[openapi.ResponseCode]*openapi.ResponseBody{
				"200": {
					Content:     responseContent,
					Description: route.Handler.Response.Description,
				},
			},
		}
		pathItem[route.Method] = op

		// add security
		var handlers []gin.HandlerFunc
		if route.Security == nil && service.Security != nil {
			route.Security = service.Security
		}
		if route.Security != nil && route.Security != NoSecurity {
			securityRequirement := route.Security.SecurityRequirement()
			op.Security = append(op.Security, securityRequirement)
			securitySchemes := route.Security.SecuritySchemes()
			for name, scheme := range securitySchemes {
				if _, ok = o.Components.SecuritySchemes[name]; !ok {
					o.Components.SecuritySchemes[name] = scheme
				}
			}
			handlers = append(handlers, route.Security.Auth)
		}
		handlers = append(handlers, route.Handler.Handle)

		// register route with gin
		r.Handle(strings.ToUpper(route.Method), fullPath, handlers...)
	}
}

func (s *APIServer) Run(ctx context.Context) {
	s.Engin.Run(ctx)
}
