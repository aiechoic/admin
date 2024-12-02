package doc

import (
	"github.com/aiechoic/admin/core/gins"
	"github.com/aiechoic/admin/core/openapi"
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

func NewDocService(api *openapi.Openapi) *gins.Service {
	for _, server := range api.Servers {
		log.Printf("serve swagger-ui at %s%s\n", server.Url, "/docs/swagger.html")
		log.Printf("serve redoc at %s%s\n", server.Url, "/docs/redoc.html")
	}
	lastModify := time.Now()
	SwaggerHtmlReader := strings.NewReader(swaggerHtml)
	RedocHtmlReader := strings.NewReader(redocHtml)
	return &gins.Service{
		Tag:  "DOCS",
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
						c.JSON(200, api)
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
						http.ServeContent(c.Writer, c.Request, "swagger.html", lastModify, SwaggerHtmlReader)
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
						http.ServeContent(c.Writer, c.Request, "redoc.html", lastModify, RedocHtmlReader)
					},
				},
			},
		},
	}
}
