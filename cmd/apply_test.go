package cmd

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var kongConfig = `services:
  - name: my-service-dev
    tags: [ ns.aps-moh-proto ]
`
var clientCredConfig = `kind: Namespace
name: ns-sampler
displayName: ns-sampler Display Name
---
kind: GatewayService
name: my-service-dev
---
kind: CredentialIssuer
name: aps-moh-proto default
---
kind: DraftDataset
name: my-service-dataset
---
kind: Product
name: my-service API
`

func TestApplyOptions(t *testing.T) {
	fileName := "gw-config.yaml"
	dir := t.TempDir()
	config, err := os.Create(filepath.Join(dir, fileName))
	io.WriteString(config, clientCredConfig)
	defer config.Close()
	if err != nil {
		t.Fatal(err)
	}
	o := &ApplyOptions{
		input: fileName,
	}
	output, err := o.Parse(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(output))
}

func TestNonYamlFile(t *testing.T) {
	fileName := "gw-config.json"
	dir := t.TempDir()
	config, err := os.Create(filepath.Join(dir, fileName))
	defer config.Close()
	if err != nil {
		t.Fatal(err)
	}
	o := &ApplyOptions{
		input: fileName,
	}
	_, err = o.Parse(dir)
	assert.Error(t, err)
}

func TestExtractResouceConfig(t *testing.T) {
	input := []byte(`kind: GatewayService
name: my-service-dev
`)
	result, err := ExtractResourceConfig(input)
	if err != nil {
		t.Fatal(err)
	}
	expect := &ResourceConfig{
		Kind: "GatewayService",
		Config: map[string]interface{}{
			"name": "my-service-dev",
		},
	}
	assert.Equal(t, expect, result)
}

func TestResourceConfigAction(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]interface{}
		expect string
	}{
		{
			name: "CredentialIssuer",
			input: map[string]interface{}{
				"kind": "CredentialIssuer",
				"name": "my-service",
			},
			expect: "issuer",
		},
		{
			name: "DraftDataset",
			input: map[string]interface{}{
				"kind": "DraftDataset",
				"name": "my-service",
			},
			expect: "dataset",
		},
		{
			name: "Product",
			input: map[string]interface{}{
				"kind": "Product",
				"name": "my-service",
			},
			expect: "product",
		},
		{
			name: "Environment",
			input: map[string]interface{}{
				"kind": "Environment",
				"name": "my-service",
			},
			expect: "environment",
		},
		{
			name: "GatewayService",
			input: map[string]interface{}{
				"kind": "CredentialIssuer",
				"name": "my-service",
			},
			expect: "publishGateway",
		},
		{
			name: "AnotherItem",
			input: map[string]interface{}{
				"kind": "CredentialIssuer",
				"name": "my-service",
			},
			expect: "skip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ResourceConfig{
				Kind:   tt.name,
				Config: tt.input,
			}
			assert.Equal(t, tt.expect, e.Action())
		})
	}
}

func TestCounter(t *testing.T) {
	c := &PublishCounter{}
	c.AddSuccess()
	c.AddSuccess()
	c.AddSkipped()
	c.AddFailed()
	assert.Equal(t, 2, c.Success)
	assert.Equal(t, 1, c.Skipped)
	assert.Equal(t, 1, c.Failed)
	assert.Equal(t, "2/3 Published, 1 Skipped", c.Print())
}

func TestPublishResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("PUT", "https://aps.gov.bc.ca/ds/api/v2/namespaces/ns-sampler/issuers", func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"result": "Issuer published",
		})
	})
	ctx := &pkg.AppContext{
		Namespace: "ns-sampler",
		Host:      "aps.gov.bc.ca",
	}
	doc := parsedConfig{
		"name": "my-service",
	}
	result, err := PublishResource(ctx, doc, "issuer")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Issuer published", result)
}

func TestPublishGatewayService(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("PUT", "https://aps.gov.bc.ca/gw/api/namespaces/ns-sampler/gateway", func(r *http.Request) (*http.Response, error) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			return nil, err
		}
		configFile := r.MultipartForm.File["configFile"]
		contents, err := configFile[0].Open()
		defer contents.Close()
		if err != nil {
			return nil, err
		}
		c, err := io.ReadAll(contents)
		if err != nil {
			return nil, err
		}
		assert.Equal(t, string(c), `{"services":[{"name":"my-service","routes":[{"name":"my-route"}]}]}`)
		return httpmock.NewJsonResponse(200, "{}")
	})
	doc := parsedConfig{
		"name": "my-service",
		"routes": []map[string]interface{}{
			{"name": "my-route"},
		},
	}
	ctx := &pkg.AppContext{
		Namespace: "ns-sampler",
		Host:      "aps.gov.bc.ca",
	}
	err := PublishGatewayService(ctx, doc)
	assert.NoError(t, err)
}
