package openapi

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"
)

type Schema struct {
	Ref         string             `json:"$ref,omitempty"`
	Description string             `json:"description,omitempty"`
	Type        string             `json:"type,omitempty"`
	Nullable    bool               `json:"nullable,omitempty"`
	Enum        []string           `json:"enum,omitempty"`
	Format      string             `json:"format,omitempty"`
	Required    []string           `json:"required,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
}

func NewSchema(v any, tag string) (schema *Schema, refs map[string]*Schema) {
	sb := newSchemaBuilder(tag)
	schema = sb.newSchema(reflect.ValueOf(v))
	refs = sb.getPointerRefs()
	return schema, refs
}

type schemaBuilder struct {
	refs      map[string]*Schema
	pointers  []string
	tag       string
	modelPath string
}

func newSchemaBuilder(tag string) *schemaBuilder {
	return &schemaBuilder{
		refs:      map[string]*Schema{},
		tag:       tag,
		modelPath: "#/components/schemas/",
	}
}

func (sb *schemaBuilder) getPointerRefs() map[string]*Schema {
	refs := map[string]*Schema{}
	for _, pointer := range sb.pointers {
		refs[strings.TrimPrefix(pointer, sb.modelPath)] = sb.refs[pointer]
	}
	return refs
}

var fileHeaderType = reflect.TypeOf(multipart.FileHeader{})
var timeType = reflect.TypeOf(time.Time{})

func (sb *schemaBuilder) newSchema(v reflect.Value) *Schema {
	if !v.IsValid() {
		return &Schema{Type: "null"}
	}
	t := v.Type()
	switch t.Kind() {
	case reflect.Bool:
		return &Schema{Type: "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Schema{Type: "integer"}
	case reflect.Float32:
		return &Schema{Type: "number", Format: "float"}
	case reflect.Float64:
		return &Schema{Type: "number", Format: "double"}
	case reflect.String:
		return &Schema{Type: "string"}
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return &Schema{Type: "string", Format: "byte"}
		} else {
			ev := reflect.Zero(t.Elem())
			return &Schema{Type: "array", Items: sb.newSchema(ev)}
		}
	case reflect.Interface:
		if v.IsValid() {
			return sb.newSchema(v.Elem())
		} else {
			return &Schema{Type: "null"}
		}
	case reflect.Struct:
		if t == fileHeaderType {
			return &Schema{
				Type:   "string",
				Format: "binary",
			}
		}
		if t == timeType {
			return &Schema{
				Type:   "string",
				Format: "date-time",
			}
		}
		return sb.newStructSchema(v)
	case reflect.Ptr:
		var ev reflect.Value
		if v.IsNil() {
			ev = reflect.Zero(t.Elem())
		} else {
			ev = v.Elem()
		}
		s := sb.newSchema(ev)
		if s.Ref == "" {
			s.Nullable = true
		}
		return s
	case reflect.Map:
		return sb.newMapSchema(v)
	default:
		return &Schema{Type: "null"}
	}
}

func setSchemaRequired(schema *Schema, property string, required bool) {
	if required {
		if !slices.Contains(schema.Required, property) {
			schema.Required = append(schema.Required, property)
		}
	} else {
		if i := slices.Index(schema.Required, property); i != -1 {
			schema.Required = append(schema.Required[:i], schema.Required[i+1:]...)
		}
	}
}

func (sb *schemaBuilder) newStructSchema(v reflect.Value) *Schema {
	uniqueName := getStructUniqueName(v.Type(), 0)
	ref := sb.modelPath + uniqueName
	if _, ok := sb.refs[ref]; ok {
		sb.pointers = append(sb.pointers, ref)
		return &Schema{Ref: ref}
	}
	schema := &Schema{Type: "object", Properties: map[string]*Schema{}}
	sb.refs[ref] = schema
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.PkgPath != "" {
			continue
		}
		if field.Anonymous {
			anonymous := sb.newSchema(v.Field(i))
			for k, s := range anonymous.Properties {
				if _, ok := schema.Properties[k]; !ok {
					setSchemaRequired(schema, k, slices.Contains(anonymous.Required, k))
					schema.Properties[k] = s
				}
			}
			continue
		}
		ftg := field.Tag
		name := ftg.Get(sb.tag)
		if name == "-" {
			continue
		}
		if name != "" {
			name = strings.Split(name, ",")[0]
		}
		if name == "" {
			name = field.Name
		}
		s := sb.newSchema(v.Field(i))
		if desc := ftg.Get("description"); desc != "" {
			s.Description = desc
		}
		var required bool
		if binding := ftg.Get("binding"); binding != "" {
			params := strings.Split(binding, ",")
			for _, param := range params {
				if param == "required" {
					required = true
				} else if param == "email" {
					s.Format = "email"
				}
			}
		}
		if format := field.Tag.Get("format"); format != "" {
			s.Format = format
		}
		schema.Properties[name] = s
		setSchemaRequired(schema, name, required)
	}
	return schema
}

func (sb *schemaBuilder) newMapSchema(v reflect.Value) *Schema {
	if v.IsNil() {
		return &Schema{Type: "object"}
	}
	if v.Type().Key().Kind() != reflect.String {
		return &Schema{Type: "null"}
	}
	s := &Schema{
		Type:       "object",
		Properties: map[string]*Schema{},
	}
	for _, mKey := range v.MapKeys() {
		s.Properties[mKey.String()] = sb.newSchema(v.MapIndex(mKey))
	}
	return s
}

var uniqueTypeCache = sync.Map{}
var uniqueNameCache = sync.Map{}

func getStructUniqueName(t reflect.Type, index int) string {
	if uniqueName, ok := uniqueTypeCache.Load(t); ok {
		return uniqueName.(string)
	}
	var uniqueName string
	if index == 0 {
		uniqueName = fmt.Sprintf("%s", t.Name())
	} else {
		uniqueName = fmt.Sprintf("%s-%d", t.Name(), index)
	}
	if _, ok := uniqueNameCache.Load(uniqueName); ok {
		return getStructUniqueName(t, index+1)
	}
	uniqueTypeCache.Store(t, uniqueName)
	uniqueNameCache.Store(uniqueName, struct{}{})
	return uniqueName
}
