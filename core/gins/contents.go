package gins

import "github.com/aiechoic/admin/core/openapi"

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

	// ContentsJson defines the response body content type is "application/json"
	ContentsJson = map[openapi.ContentType]*openapi.MediaType{
		"application/json": {
			Schema: &openapi.Schema{Type: "object"},
		},
	}

	// ContentsXML defines the response body content type is "application/xml"
	ContentsXML = map[openapi.ContentType]*openapi.MediaType{
		"application/xml": {
			Schema: &openapi.Schema{Type: "object"},
		},
	}

	// ContentsTextPlain defines the response body content type is "text/plain"
	ContentsTextPlain = map[openapi.ContentType]*openapi.MediaType{
		"text/plain": {
			Schema: &openapi.Schema{Type: "string"},
		},
	}

	// ContentsImages defines the response body content type is "image/*"
	ContentsImages = map[openapi.ContentType]*openapi.MediaType{
		"image/*": {
			Schema: &openapi.Schema{Type: "string", Format: "binary"},
		},
	}
)
