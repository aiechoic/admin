package gins

import (
	"github.com/aiechoic/admin/core/openapi"
	"github.com/gin-gonic/gin"
)

var (
	// ContentsTextHtml defines the response body content type is "text/html"
	ContentsTextHtml = map[openapi.ContentType]*openapi.MediaType{
		"text/html": {
			Schema: &openapi.Schema{Type: "string"},
		},
	}

	// ContentsOctetStream defines the response body content type is "application/octet-stream"
	ContentsOctetStream = map[openapi.ContentType]*openapi.MediaType{
		"application/octet-stream": {
			Schema: &openapi.Schema{Type: "string", Format: "binary"},
		},
	}

	// ContentsImages defines the response body content type is "image/*"
	ContentsImages = map[openapi.ContentType]*openapi.MediaType{
		"image/*": {
			Schema: &openapi.Schema{Type: "string", Format: "binary"},
		},
	}
)

var NoSecurity = &noSecurity{}

type noSecurity struct{}

func (n *noSecurity) Auth(c *gin.Context) {}

func (n *noSecurity) SecurityScheme() map[string][]string {
	return nil
}

type Request struct {
	// request summary
	Description string

	// headers parameters
	Headers any

	// url query parameters
	Query any

	// request body encode type: "application/json"
	// example:
	// type User struct {
	// 	 ID       int    `json:"id"`
	//   Name     string `json:"name"`
	// }
	Json any

	// request body encode type: "application/x-www-form-urlencoded" or "multipart/form-data"
	// if the schema has file property, the encode type is "multipart/form-data", otherwise
	// it is "application/x-www-form-urlencoded".
	//
	// example:
	// type Upload struct {
	// 	 File []byte `form:"file" format:"binary" description:"The binary of the upload file"`
	//   Image []byte `form:"image" format:"byte" description:"The base64 of the upload file"`
	//   Avatar *multipart.FileHeader `form:"avatar" description:"The binary of the upload file"`
	// }
	Form any

	// request body encode type: application/xml
	// example:
	// type User struct {
	// 	 ID       int    `xml:"id"`
	//   Name     string `xml:"name"`
	// }
	Xml any

	// Contents defines the request body content type, such as "application/xml", "application/json"
	// "application/x-www-form-urlencoded", "multipart/form-data", etc. most of the time, you don't need
	// to use this field, just use Query, Json, Form, Xml fields is enough, these fields will automatically
	// generate the Contents field and set the corresponding schema.
	// See https://swagger.io/specification/#request-body-object for more information.
	Contents map[openapi.ContentType]*openapi.MediaType

	OmitFields []string
}

func (r *Request) getContents(key string) (map[openapi.ContentType]*openapi.MediaType, map[string]*openapi.Schema) {
	var contents = map[openapi.ContentType]*openapi.MediaType{}
	var refs map[string]*openapi.Schema
	var schema *openapi.Schema
	if r.Json != nil {
		schema, refs = openapi.NewSchema(r.Json, key, "json", r.OmitFields)
		contents[openapi.ContentTypeJson] = &openapi.MediaType{
			Schema: schema,
		}
	} else if r.Form != nil {
		schema, refs = openapi.NewSchema(r.Form, key, "form", r.OmitFields)
		var ct openapi.ContentType
		if hasFileProperty(schema) {
			ct = openapi.ContentTypeMultipartForm
		} else {
			ct = openapi.ContentTypeForm
		}
		contents[ct] = &openapi.MediaType{
			Schema: schema,
		}
	} else if r.Xml != nil {
		schema, refs = openapi.NewSchema(r.Xml, key, "xml", r.OmitFields)
		contents[openapi.ContentTypeXml] = &openapi.MediaType{
			Schema: schema,
		}
	}

	for ct, media := range r.Contents {
		contents[ct] = media
	}
	return contents, refs
}

type Response struct {
	// response summary
	Description string

	// response body encode type: "application/json"
	Json any

	// response body encode type: "application/xml"
	Xml any

	// Contents defines the response body content type, such as "application/octet-stream", "image/*", "text/html", etc.
	// Most of the time, you don't need to use this field, just use Json, Xml fields is enough, these fields will
	// automatically generate the Contents field and set the corresponding schema. There are some common contents
	// defined in this package, such as ContentsTextHtml, ContentsOctetStream, ContentsImages, you can use them directly.
	// See https://swagger.io/specification/#response-object for more information.
	Contents map[openapi.ContentType]*openapi.MediaType

	OmitFields []string
}

func (r *Response) getContents(key string) (map[openapi.ContentType]*openapi.MediaType, map[string]*openapi.Schema) {
	var contents = map[openapi.ContentType]*openapi.MediaType{}
	var schema *openapi.Schema
	var refs map[string]*openapi.Schema
	if r.Json != nil {
		schema, refs = openapi.NewSchema(r.Json, key, "json", r.OmitFields)
		contents[openapi.ContentTypeJson] = &openapi.MediaType{
			Schema: schema,
		}
	} else if r.Xml != nil {
		schema, refs = openapi.NewSchema(r.Xml, key, "xml", r.OmitFields)
		contents[openapi.ContentTypeXml] = &openapi.MediaType{
			Schema: schema,
		}
	}
	for ct, media := range r.Contents {
		contents[ct] = media
	}
	return contents, refs
}

type Route struct {
	Method      string
	Path        string
	Summary     string
	Description string
	Deprecated  bool
	Security    Security
	Handler     Handler
}

type Handler struct {
	Request  Request
	Response Response
	Handle   func(c *gin.Context)
}

type Security interface {
	Auth(c *gin.Context)
	SecurityScheme() map[string][]string
}

func (r Route) Use(h Security) Route {
	r.Security = h
	return r
}

type Service struct {
	Tag         string
	Description string
	Path        string
	Security    Security
	Routes      []Route
}

func hasFileProperty(s *openapi.Schema) bool {
	if s.Type == "object" {
		for _, prop := range s.Properties {
			if hasFileProperty(prop) {
				return true
			}
		}
		return false
	}
	if s.Type == "array" {
		return hasFileProperty(s.Items)
	}
	return s.Format == "binary"
}
