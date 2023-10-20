package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"

	"github.com/bcgov/gwa-cli/pkg"
)

var apiHost = "api.gov.bc.ca"
var authHost = "https://api.gov.bc.ca/auth/token"

func setupConfig(dir string) error {
	fileName := ".gwa-config.yaml"
	path := path.Join(dir, fileName)
	configFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer configFile.Close()
	viper.AddConfigPath(dir)
	viper.SetConfigFile(path)
	viper.SetDefault("token_endpoint", authHost)
	return nil
}

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
			"license_title": "Access Only",
			"name":          "a-unit-test-dataset",
			"notes":         "Used for unit tests :)",
			"organization": map[string]interface{}{
				"name":  "ministry-of-citizen-services",
				"title": "Ministry of Citizen Services",
			},
			"organizationUnit": map[string]interface{}{
				"name":  "data-innovation-program",
				"title": "Data Innovation Program",
			},
			"products": map[string]interface{}{
				"environments": []map[string]interface{}{
					{"active": false,
						"flow": "client-credentials",
						"name": "dev"},
				},
				"id":   123,
				"name": "ABCDEF",
			},
			"record_publish_date": "2014-12-15",
			"security_class":      "LOW-PUBLIC",
			"tags":                []string{"BC", "Canada"},
			"title":               "A Unit Test Dataset",
			"view_audience":       "Government",
		},
	})

}

func issuersResponse(r *http.Request) (*http.Response, error) {
	return httpmock.NewJsonResponse(200, []map[string]interface{}{
		{
			"apiKeyName":          "X-API-KEY",
			"availableScopes":     []string{},
			"clientAuthenticator": "client-jwt-jwks-url",
			"clientMappers": []map[string]interface{}{
				{
					"defaultValue": "https://aps.gov.bc.ca",
					"name":         "audience",
				},
			},
			"clientRoles": []string{"read", "write"},
			"environmentDetails": []map[string]interface{}{
				{
					"clientId":           "aps-team",
					"clientRegistration": "managed",
					"clientSecret":       "****",
					"environment":        "dev",
					"exists":             true,
					"issuerUrl":          "https://aps.gov.bc.ca/auth/realms/issuer",
				},
			},
			"flow":           "client-credentials",
			"isShared":       false,
			"mode":           "auto",
			"name":           "APS IdP",
			"owner":          "janis@idir",
			"resourceScopes": []string{},
		},
	})
}

func productsResponse(r *http.Request) (*http.Response, error) {
	return httpmock.NewJsonResponse(200, []map[string]interface{}{
		{
			"name":  "DemoNet",
			"appId": "132QWE",
			"environments": []map[string]interface{}{
				{
					"name":     "dev",
					"active":   false,
					"approval": false,
					"flow":     "public",
					"appId":    "00000000",
				},
				{
					"name":     "prod",
					"active":   false,
					"approval": true,
					"flow":     "public",
					"appId":    "00000001",
				},
			},
		},
	})
}

func orgUnitsResponse(r *http.Request) (*http.Response, error) {
	return httpmock.NewJsonResponse(200, []map[string]interface{}{
		{
			"id":    "1",
			"name":  "planning-and-innovation-division",
			"title": "Planning and Innovation Division",
		},
	})
}

