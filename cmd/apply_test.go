package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

var kongConfig = `services:
  - name: my-service-dev
    tags: [ ns.aps-moh-proto ]
`

var input = `kind: Namespace
name: ns-sampler
displayName: ns-sampler Display Name
---
kind: GatewayService
name: service-1
host: api.co1.com
---
kind: GatewayService
name: service-2
host: api.co2.com
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
	io.WriteString(config, input)
	defer config.Close()
	if err != nil {
		t.Fatal(err)
	}
	o := &ApplyOptions{
		cwd:   dir,
		input: fileName,
	}
	err = o.Parse()
	if err != nil {
		t.Fatal(err)
	}

	expected := []interface{}{
		GatewayService{Config: []map[string]interface{}{
			{
				"name": "service-1",
				"host": "api.co1.com",
			},
			{
				"name": "service-2",
				"host": "api.co2.com",
			},
		}},
		Skipped{Name: "ns-sampler", Kind: "Namespace"},
		Resource{Kind: "CredentialIssuer", Config: map[string]interface{}{"name": "aps-moh-proto default"}},
		Resource{Kind: "DraftDataset", Config: map[string]interface{}{"name": "my-service-dataset"}},
		Resource{Kind: "Product", Config: map[string]interface{}{"name": "my-service API"}},
	}

	assert.Equal(t, expected, o.output, "outputs a map keyed by type, with grouped gateways")
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
		cwd:   dir,
		input: fileName,
	}
	err = o.Parse()
	assert.Error(t, err)
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
			name: "AnotherItem",
			input: map[string]interface{}{
				"kind": "CredentialIssuer",
				"name": "my-service",
			},
			expect: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Resource{
				Kind:   tt.name,
				Config: tt.input,
			}
			assert.Equal(t, tt.expect, e.GetAction())
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
	httpmock.RegisterResponder(
		"PUT",
		"https://aps.gov.bc.ca/ds/api/v2/namespaces/ns-sampler/issuers",
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"result": "Issuer published",
			})
		},
	)
	ctx := &pkg.AppContext{
		ApiVersion: "v2",
		Namespace:  "ns-sampler",
		Host:       "aps.gov.bc.ca",
	}
	doc := map[string]interface{}{
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
	httpmock.RegisterResponder(
		"PUT",
		"https://aps.gov.bc.ca/gw/api/v2/namespaces/ns-sampler/gateway",
		func(r *http.Request) (*http.Response, error) {
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
			assert.Equal(
				t,
				string(c),
				`{"services":[{"name":"service-1","routes":[{"name":"api.co1.com/route"}]},{"name":"service-2","routes":[{"name":"api.co2.com/route"}]}]}`,
			)
			return httpmock.NewJsonResponse(200, "{}")
		},
	)
	doc := []map[string]interface{}{
		{
			"name": "service-1",
			"routes": []map[string]interface{}{
				{"name": "api.co1.com/route"},
			},
		},
		{
			"name": "service-2",
			"routes": []map[string]interface{}{
				{"name": "api.co2.com/route"},
			},
		},
	}
	ctx := &pkg.AppContext{
		ApiVersion: "v2",
		Namespace:  "ns-sampler",
		Host:       "aps.gov.bc.ca",
	}
	res, err := PublishGatewayService(ctx, doc)
	assert.NoError(t, err)
	assert.Equal(t, "", res.Results, "returns a successful gateway service resposne like in publish-gateway")
}

func TestApplyStdout(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
		expected     []string
	}{
		{
			name:         "Success output",
			responseCode: 200,
			expected: []string{
				"↑ Publishing Gateway Services",
				"✓ Gateway Services published",
				"Pubished: 2\nSkipped: 1",
				"4/4 Published, 1 Skipped",
				"- [Namespace] ns-sampler",
				"↑ [CredentialIssuer] aps-moh-proto default",
				"✓ [CredentialIssuer] aps-moh-proto default: Published",
				"↑ [DraftDataset] my-service-dataset",
				"✓ [DraftDataset] my-service-dataset: Published",
				"↑ [Product] my-service API",
				"✓ [Product] my-service API: Published",
			},
		},
		{
			name:         "Failed output",
			responseCode: 401,
			expected: []string{
				"↑ Publishing Gateway Services",
				"x Gateway Services publish failed",
				"0/4 Published, 1 Skipped",
				"- [Namespace] ns-sampler",
				"↑ [CredentialIssuer] aps-moh-proto default",
				"x [CredentialIssuer] aps-moh-proto default failed",
				"↑ [DraftDataset] my-service-dataset",
				"x [DraftDataset] my-service-dataset failed",
				"↑ [Product] my-service API",
				"x [Product] my-service API failed",
			},
		},
	}

	for _, tt := range tests {
		// Mocks
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		regexPattern := `=~^https://api\.gov\.bc\.ca/ds/api/v2/namespaces/ns-sampler/\w+$`
		httpmock.RegisterResponder("PUT", regexPattern, func(_ *http.Request) (*http.Response, error) {
			fmt.Println("Resource publish")
			return httpmock.NewJsonResponse(tt.responseCode, map[string]interface{}{
				"result": "Published",
			})
		})
		httpmock.RegisterResponder(
			"PUT",
			"https://api.gov.bc.ca/gw/api/v2/namespaces/ns-sampler/gateway",
			func(_ *http.Request) (*http.Response, error) {
				fmt.Println("gateway publish")
				return httpmock.NewJsonResponse(tt.responseCode, map[string]interface{}{
					"results": "Pubished: 2\nSkipped: 1",
				})
			},
		)

		// Setup
		cwd := t.TempDir()
		ctx := &pkg.AppContext{
			Cwd:        cwd,
			Namespace:  "ns-sampler",
			ApiHost:    "api.gov.bc.ca",
			ApiVersion: "v2",
		}
		filename := "gw-config.yaml"
		os.WriteFile(filepath.Join(cwd, filename), []byte(input), 0644)

		args := []string{"apply", "--input", filename}

		mainCmd := &cobra.Command{
			Use: "gwa",
		}
		mainCmd.AddCommand(NewApplyCmd(ctx))
		mainCmd.SetArgs(args)
		out := capturer.CaptureOutput(func() {
			mainCmd.Execute()
		})
		for _, e := range tt.expected {
			assert.Contains(t, out, e)
		}
	}
}
