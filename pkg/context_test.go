package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContxtUrl(t *testing.T) {
	x := AppContext{
		ApiHost: "api.bc.gov.ca",
	}
	tests := []struct {
		name   string
		params interface{}
		path   string
		expect string
	}{
		{
			name:   "builds a URL without params",
			params: nil,
			path:   "/status",
			expect: "https://api.bc.gov.ca/status",
		},
		{
			name: "builds a URL with one param",
			params: struct {
				Name string `url:"name"`
			}{
				"ns-sampler",
			},
			path:   "/namespace",
			expect: "https://api.bc.gov.ca/namespace?name=ns-sampler",
		},
		{
			name: "builds a URL with multiple param types",
			params: struct {
				X int  `url:"x"`
				Y bool `url:"y"`
			}{
				1,
				true,
			},
			path:   "/path/with/segments/123",
			expect: "https://api.bc.gov.ca/path/with/segments/123?x=1&y=true",
		},
	}
	for _, tt := range tests {
		url, _ := x.CreateUrl(tt.path, tt.params)
		assert.Equal(t, tt.expect, url, tt.name)
	}
}
