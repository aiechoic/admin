package gins

import (
	"context"
	"errors"
	"fmt"
	"github.com/aiechoic/admin/core/openapi"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	API       *openapi.Openapi
	Port      int
	Engine    *gin.Engine
	APIRouter gin.IRouter
}

func (s *Server) Run(ctx context.Context) {

	address := fmt.Sprintf(":%d", s.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: s.Engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("received system signal, shutting down gracefully")
	case <-ctx.Done():
		log.Println("context cancelled, shutting down gracefully")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
}

func (s *Server) Register(services ...*Service) {
	for _, service := range services {
		s.register(service)
	}
}

func (s *Server) SetSecuritySchemes(schemes openapi.SecuritySchemes) {
	if s.API.Components.SecuritySchemes == nil {
		s.API.Components.SecuritySchemes = make(openapi.SecuritySchemes)
	}
	for name, scheme := range schemes {
		if _, ok := s.API.Components.SecuritySchemes[name]; ok {
			panic(fmt.Sprintf("security scheme %s already exists", name))
		}
		s.API.Components.SecuritySchemes[name] = scheme
	}
}

func (s *Server) register(service *Service) {
	o := s.API
	r := s.APIRouter
	o.Tags = append(o.Tags, &openapi.Tag{
		Name:        service.Tag,
		Description: service.Description,
	})
	for _, route := range service.Routes {
		// openapi spec requires method to be lowercase
		route.Method = strings.ToLower(route.Method)
		var fullPath string
		if strings.HasPrefix(route.Path, "/") {
			fullPath = route.Path
		} else {
			if route.Path != "" {
				fullPath = service.Path + "/" + route.Path
			} else {
				fullPath = service.Path
			}
		}
		uniqueKey := fmt.Sprintf("%s-%s", strings.Replace(strings.TrimPrefix(fullPath, "/"), "/", "-", -1), route.Method)
		var parameters []*openapi.Parameter
		// Add header parameters to the operation
		if route.Handler.Request.Headers != nil {
			schema, refs := openapi.NewSchema(route.Handler.Request.Headers, uniqueKey, "header", route.Handler.Request.OmitFields)
			for name, prop := range schema.Properties {
				parameters = append(parameters, &openapi.Parameter{
					Name:        name,
					In:          "header",
					Schema:      prop,
					Description: prop.Description,
					Required:    slices.Contains(schema.Required, name),
				})
			}
			o.AddComponentsSchemas(refs)
		}

		// Add query parameters to the operation
		if route.Handler.Request.Query != nil {
			schema, refs := openapi.NewSchema(route.Handler.Request.Query, uniqueKey, "form", route.Handler.Request.OmitFields)
			for name, prop := range schema.Properties {
				parameters = append(parameters, &openapi.Parameter{
					Name:        name,
					In:          "query",
					Schema:      prop,
					Description: prop.Description,
					Required:    slices.Contains(schema.Required, name),
				})
			}
			o.AddComponentsSchemas(refs)
		}
		// Add path parameters to the operation
		pathParts := strings.Split(fullPath, "/")
		var hasPathParams bool
		for i, part := range pathParts {
			if strings.HasPrefix(part, ":") {
				paramName := part[1:]
				parameters = append(parameters, &openapi.Parameter{
					Name:     paramName,
					In:       "path",
					Required: true,
					Schema: &openapi.Schema{
						Type: "string",
					},
				})
				part = fmt.Sprintf("{%s}", paramName)
				hasPathParams = true
			}
			pathParts[i] = part
		}
		var swaggerPath string
		if hasPathParams {
			swaggerPath = strings.Join(pathParts, "/")
		} else {
			swaggerPath = fullPath
		}
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
		content, refs := route.Handler.Response.getContents(uniqueKey)
		o.AddComponentsSchemas(refs)
		op := &openapi.Operation{
			Tags:        []string{service.Tag},
			Summary:     route.Summary,
			Description: route.Description,
			Deprecated:  route.Deprecated,
			Parameters:  parameters,
			Responses: map[openapi.ResponseCode]*openapi.ResponseBody{
				"200": {
					Content:     content,
					Description: route.Handler.Response.Description,
				},
			},
		}
		content, refs = route.Handler.Request.getContents(uniqueKey)
		o.AddComponentsSchemas(refs)
		if route.Handler.Request.Json != nil || route.Handler.Request.Form != nil {
			if route.Handler.Request.Json != nil && route.Handler.Request.Form != nil {
				panic(fmt.Sprintf(
					"service %s route %s: cannot have both json and form Body parameters",
					service.Tag, route.Path,
				))
			}
			if route.Method == "get" {
				panic(fmt.Sprintf(
					"service %s route %s: request json/form Body not allowed for GET method",
					service.Tag, route.Path,
				))
			}
			op.RequestBody = &openapi.RequestBody{
				Content:     content,
				Description: route.Handler.Request.Description,
			}
		}
		pathItem[route.Method] = op
		var handlers []gin.HandlerFunc
		if route.Security == nil && service.Security != nil {
			route.Security = service.Security
		}
		if route.Security != nil && route.Security != NoSecurity {
			op.Security = append(op.Security, route.Security.SecurityScheme())
			handlers = append(handlers, route.Security.Auth)
		}
		handlers = append(handlers, route.Handler.Handle)
		// gin router requires method to be uppercase
		r.Handle(strings.ToUpper(route.Method), fullPath, handlers...)
	}
}
