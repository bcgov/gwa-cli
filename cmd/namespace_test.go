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
		name      string
		args      []string
		expect    string
		namespace string
		method    string
		response  httpmock.Responder
	}{
		{
			name: "list namespaces",
			args: []string{"list"},
			expect: `ns-123
ns-456`,
			namespace: "ns-sampler",
			method:    "GET",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []string{
					"ns-123",
					"ns-456",
				})
			},
		},
		{
			name:      "no namespaces",
			args:      []string{"list"},
			expect:    "You have no namespaces",
			namespace: "ns-sampler",
			method:    "GET",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []string{})
			},
		},
		{
			name:      "new license plate name",
			args:      []string{"create"},
			expect:    "ns-qwerty",
			namespace: "ns-sampler",
			method:    "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"name": "ns-qwerty",
				})
			},
		},
		{
			name:      "new license plate with description",
			args:      []string{"create", "--description", "my description"},
			expect:    "ns-qwerty",
			namespace: "ns-sampler",
			method:    "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"name": "ns-qwerty",
				})
			},
		},
		{
			name:      "new name",
			args:      []string{"create", "--name", "ns-sampler", "--description", "my description"},
			expect:    "ns-sampler",
			namespace: "ns-sampler",
			method:    "POST",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"name": "ns-sampler",
				})
			},
		},
		{
			name:      "new namespace fails",
			args:      []string{"create"},
			expect:    "Error: Validation Failed\nYou do not have access to this resource",
			method:    "POST",
			namespace: "ns-sampler",
			response: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(400, map[string]interface{}{
					"error":             "Validation Failed",
					"error_description": "You do not have access to this resource",
				})
			},
		},
		{
			name:      "show current namespace",
			args:      []string{"current"},
			namespace: "ns-sampler",
			expect:    "ns-sampler",
		},
		{
			name:      "no current namespace to show",
			args:      []string{"current"},
			expect:    "no namespace has been defined",
			namespace: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.response != nil {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				httpmock.RegisterResponder(tt.method, "https://api.gov.ca/ds/api/v2/namespaces", tt.response)
			}
			ctx := &pkg.AppContext{
				ApiHost:   "api.gov.ca",
				Namespace: tt.namespace,
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
