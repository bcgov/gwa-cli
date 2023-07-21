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

var inputContent = `name: my-new-product
appId: "000000000000"
environments:
  - name: dev
    active: false
    approval: false
    flow: public
    appId: "00000000"`
var commands = []struct {
	args   []string
	expect string
	slug   string
}{
	{
		args:   []string{"content", "--input", "./target.yaml"},
		slug:   "contents",
		expect: "Content successfully published",
	},
	{
		args:   []string{"dataset", "--input", "./target.yaml"},
		slug:   "datasets",
		expect: "Dataset successfully published",
	},
	{
		args:   []string{"product", "--input", "./target.yaml"},
		slug:   "products",
		expect: "Product successfully published",
	},
	{
		args:   []string{"issuer", "--input", "./target.yaml"},
		slug:   "issuers",
		expect: "Issuer successfully published",
	},
}

func WriteTestFile(dir string, content string) error {
	fileName := filepath.Join(dir, "target.yaml")
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	io.WriteString(f, content)
	return nil
}

func TestCommands(t *testing.T) {
	for _, tt := range commands {
		t.Run(tt.slug+" Command", func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			dir := t.TempDir()
			err := WriteTestFile(dir, inputContent)
			if err != nil {
				t.Fatal(err)
			}
			ctx := &pkg.AppContext{
				Namespace: "ns-sampler",
				ApiHost:   "api.gov.bc.ca",
				Cwd:       dir,
			}
			route := fmt.Sprintf("/ds/api/v2/namespaces/%s/%ss", ctx.Namespace, tt.args[0])
			URL, _ := ctx.CreateUrl(route, nil)

			httpmock.RegisterResponder("PUT", URL, func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"name": "AppName",
				})
			})

			args := append([]string{"publish"}, tt.args...)
			mainCmd := &cobra.Command{
				Use: "gwa",
			}
			mainCmd.AddCommand(NewPublishCmd(ctx))
			mainCmd.SetArgs(args)
			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})
			assert.Contains(t, out, tt.expect)
		})
	}
}

func TestPublishUrls(t *testing.T) {
	for _, tt := range commands {
		t.Run(tt.slug, func(t *testing.T) {
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
			err := Publish(ctx, nil, tt.args[0])
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
			name:   "product yaml",
			input:  inputContent,
			expect: `{"appId":"000000000000","environments":[{"active":false,"appId":"00000000","approval":false,"flow":"public","name":"dev"}],"name":"my-new-product"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			err := WriteTestFile(dir, tt.input)
			if err != nil {
				t.Fatal(err)
			}

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
