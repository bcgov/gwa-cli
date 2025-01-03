package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func setup(dir string) error {
	fileName := ".gwa-config.yaml"
	path := path.Join(dir, fileName)
	configFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer configFile.Close()
	viper.AddConfigPath(dir)
	viper.SetConfigFile(path)
	return nil
}

func TestGatewayCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expect   string
		method   string
		gateway  string
		response httpmock.Responder
	}{
		{
			name: "list gateways",
			args: []string{"list"},
			expect: `Display Name     Gateway ID  
janis's Gateway  gw-1        
janis's Gateway  gw-2        `,
			method: "GET",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []GatewayFormData{
					{
						GatewayId: "gw-1",
						DisplayName: "janis's Gateway",
					},
					{
						GatewayId: "gw-2",
						DisplayName: "janis's Gateway",
					},
				})
			},
		},
		{
			name:   "no gateways",
			args:   []string{"list"},
			expect: "You have no gateways",
			method: "GET",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []string{})
			},
		},
		{
			name:   "new license plate name",
			args:   []string{"create", "--generate"},
			expect: "ns-qwerty",
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"gatewayId": "ns-qwerty",
				})
			},
		},
		{
			name:   "new license plate with display name",
			args:   []string{"create", "--generate", "--display-name", "my display name"},
			expect: "ns-qwerty",
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"gatewayId": "ns-qwerty",
				})
			},
		},
		{
			name:   "new name",
			args:   []string{"create", "--gateway-id", "ns-sampler", "--display-name", "my display name"},
			expect: "ns-sampler",
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"gatewayId": "ns-sampler",
				})
			},
		},
		{
			name: "new name fails",
			args: []string{"create", "--gateway-id", "ns"},
			expect: heredoc.Doc(`
        Error: Validation Failed
        Gateway name must be between 5 and 15 alpha-numeric lowercase characters and start and end with an alphabet.
      `),
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(500, map[string]interface{}{
					"message": "Validation Failed",
					"details": map[string]interface{}{
						"d0": map[string]interface{}{
							"message": "Gateway name must be between 5 and 15 alpha-numeric lowercase characters and start and end with an alphabet.",
						},
					},
				})
			},
		},
		{
			name:   "new gateway fails",
			args:   []string{"create", "--generate"},
			expect: "Error: Validation Failed\nYou do not have access to this resource",
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(400, map[string]interface{}{
					"error":             "Validation Failed",
					"error_description": "You do not have access to this resource",
				})
			},
		},
		{
			name:    "destroy gateway",
			args:    []string{"destroy"},
			expect:  "Gateway destroyed: ns-sampler",
			method:  "DELETE",
			gateway: "/ns-sampler",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{})
			},
		},
		{
			name:   "show current gateway",
			args:   []string{"current"},
			expect: `Display Name     Gateway ID  
janis's Gateway  ns-sampler  
`,
			method: "GET",
			gateway: "/ns-sampler",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"displayName": "janis's Gateway",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			dir := t.TempDir()
			setup(dir)
			if tt.response != nil {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				URL := fmt.Sprintf("https://api.gov.ca/ds/api/v3/gateways%s", tt.gateway)
				httpmock.RegisterResponder(tt.method, URL, tt.response)
			}
			ctx := &pkg.AppContext{
				ApiHost:    "api.gov.ca",
				ApiVersion: "v3",
				Gateway:    "ns-sampler",
			}
			args := append([]string{"gateway"}, tt.args...)
			mainCmd := &cobra.Command{
				Use:          "gwa",
				SilenceUsage: true,
			}
			mainCmd.AddCommand(NewGatewayCmd(ctx, buf))
			mainCmd.SetArgs(args)

			// Use buffer to capture table output
			if (tt.name == "list gateways" || tt.name == "show current gateway") {
				mainCmd.Execute()
				out := buf.String()
				assert.Contains(t, out, tt.expect, "Expect: %v\nActual: %v\n", tt.expect, out)
			} else {
				out := capturer.CaptureOutput(func() {
					mainCmd.Execute()
				})
				assert.Contains(t, out, tt.expect, "Expect: %v\nActual: %v\n", tt.expect, out)
			}
		})
	}
}