func TestGetCmdTables(t *testing.T) {
	setupConfig(t.TempDir())
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
				"Name                 Title",
				"a-unit-test-dataset  A Unit Test Dataset",
			},
			response: datasetsResponse,
		},
		{
			name: "get issuers",
			args: []string{"issuers"},
			expect: []string{
				"Name     Flow                Mode  Owner",
				"APS IdP  client-credentials  auto  janis@idir",
			},
			response: issuersResponse,
		},
		{
			name: "get products",
			args: []string{"products"},
			expect: []string{
				"Name     App ID  Environments",
				"DemoNet  132QWE  2",
			},
			response: productsResponse,
		},
		{
			name: "get org-units",
			args: []string{"org-units", "--org", "ministry-of-citizens-services"},
			expect: []string{
				"Name                              Title",
				"planning-and-innovation-division  Planning and Innovation Division",
			},
			response: orgUnitsResponse,
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
			httpmock.RegisterResponder("POST", authHost, func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"access_token":       "123ABC",
					"refresh_token":      "refresh",
					"refresh_expires_in": 0,
				})
			})
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
	setupConfig(t.TempDir())
	tests := []struct {
		name     string
		args     []string
		expect   string
		response httpmock.Responder
	}{
		{
			name: "get datasets json",
			args: []string{"datasets", "--json"},
			expect: `[{"license_title":"Access Only","name":"a-unit-test-dataset","notes":"Used for unit tests :)","organization":{"name":"ministry-of-citizen-services","title":"Ministry of Citizen Services"},"organizationUnit":{"name":"data-innovation-program","title":"Data Innovation Program"},"products":{"environments":[{"active":false,"flow":"client-credentials","name":"dev"}],"id":123,"name":"ABCDEF"},"record_publish_date":"2014-12-15","security_class":"LOW-PUBLIC","tags":["BC","Canada"],"title":"A Unit Test Dataset","view_audience":"Government"}]
`,
			response: datasetsResponse,
		},
		{
			name: "get datasets yaml",
			args: []string{"datasets", "--yaml"},
			expect: `- license_title: Access Only
  name: a-unit-test-dataset
  notes: Used for unit tests :)
  organization:
    name: ministry-of-citizen-services
    title: Ministry of Citizen Services
  organizationUnit:
    name: data-innovation-program
    title: Data Innovation Program
  products:
    environments:
        - active: false
          flow: client-credentials
          name: dev
    id: 123
    name: ABCDEF
  record_publish_date: "2014-12-15"
  security_class: LOW-PUBLIC
  tags:
    - BC
    - Canada
  title: A Unit Test Dataset
  view_audience: Government

`,
			response: datasetsResponse,
		},
		{
			name: "get issuers json",
			args: []string{"issuers", "--json"},
			expect: `[{"apiKeyName":"X-API-KEY","availableScopes":[],"clientAuthenticator":"client-jwt-jwks-url","clientMappers":[{"defaultValue":"https://aps.gov.bc.ca","name":"audience"}],"clientRoles":["read","write"],"environmentDetails":[{"clientId":"aps-team","clientRegistration":"managed","clientSecret":"****","environment":"dev","exists":true,"issuerUrl":"https://aps.gov.bc.ca/auth/realms/issuer"}],"flow":"client-credentials","isShared":false,"mode":"auto","name":"APS IdP","owner":"janis@idir","resourceScopes":[]}]
`,
			response: issuersResponse,
		},
		{
			name: "get issuers yaml",
			args: []string{"issuers", "--yaml"},
			expect: `- apiKeyName: X-API-KEY
  availableScopes: []
  clientAuthenticator: client-jwt-jwks-url
  clientMappers:
    - defaultValue: https://aps.gov.bc.ca
      name: audience
  clientRoles:
    - read
    - write
  environmentDetails:
    - clientId: aps-team
      clientRegistration: managed
      clientSecret: '****'
      environment: dev
      exists: true
      issuerUrl: https://aps.gov.bc.ca/auth/realms/issuer
  flow: client-credentials
  isShared: false
  mode: auto
  name: APS IdP
  owner: janis@idir
  resourceScopes: []

`,
			response: issuersResponse,
		},
		{
			name: "get products json",
			args: []string{"products", "--json"},
			expect: `[{"appId":"132QWE","environments":[{"active":false,"appId":"00000000","approval":false,"flow":"public","name":"dev"},{"active":false,"appId":"00000001","approval":true,"flow":"public","name":"prod"}],"name":"DemoNet"}]
`,
			response: productsResponse,
		},
		{
			name: "get products yaml",
			args: []string{"products", "--yaml"},
			expect: `- appId: 132QWE
  environments:
    - active: false
      appId: "00000000"
      approval: false
      flow: public
      name: dev
    - active: false
      appId: "00000001"
      approval: true
      flow: public
      name: prod
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
			httpmock.RegisterResponder("POST", authHost, func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"access_token":       "123ABC",
					"refresh_token":      "refresh",
					"refresh_expires_in": 0,
				})
			})
			httpmock.RegisterResponder("GET", URL, tt.response)
			mainCmd := setupGetTests(tt.args, tt.response, nil)
			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})
			assert.Equal(t, tt.expect, out)
		})
	}
}
