package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"

	"github.com/bcgov/gwa-cli/pkg"
)

var apiHost = "api.gov.bc.ca"

func setupGetTests(args []string, response httpmock.Responder, buf *bytes.Buffer) *cobra.Command {
	ctx := &pkg.AppContext{
		Namespace: "ns-sampler",
		ApiHost:   apiHost,
	}

	args = append([]string{"get"}, args...)

	mainCmd := &cobra.Command{
		Use: "gwa",
	}
	mainCmd.AddCommand(NewGetCmd(ctx, buf))
	mainCmd.SetArgs(args)
	return mainCmd
}

func datasetsResponse(r *http.Request) (*http.Response, error) {
	return httpmock.NewJsonResponse(200, []map[string]interface{}{
		{
			"name":  "full-dataset",
			"title": "Full Dataset",
		},
	})
}

func issuersResponse(r *http.Request) (*http.Response, error) {
	return httpmock.NewJsonResponse(200, []map[string]interface{}{
		{
			"name":  "Ministry IdP",
			"flow":  "client-credentials",
			"mode":  "auto",
			"owner": "janis@idir",
		},
	})
}

func productsResponse(r *http.Request) (*http.Response, error) {
	return httpmock.NewJsonResponse(200, []map[string]interface{}{
		{
			"name":  "DemoNet",
			"appId": "132QWE",
			"environments": []map[string]interface{}{
				{"id": "e1", "name": "dev"},
			},
		},
	})
}

func TestGetCmdTables(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expect   []string
		response httpmock.Responder
	}{
		{
			name: "get datasets",
			args: []string{"datasets"},
			expect: []string{
				"Name          Title",
				"full-dataset  Full Dataset",
			},
			response: datasetsResponse,
		},
		{
			name: "get issuers",
			args: []string{"issuers"},
			expect: []string{
				"Name          Flow                Mode  Owner",
				"Ministry IdP  client-credentials  auto  janis@idir",
			},
			response: issuersResponse,
		},
		{
			name: "get products",
			args: []string{"products"},
			expect: []string{
				"Name     AppId   Environments",
				"DemoNet  132QWE  1",
			},
			response: productsResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			var operator = tt.args[0]
			if operator == "datasets" {
				operator = "directory"
			}
			URL := fmt.Sprintf("https://%s/ds/api/v2/namespaces/ns-sampler/%s", host, operator)
			httpmock.RegisterResponder("GET", URL, tt.response)
			buf := &bytes.Buffer{}
			mainCmd := setupGetTests(tt.args, tt.response, buf)
			mainCmd.Execute()
			out := buf.String()
			for _, expected := range tt.expect {
				assert.Contains(t, out, expected)
			}
		})
	}
}

func TestGetJsonYamlCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expect   string
		response httpmock.Responder
	}{
		{
			name: "get datasets json",
			args: []string{"datasets", "--json"},
			expect: `[{"name":"full-dataset","title":"Full Dataset"}]
`,
			response: datasetsResponse,
		},
		{
			name: "get datasets yaml",
			args: []string{"datasets", "--yaml"},
			expect: `- name: full-dataset
  organization: ""
  tags: []
  title: Full Dataset

`,
			response: datasetsResponse,
		},
		{
			name: "get issuers json",
			args: []string{"issuers", "--json"},
			expect: `[{"name":"Ministry IdP","description":"","flow":"client-credentials","clientAuthenicator":"","mode":"auto","environmentDetails":null,"owner":"janis@idir"}]
`,
			response: issuersResponse,
		},
		{
			name: "get issuers yaml",
			args: []string{"issuers", "--yaml"},
			expect: `- name: Ministry IdP
  description: ""
  flow: client-credentials
  clientAuthenticator: ""
  mode: auto
  environmentDetails: []
  owner: janis@idir

`,
			response: issuersResponse,
		},
		{
			name: "get products json",
			args: []string{"products", "--json"},
			expect: `[{"appId":"132QWE","environments":[{"name":"dev"}],"name":"DemoNet"}]
`,
			response: productsResponse,
		},
		{
			name: "get products yaml",
			args: []string{"products", "--yaml"},
			expect: `- appId: 132QWE
  environments:
    - active: false
      appId: ""
      approval: false
      flow: ""
      name: dev
  name: DemoNet

`,
			response: productsResponse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			var operator = tt.args[0]
			if operator == "datasets" {
				operator = "directory"
			}
			URL := fmt.Sprintf("https://%s/ds/api/v2/namespaces/ns-sampler/%s", host, operator)
			httpmock.RegisterResponder("GET", URL, tt.response)
			mainCmd := setupGetTests(tt.args, tt.response, nil)
			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})
			assert.Equal(t, tt.expect, out)
		})
	}
}
