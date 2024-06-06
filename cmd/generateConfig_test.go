package cmd

import (
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/stretchr/testify/assert"
)

func TestParseUpstream(t *testing.T) {
	tests := []struct {
		name   string
		input  *GenerateConfigOptions
		expect string
	}{
		{
			name: "upstream port set",
			input: &GenerateConfigOptions{
				Upstream: "https://test.com:8000",
			},
			expect: "8000",
		},
		{
			name: "upstream https set",
			input: &GenerateConfigOptions{
				Upstream: "https://test.com",
			},
			expect: "443",
		},
		{
			name: "upstream http set",
			input: &GenerateConfigOptions{
				Upstream: "http://test.com",
			},
			expect: "80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.ParseUpstream()
			assert.Equal(t, tt.expect, tt.input.UpstreamPort)
		})
	}
}

func TestValidateTemplate(t *testing.T) {
	kong := &GenerateConfigOptions{
		Template: "kong-httpbin",
	}
	assert.NoError(t, kong.ValidateTemplate())
	clientCreds := &GenerateConfigOptions{
		Template: "client-credentials-shared-idp",
	}
	assert.NoError(t, clientCreds.ValidateTemplate())
	bad := &GenerateConfigOptions{
		Template: "asdf",
	}
	assert.Error(t, bad.ValidateTemplate())
}

func TestGenerateKongConfig(t *testing.T) {
	dir := t.TempDir()
	ctx := &pkg.AppContext{
		Cwd: dir,
	}
	opts := &GenerateConfigOptions{
		Namespace:    "sampler",
		Template:     "kong-httpbin",
		Service:      "my-service",
		UpstreamPort: "443",
		UpstreamUrl: &url.URL{
			Host:   "httpbin.org",
			Scheme: "https",
		},
		Out: "gw-config.yaml",
	}
	err := GenerateConfig(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.ReadFile(path.Join(ctx.Cwd, opts.Out))
	if err != nil {
		t.Fatal(err)
	}
	compare := string(file)
	assert.Contains(t, compare, "name: my-service-dev")
	assert.Contains(t, compare, "tags: [ ns.sampler ]")
	assert.Contains(t, compare, "host: httpbin.org")
	assert.Contains(t, compare, "port: 443")
	assert.Contains(t, compare, "protocol: https")
	assert.Contains(t, compare, "- my-service.dev.api.gov.bc.ca")
}

func TestClientCredentialsGenerator(t *testing.T) {
	dir := t.TempDir()
	ctx := &pkg.AppContext{
		Cwd: dir,
	}
	opts := &GenerateConfigOptions{
		Namespace:    "cc-sampler",
		Template:     "client-credentials-shared-idp",
		Service:      "my-service",
		UpstreamPort: "443",
		UpstreamUrl: &url.URL{
			Host:   "httpbin.org",
			Path:   "/post",
			Scheme: "https",
		},
		Organization:     "ministry-of-citizens-services",
		OrganizationUnit: "databc",
		Out:              "gw-config.yaml",
	}
	err := GenerateConfig(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.ReadFile(path.Join(ctx.Cwd, opts.Out))
	if err != nil {
		t.Fatal(err)
	}
	compare := string(file)
	assert.Contains(t, compare, "name: my-service-dev")
	assert.Contains(t, compare, "tags: [ns.cc-sampler]")
	assert.Contains(t, compare, "host: httpbin.org")
	assert.Contains(t, compare, "port: 443")
	assert.Contains(t, compare, "protocol: https")
	assert.Contains(t, compare, "- my-service.dev.api.gov.bc.ca")
	assert.Contains(t, compare, "paths: [/post]")
	assert.Contains(t, compare, "allowed_aud: ap-cc-sampler-default-dev")
}
