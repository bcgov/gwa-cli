package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func createCmd(subCmd *cobra.Command, args []string) *cobra.Command {
	mainCmd := &cobra.Command{
		Use: "gwa",
	}
	mainArg := []string{"login"}
	mainCmd.AddCommand(subCmd)
	mainCmd.SetArgs(append(mainArg, args...))
	return mainCmd
}

var host = "api.gov.bc.ca"
var tokenUrl = "https://api.gov.bc.ca/auth/token"
var ctx = &pkg.AppContext{
	ApiHost: host,
}

func SetupLogin(dir string) error {
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

func TestDeviceLoginSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	dir := t.TempDir()
	SetupLogin(dir)

	// First request to obtain OpenAPI config
	linkUrl := fmt.Sprintf("https://%s/ds/api", host)
	httpmock.RegisterResponder("GET", linkUrl, func(_ *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(204, "")
		res.Header.Set("link", `</ds/api/v2/openapi.yaml>; rel="service-desc"`)
		return res, nil
	})
	// Fetch the URL for well-known OpenID Config
	openApiUrl := fmt.Sprintf("https://%s/ds/api/v2/openapi.yaml", host)
	httpmock.RegisterResponder("GET", openApiUrl, func(_ *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, `components:
    securitySchemes:
      openid:
        openIdConnectUrl: >-
          https://authz-api.gov.bc.ca/auth/realms/app/.well-known/openid-configuration`), nil
	})
	// Fetch the URL for well-known OpenID Config
	httpmock.RegisterResponder("GET", "https://authz-api.gov.bc.ca/auth/realms/app/.well-known/openid-configuration", func(_ *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"device_authorization_endpoint": fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/auth/device", host),
			"token_endpoint":                fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/token", host),
		})
	})
	// Grab the auth endpoints
	deviceAuthEndpointUrl := fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/auth/device", host)
	tokenEndpointUrl := fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/token", host)
	httpmock.RegisterResponder("POST", "https://authz-api.gov.bc.ca/auth/realms/app/.well-known/openid-configuration", func(_ *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"device_authorization_endpoint": deviceAuthEndpointUrl,
			"token_endpoint":                tokenEndpointUrl,
		})
	})
	// Loging URL
	verificationUrl := fmt.Sprintf("https://authz-%s/auth/realms/app/device", host)
	httpmock.RegisterResponder("POST", deviceAuthEndpointUrl, func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"device_code":               "1q2w3e4r",
			"user_code":                 "ABCD-EFGH",
			"verification_uri":          verificationUrl,
			"verification_uri_complete": verificationUrl + "?user_code=1q2w3e4r",
			"expires_in":                600,
			"interval":                  5,
		})
	})

	// Polling URL
	httpmock.RegisterResponder("POST", tokenEndpointUrl, func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"access_token":       "q1w2e3r4t5y6",
			"refresh_expires_in": 300,
			"refresh_token":      "r5t6y7u8i9o0",
		})
	})

	loginCmd := NewLoginCmd(ctx)
	mainCmd := createCmd(loginCmd, nil)
	out := capturer.CaptureOutput(func() {
		mainCmd.Execute()
	})
	assert.Contains(t, out, fmt.Sprintf(`

Please sign in at %s
Input the following code ABCD-EFGH

Waiting for authentication handshake...`, verificationUrl))
}

func TestClientCredentialLogin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	dir := t.TempDir()
	SetupLogin(dir)

	// First request to obtain OpenAPI config
	linkUrl := fmt.Sprintf("https://%s/ds/api", host)
	httpmock.RegisterResponder("GET", linkUrl, func(_ *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(204, "")
		res.Header.Set("link", `</ds/api/v2/openapi.yaml>; rel="service-desc"`)
		return res, nil
	})
	// Fetch the URL for well-known OpenID Config
	openApiUrl := fmt.Sprintf("https://%s/ds/api/v2/openapi.yaml", host)
	httpmock.RegisterResponder("GET", openApiUrl, func(_ *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, `components:
    securitySchemes:
      openid:
        openIdConnectUrl: >-
          https://authz-api.gov.bc.ca/auth/realms/app/.well-known/openid-configuration`), nil
	})
	// Fetch the URL for well-known OpenID Config
	tokenEndpointUrl := fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/token", host)
	httpmock.RegisterResponder("GET", "https://authz-api.gov.bc.ca/auth/realms/app/.well-known/openid-configuration", func(_ *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"token_endpoint": tokenEndpointUrl,
		})
	})
	//  Auth URL
	httpmock.RegisterResponder("POST", tokenEndpointUrl, func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"access_token":       "q1w2e3r4t5y6",
			"refresh_expires_in": 300,
			"refresh_token":      "r5t6y7u8i9o0",
		})
	})

	loginCmd := NewLoginCmd(ctx)
	mainCmd := createCmd(loginCmd, []string{"--client-id", "client123", "--client-secret", "$3cr3t"})
	out := capturer.CaptureOutput(func() {
		mainCmd.Execute()
	})
	assert.Contains(t, out, "Successfully logged in")
}
