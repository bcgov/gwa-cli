package cmd

import (
	// "encoding/json"
	"net/http"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func TestNamespaceCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expect   string
		method   string
		response httpmock.Responder
	}{
		{
			name: "list namespaces",
			args: []string{"list"},
			expect: `ns-123
ns-456`,
			method: "GET",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []string{
					"ns-123",
					"ns-456",
				})
			},
		},
		{
			name:   "no namespaces",
			args:   []string{"list"},
			expect: "You have no namespaces",
			method: "GET",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []string{})
			},
		},
		{
			name:   "new license plate name",
			args:   []string{"create"},
			expect: "ns-qwerty",
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"name": "ns-qwerty",
				})
			},
		},
		{
			name:   "new license plate with description",
			args:   []string{"create", "--description", "my description"},
			expect: "ns-qwerty",
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"name": "ns-qwerty",
				})
			},
		},
		{
			name:   "new name",
			args:   []string{"create", "--name", "ns-sampler", "--description", "my description"},
			expect: "ns-sampler",
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"name": "ns-sampler",
				})
			},
		},
		{
			name:   "new namespace fails",
			args:   []string{"create"},
			expect: "Error: Validation Failed\nYou do not have access to this resource",
			method: "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(400, map[string]interface{}{
					"error":             "Validation Failed",
					"error_description": "You do not have access to this resource",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(tt.method, "https://api.gov.ca/ds/api/v2/namespaces", tt.response)
			ctx := &pkg.AppContext{
				ApiHost:   "api.gov.ca",
				Namespace: "ns-sampler",
			}
			args := append([]string{"namespace"}, tt.args...)
			mainCmd := &cobra.Command{
				Use: "gwa",
			}
			mainCmd.AddCommand(NewNamespaceCmd(ctx))
			mainCmd.SetArgs(args)
			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})

			assert.Contains(t, out, tt.expect, "Expect: %v\nActual: %v\n", tt.expect, out)
		})
	}

}
