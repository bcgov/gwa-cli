package cmd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func TestStatusCmds(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expect      string
		noNamespace bool
		response    httpmock.Responder
	}{
		{
			name:   "no services",
			expect: "You currently do not have any services",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []map[string]interface{}{})
			},
		},
		{
			name:        "no namespace",
			expect:      "no namespace has been defined",
			noNamespace: true,
		},
		{
			name:   "prints json",
			args:   []string{"--json"},
			expect: `[{"name":"my-awesome-service","upstream":"upstream.host.com","status":"UP","reason":"No reason at all","host":"host.com","env_host":"host.com"}]`,
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []map[string]interface{}{
					{
						"name":     "my-awesome-service",
						"upstream": "upstream.host.com",
						"status":   "UP",
						"reason":   "No reason at all",
						"env_host": "host.com",
						"host":     "host.com",
					},
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			host := "api.gov.bc.ca"
			URL := fmt.Sprintf("https://%s/gw/api/namespaces/ns-sampler/services", host)
			httpmock.RegisterResponder("GET", URL, tt.response)

			args := append([]string{"status"}, tt.args...)
			ctx := &pkg.AppContext{
				Namespace: "ns-sampler",
				ApiHost:   host,
			}

			if tt.noNamespace {
				ctx.Namespace = ""
			}

			mainCmd := &cobra.Command{
				Use: "gwa",
			}
			mainCmd.AddCommand(NewStatusCmd(ctx))
			mainCmd.SetArgs(args)
			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})
			assert.Contains(t, out, tt.expect)
		})
	}
}
