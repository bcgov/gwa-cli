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
	"github.com/stretchr/testify/assert"
)

func TestPublishUrls(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		slug string
	}{
		{
			arg:  "content",
			slug: "contents",
		},
		{
			arg:  "dataset",
			slug: "datasets",
		},
		{
			arg:  "product",
			slug: "products",
		},
		{
			arg:  "issuer",
			slug: "issuers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			ctx := &pkg.AppContext{
				ApiHost:   "api.gov.bc.ca",
				Namespace: "ns-sampler",
			}
			URL := fmt.Sprintf("https://%s/ds/api/v2/namespaces/ns-sampler/%s", ctx.ApiHost, tt.slug)
			httpmock.RegisterResponder("PUT", URL, func(r *http.Request) (*http.Response, error) {
				assert.Contains(t, r.URL.Path, tt.slug)
				return httpmock.NewJsonResponse(200, "{}")
			})
			err := Publish(ctx, nil, tt.arg)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestConvertYamlToJson(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name: "product yaml",
			input: `name: my-new-product
appId: "000000000000"
environments:
  - name: dev
    active: false
    approval: false
    flow: public
    appId: "00000000"`,
			expect: `{"appId":"000000000000","environments":[{"active":false,"appId":"00000000","approval":false,"flow":"public","name":"dev"}],"name":"my-new-product"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			fileName := filepath.Join(dir, "target.yaml")
			f, err := os.Create(fileName)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			io.WriteString(f, tt.input)

			ctx := &pkg.AppContext{
				Cwd: dir,
			}
			opts := &PublishOptions{
				input: "./target.yaml",
			}
			out, err := opts.ParseInput(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expect, string(out))
		})
	}
}
