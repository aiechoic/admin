package openapi

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
	"time"
)

func TestNewSchema(t *testing.T) {
	type testcase struct {
		v      any
		schema *Schema
	}
	type NestedC struct {
		C bool `json:"c" binding:"required" description:"c"`
	}
	type NestedD struct {
		D []byte `json:"d" binding:"required" description:"nested d"`
	}
	var testcases = []testcase{
		{int(0), &Schema{Type: "integer"}},
		{int8(0), &Schema{Type: "integer"}},
		{int16(0), &Schema{Type: "integer"}},
		{int32(0), &Schema{Type: "integer"}},
		{int64(0), &Schema{Type: "integer"}},
		{uint(0), &Schema{Type: "integer"}},
		{uint8(0), &Schema{Type: "integer"}},
		{uint16(0), &Schema{Type: "integer"}},
		{uint32(0), &Schema{Type: "integer"}},
		{uint64(0), &Schema{Type: "integer"}},
		{float32(0), &Schema{Type: "number", Format: "float"}},
		{float64(0), &Schema{Type: "number", Format: "double"}},
		{"", &Schema{Type: "string"}},
		{false, &Schema{Type: "boolean"}},
		{[]byte{}, &Schema{Type: "string", Format: "byte"}},
		{[]int{}, &Schema{Type: "array", Items: &Schema{Type: "integer"}}},
		{[1]int{}, &Schema{Type: "array", Items: &Schema{Type: "integer"}}},
		{interface{}(nil), &Schema{Type: "null"}},
		{map[string]int{}, &Schema{Type: "object"}},
		{map[string]string{"k1": "v1"}, &Schema{Type: "object", Properties: map[string]*Schema{"k1": {Type: "string"}}}},
		{time.Time{}, &Schema{Type: "string", Format: "date-time"}},
		{multipart.FileHeader{}, &Schema{Type: "string", Format: "binary"}},
		{&time.Time{}, &Schema{Type: "string", Format: "date-time", Nullable: true}},
		{getNew(0), &Schema{Type: "integer", Nullable: true}},
		{struct {
			A int    `json:"a" binding:"required" format:"int64" description:"a"`
			B string `json:"b" binding:"required,email" description:"b"`
			*NestedC
			NestedD
			D []int `json:"d" description:"outer d"`
		}{}, &Schema{
			Type:     "object",
			Required: []string{"a", "b", "c"},
			Properties: map[string]*Schema{
				"a": {Type: "integer", Description: "a", Format: "int64"},
				"b": {Type: "string", Description: "b", Format: "email"},
				"c": {Type: "boolean", Description: "c"},
				"d": {Type: "array", Description: "outer d", Items: &Schema{Type: "integer"}},
			},
		},
		},
	}
	t.Run("testcases", func(t *testing.T) {
		for _, tc := range testcases {
			schema, _ := NewSchema(tc.v, "json")
			assert.Equal(t, getJson(tc.schema), getJson(schema))
		}
	})

	t.Run("recursive", func(t *testing.T) {
		// 创建测试数据
		a := &RecursiveTypeA{}
		// 调用生成 schema
		gotSchema, gotRefs := NewSchema(a, "json")
		needSchema := &Schema{
			Type:     "object",
			Nullable: true,
			Properties: map[string]*Schema{
				"b": {
					Type:     "object",
					Nullable: true,
					Properties: map[string]*Schema{
						"a": {
							Ref: "#/components/schemas/RecursiveTypeA",
						},
					},
				},
			},
		}
		needRefs := map[string]*Schema{
			"RecursiveTypeA": {
				Type:     "object",
				Nullable: true,
				Properties: map[string]*Schema{
					"b": {
						Type:     "object",
						Nullable: true,
						Properties: map[string]*Schema{
							"a": {
								Ref: "#/components/schemas/RecursiveTypeA",
							},
						},
					},
				},
			},
		}
		assert.Equal(t, getJson(needSchema), getJson(gotSchema))
		assert.Equal(t, getJson(needRefs), getJson(gotRefs))
	})
}

type RecursiveTypeA struct {
	B *RecursiveTypeB `json:"b"`
}

type RecursiveTypeB struct {
	A *RecursiveTypeA `json:"a"`
}

func getNew[T any](v T) *T {
	return &v
}

func getJson(v any) string {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(v)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
