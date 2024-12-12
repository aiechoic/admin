package openapi

type Contact struct {
	Name  string `json:"name,omitempty" mapstructure:"name"`
	Url   string `json:"url,omitempty" mapstructure:"url"`
	Email string `json:"email,omitempty" mapstructure:"email"`
}

type License struct {
	Name       string `json:"name,omitempty" mapstructure:"name"`
	Url        string `json:"url,omitempty" mapstructure:"url"`
	Identifier string `json:"identifier,omitempty" mapstructure:"identifier"`
}

type Info struct {
	Title          string   `json:"title,omitempty" mapstructure:"title"`
	Summary        string   `json:"summary,omitempty" mapstructure:"summary"`
	Description    string   `json:"description,omitempty" mapstructure:"description"`
	TermsOfService string   `json:"termsOfService,omitempty" mapstructure:"terms_of_service"`
	Contact        *Contact `json:"contact,omitempty" mapstructure:"contact"`
	License        *License `json:"license,omitempty" mapstructure:"license"`
	Version        string   `json:"version,omitempty" mapstructure:"version"`
}

type Server struct {
	Url         string `json:"url,omitempty" mapstructure:"url"`
	Description string `json:"description,omitempty" mapstructure:"description"`
}

type SecuritySchemes map[string]*SecurityScheme

type SecuritySchemeType string

const (
	SecuritySchemeTypeApiKey    SecuritySchemeType = "apiKey"
	SecuritySchemeTypeHttp      SecuritySchemeType = "http"
	SecuritySchemeTypeOauth     SecuritySchemeType = "oauth2"
	SecuritySchemeTypeOpenId    SecuritySchemeType = "openIdConnect"
	SecuritySchemeTypeMutualTLS SecuritySchemeType = "mutualTLS"
)

type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type SecurityScheme struct {
	Type             SecuritySchemeType `json:"type,omitempty"`
	Description      string             `json:"description,omitempty"`
	Name             string             `json:"name,omitempty"`
	In               string             `json:"in,omitempty"` // header, query, cookie
	Scheme           string             `json:"scheme,omitempty"`
	BearerFormat     string             `json:"bearerFormat,omitempty"`
	OpenIdConnectUrl string             `json:"openIdConnectUrl,omitempty"`
	Flows            *OAuthFlows        `json:"flows,omitempty"`
}

type Components struct {
	SecuritySchemes SecuritySchemes    `json:"securitySchemes,omitempty"`
	Schemas         map[string]*Schema `json:"schemas,omitempty"`
}

type ContentType string

const (
	ContentTypeJson          ContentType = "application/json"
	ContentTypeXml           ContentType = "application/xml"
	ContentTypeHtml          ContentType = "text/html"
	ContentTypeForm          ContentType = "application/x-www-form-urlencoded"
	ContentTypeMultipartForm ContentType = "multipart/form-data"
)

type MediaType struct {
	Schema  *Schema `json:"schema,omitempty"`
	Example any     `json:"example,omitempty"`
}

type RequestBody struct {
	Description string                     `json:"description,omitempty"`
	Content     map[ContentType]*MediaType `json:"content,omitempty"`
	Required    bool                       `json:"required,omitempty"`
}

type ResponseCode string

type Parameter struct {
	Name            string  `json:"name,omitempty"`
	In              string  `json:"in,omitempty"` // query, path, header, cookie
	Schema          *Schema `json:"schema,omitempty"`
	Description     string  `json:"description,omitempty"`
	Required        bool    `json:"required,omitempty"`
	Deprecated      bool    `json:"deprecated,omitempty"`
	AllowEmptyValue bool    `json:"allowEmptyValue,omitempty"`
}

type ResponseBody struct {
	Description string                     `json:"description"`
	Content     map[ContentType]*MediaType `json:"content,omitempty"`
}

type Operation struct {
	Tags        []string                       `json:"tags,omitempty"`
	Summary     string                         `json:"summary,omitempty"`
	Description string                         `json:"description,omitempty"`
	RequestBody *RequestBody                   `json:"requestBody,omitempty"`
	Parameters  []*Parameter                   `json:"parameters,omitempty"`
	Responses   map[ResponseCode]*ResponseBody `json:"responses,omitempty"`
	Deprecated  bool                           `json:"deprecated,omitempty"`
	Security    []map[string][]string          `json:"security,omitempty"`
}

type PathItem map[string]*Operation

type Tag struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Openapi struct {
	Openapi    string              `json:"openapi,omitempty" mapstructure:"openapi"`
	Info       *Info               `json:"info,omitempty" mapstructure:"info"`
	Servers    []*Server           `json:"servers,omitempty" mapstructure:"servers"`
	Components Components          `json:"components,omitempty" mapstructure:"-"`
	Paths      map[string]PathItem `json:"paths,omitempty" mapstructure:"-"`
	Tags       []*Tag              `json:"tags,omitempty" mapstructure:"-"`
}

type SecurityRoute struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (o *Openapi) AddComponentsSchemas(schemas map[string]*Schema) {
	if o.Components.Schemas == nil {
		o.Components.Schemas = map[string]*Schema{}
	}
	for k, v := range schemas {
		o.Components.Schemas[k] = v
	}
}

// GetAllSecurityRoutes 获取所有具有权限控制的路由, 按照 tag 分组, 用于生成权限控制文档
func (o *Openapi) GetAllSecurityRoutes() map[string][]*SecurityRoute {
	var routes = map[string][]*SecurityRoute{}
	for path, item := range o.Paths {
		for method, op := range item {
			if len(op.Security) > 0 {
				routes[op.Tags[0]] = append(routes[op.Tags[0]], &SecurityRoute{
					Method: method,
					Path:   path,
				})
			}
		}
	}
	return routes
}
