package doc

import (
	"bytes"
	"encoding/json"
	"github.com/aiechoic/admin/core/gins"
	"github.com/aiechoic/admin/core/openapi"
	"github.com/aiechoic/admin/pkg/errs"
	"github.com/aiechoic/admin/pkg/rsp"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

var redocHtml = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation</title>
	<redoc spec-url="openapi.json"></redoc>
    <script src="https://unpkg.com/redoc@2.2.0/bundles/redoc.standalone.js"> </script>
</head>
<body>
</body>
</html>
`

var swaggerHtml = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
  </head>
  <body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js" crossorigin></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-standalone-preset.js" crossorigin></script>
  <script>
    window.onload = () => {
	  window.ui = SwaggerUIBundle({
		url: "openapi.json",
		dom_id: '#swagger-ui',
		deepLinking: true,
		presets: [
		  SwaggerUIBundle.presets.apis,
		  SwaggerUIStandalonePreset
		],
		plugins: [
		  SwaggerUIBundle.plugins.DownloadUrl
		],
		layout: "StandaloneLayout",
		persistAuthorization: true,
	  });
    };
  </script>
  </body>
</html>
`

func NewService(api *openapi.Openapi) *gins.Service {
	for _, server := range api.Servers {
		log.Printf("serve swagger-ui at %s%s\n", server.Url, "/docs/swagger.html")
		log.Printf("serve redoc at %s%s\n", server.Url, "/docs/redoc.html")
	}
	lastModify := time.Now()
	swaggerHtmlReader := strings.NewReader(swaggerHtml)
	redocHtmlReader := strings.NewReader(redocHtml)

	// initialize when first request comes, because data may not be ready.
	var openapiJsonReader *bytes.Reader
	var errorCodesReader *bytes.Reader
	return &gins.Service{
		Tag:  "Docs",
		Path: "/docs",
		Routes: []gins.Route{
			{
				Method: "GET",
				Path:   "openapi.json",
				Handler: gins.Handler{
					Response: gins.Response{
						Json: openapi.Openapi{},
					},
					Handle: func(c *gin.Context) {
						if openapiJsonReader == nil {
							jsonData, err := json.Marshal(api)
							if err != nil {
								rsp.SendError(c, errs.InternalServerError, err)
								return
							}
							openapiJsonReader = bytes.NewReader(jsonData)
						}
						http.ServeContent(c.Writer, c.Request, "openapi.json", lastModify, openapiJsonReader)
					},
				},
			},
			{
				Method: "GET",
				Path:   "swagger.html",
				Handler: gins.Handler{
					Response: gins.Response{
						Contents: gins.ContentsTextHtml,
					},
					Handle: func(c *gin.Context) {
						http.ServeContent(c.Writer, c.Request, "swagger.html", lastModify, swaggerHtmlReader)
					},
				},
			},
			{
				Method: "GET",
				Path:   "redoc.html",
				Handler: gins.Handler{
					Response: gins.Response{
						Contents: gins.ContentsTextHtml,
					},
					Handle: func(c *gin.Context) {
						http.ServeContent(c.Writer, c.Request, "redoc.html", lastModify, redocHtmlReader)
					},
				},
			},
			{
				Method: "GET",
				Path:   "error_codes.json",
				Handler: gins.Handler{
					Response: gins.Response{
						Contents: gins.ContentsJson,
					},
					Handle: func(c *gin.Context) {
						if errorCodesReader == nil {
							errorCodes, err := json.Marshal(errs.GetCodes())
							if err != nil {
								rsp.SendError(c, errs.InternalServerError, err)
								return
							}
							errorCodesReader = bytes.NewReader(errorCodes)
						}
						http.ServeContent(c.Writer, c.Request, "error_codes.json", lastModify, errorCodesReader)
					},
				},
			},
		},
	}
}
