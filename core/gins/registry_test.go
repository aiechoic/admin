package gins

import "testing"

func TestConvertGinPathToSwaggerPath(t *testing.T) {
	tests := []struct {
		ginPath     string
		swaggerPath string
	}{
		{"/foo/:id", "/foo/{id}"},
		{"/foo/*action", "/foo/{action}"},
		{"/foo/bar", "/foo/bar"},
		{"/:id/foo/*action", "/{id}/foo/{action}"},
		{"/", "/"},
	}

	for _, test := range tests {
		result := convertGinPathToSwaggerPath(test.ginPath)
		if result != test.swaggerPath {
			t.Errorf("convertGinPathToSwaggerPath(%s) = %s; want %s", test.ginPath, result, test.swaggerPath)
		}
	}
}
