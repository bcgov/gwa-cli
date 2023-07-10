package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextUrl(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			url, _ := x.CreateUrl(tt.path, tt.params)
			assert.Equal(t, tt.expect, url, tt.name)
		})
	}
}

func TestCreateUrl(t *testing.T) {
	tests := []struct {
		name   string
		expect string
		ctx    AppContext
		hasErr bool
		params interface{}
	}{
		{
			name:   "ApiHost set",
			expect: "https://api.gov.bc.ca/status",
			ctx: AppContext{
				ApiHost: "api.gov.bc.ca",
			},
		},
		{
			name:   "Host set",
			expect: "https://local.test/status",
			ctx: AppContext{
				ApiHost: "api.gov.bc.ca",
				Host:    "local.test",
			},
		},
		{
			name:   "Correctly formated params",
			expect: "https://local.test/status?hello=world",
			ctx: AppContext{
				ApiHost: "api.gov.bc.ca",
				Host:    "local.test",
			},
			params: struct {
				Hello string `url:"hello"`
			}{
				Hello: "world",
			},
		},
		{
			name:   "no host",
			ctx:    AppContext{},
			hasErr: true,
		},
		{
			name:   "incorrectly formatted params",
			ctx:    AppContext{},
			hasErr: true,
			params: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := tt.ctx.CreateUrl("/status", tt.params)
			if tt.hasErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expect, url)
			}
		})
	}
}
