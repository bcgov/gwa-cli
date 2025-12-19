package cmd

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func TestGatewayPatternCommands(t *testing.T) {
	tests := []struct {
		name     string
		expect   string
		method   string
		gateway  string
		payload  string
		response httpmock.Responder
	}{
		{
			name: "eval gateway pattern",
			expect: `kind: GatewayService
name: sdx.test-abc
`,
			method: "PUT",
			payload: `{
				"pattern": "simple-service.r1",
				"parameters": {
					"gateway_id": "gw-1",
					"service_name": "test-abc",
					"service_url": "https://httpbun.com"
				}
			}`,
			gateway: "/ns-sampler/pattern",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, GatewayPatternResponse{
					Documents: []interface{}{
						map[string]interface{}{
							"kind": "GatewayService",
							"name": "sdx.test-abc",
						},
					},
				},
				)
			},
		},

		{
			name:   "eval missing parameter",
			expect: `Error: Invalid input`,
			method: "PUT",
			payload: `{
				"pattern": "simple-service.r1",
				"parameters": {
					"gateway_id": "gw-1",
					"service_url": "https://httpbun.com"
				}
			}`,
			gateway: "/ns-sampler/pattern",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(422, map[string]interface{}{
					"message": "Invalid input",
					"fields": map[string]interface{}{
						"service_name": map[string]interface{}{
							"message": "service_name is required",
						},
					},
				},
				)
			},
		},

		{
			name:   "eval invalid gateway pattern",
			expect: `Error: Invalid input`,
			method: "PUT",
			payload: `{
				"pattern": "simple-service.r1",
				"parameters": {
					"gateway_id": "gw-1",
					"service_url": "https://httpbun.com"
				}
			}`,
			gateway: "/ns-sampler/pattern",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(422, map[string]interface{}{
					"message": "Invalid input",
					"fields": map[string]interface{}{
						"pattern": map[string]interface{}{
							"message": "pattern not found",
						},
					},
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			setup(dir)

			// prepare pattern.yaml
			patternFilePath := fmt.Sprintf("%s/pattern.yaml", dir)
			err := os.WriteFile(patternFilePath, []byte(tt.payload), 0644)
			if err != nil {
				t.Fatalf("failed to write pattern file: %v", err)
			}

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
			args := append([]string{"gateway-pattern"}, patternFilePath)
			mainCmd := &cobra.Command{
				Use:          "gwa",
				SilenceUsage: true,
			}
			mainCmd.AddCommand(GatewayPatternCmd(ctx))
			mainCmd.SetArgs(args)

			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})
			assert.Contains(t, out, tt.expect, "Expect: %v\nActual: %v\n", tt.expect, out)
		})
	}
}
