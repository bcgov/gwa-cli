package pkg

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var host = "api.gov.bc.ca"
var tokenUrl = "https://api.gov.bc.ca/auth/token"

// TODO: The tmp dir and viper aren't getting along in this test, try and investigate
// for some more reliable results. The happy path is covered in the meantime
func SetupAuthConfig(dir string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName(".gwa-config")
	viper.SetConfigType("yaml")
	viper.SafeWriteConfig()
	err := viper.ReadInConfig()
	return err
}

func TestFetchConfigUrl(t *testing.T) {
	tests := []struct {
		name      string
		expect    string
		responder httpmock.Responder
	}{
		{
			name:   "success",
			expect: "/ds/api/v2/openapi.yaml",
			responder: func(_ *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(204, "")
				res.Header.Set("link", `</ds/api/v2/openapi.yaml>; rel="service-desc"`)
				return res, nil
			},
		},
		{
			name:   "failure",
			expect: "host is not configured correctly",
			responder: func(_ *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, "")
				return res, nil
			},
		},
		{
			name:   "failure",
			expect: "unable to find OpenAPI configuration",
			responder: func(_ *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(204, "")
				return res, nil
			},
		},
	}
	for _, tt := range tests {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		t.Run(tt.name, func(t *testing.T) {

			url := fmt.Sprintf("https://%s/ds/api", host)
			httpmock.RegisterResponder("GET", url, tt.responder)
			ctx := &AppContext{
				Host: host,
			}
			link, err := fetchConfigUrl(ctx)
			if err != nil {
				assert.ErrorContains(t, err, tt.expect)
			} else {
				assert.Equal(t, tt.expect, link)
			}
		})
	}
}

func TestParseLinkHeader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expect  string
		isError bool
	}{
		{
			name:   "happy",
			input:  `</ds/api/v2/openapi.yaml>; rel="service-desc"`,
			expect: "/ds/api/v2/openapi.yaml",
		},
		{
			name:    "happy",
			input:   `rel="service-desc"`,
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseLinkHeader(tt.input)
			if tt.isError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expect, result)

			}
		})
	}
}

func TestFetchOpenAPIConfig(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s/ds/api/v2/openapi.yaml", host)
	expect := fmt.Sprintf("https://authz-%s/auth/realms/aps-v2/.well-known/openid-configuration", host)

	httpmock.RegisterResponder("GET", url, func(_ *http.Request) (*http.Response, error) {
		yamlResult := fmt.Sprintf(`components:
  securitySchemes:
    openid:
      type: openIdConnect
      description: OIDC Login
      openIdConnectUrl: >-
        %s`, expect)
		return httpmock.NewStringResponse(200, yamlResult), nil
	})
	ctx := &AppContext{
		Host: host,
	}
	result, err := fetchOpenApiConfig(ctx, "/ds/api/v2/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expect, result)
}

func TestDeviceLogin(t *testing.T) {
	dir := t.TempDir()
	SetupAuthConfig(dir)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	clientId := "client123"
	wellKnownConfig := WellKnownConfig{
		DeviceAuthorizationEndpoint: fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/auth/device", host),
		TokenEndpoint:               fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/token", host),
	}
	verificationUrl := fmt.Sprintf("https://authz-%s/auth/realms/app/device", host)
	deviceCode := "1q2w3e4r"
	httpmock.RegisterResponder("POST", wellKnownConfig.DeviceAuthorizationEndpoint, func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, clientId, r.PostFormValue("client_id"))
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"device_code":               deviceCode,
			"user_code":                 "ABCD-EFGH",
			"verification_uri":          verificationUrl,
			"verification_uri_complete": verificationUrl + "?user_code=1q2w3e4r",
			"expires_in":                600,
			"interval":                  5,
		})
	})
	count := 2
	httpmock.RegisterResponder("POST", wellKnownConfig.TokenEndpoint, func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, deviceCode, r.FormValue("device_code"))
		assert.Equal(t, clientId, r.FormValue("client_id"))
		fmt.Println("count", count)
		if count > 0 {
			count -= 1
			return httpmock.NewStringResponse(401, ""), nil
		}
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"access_token":  "q1w2e3r4t5",
			"refresh_token": "y6u7i8o9p0",
		})
	})
	deviceLogin(wellKnownConfig, clientId, 0)
	assert.Equal(t, "q1w2e3r4t5", viper.GetString("api_key"))
	assert.Equal(t, "y6u7i8o9p0", viper.GetString("refresh_token"))
}

func TestFetchWellKnown(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := fmt.Sprintf("https://%s/.well-known/openid-configuration", host)
	httpmock.RegisterResponder("GET", url, func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"device_authorization_endpoint": fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/auth/device", host),
			"token_endpoint":                fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/token", host),
		})
	})
	result, err := fetchWellKnown(url)
	assert.NoError(t, err)
	assert.Equal(t, WellKnownConfig{
		DeviceAuthorizationEndpoint: fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/auth/device", host),
		TokenEndpoint:               fmt.Sprintf("https://authz-%s/auth/realms/app/protocol/openid-connect/token", host),
	}, result)
}

func TestPollAuthStatus(t *testing.T) {
	dir := t.TempDir()
	SetupAuthConfig(dir)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	i := 0
	url := fmt.Sprintf("https://%s/auth/token", host)
	httpmock.RegisterResponder("POST", url, func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "client123", r.FormValue("client_id"))
		assert.Equal(t, "ABCD-EFGH", r.FormValue("device_code"))
		assert.Equal(t, "urn:ietf:params:oauth:grant-type:device_code", r.FormValue("grant_type"))

		if i == 2 {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"access_token":       "q1w2e3r4t5y6",
				"refresh_expires_in": 300,
				"refresh_token":      "r5t6y7u8i9o0",
			})
		}
		i += 1
		return httpmock.NewJsonResponse(401, "")
	})
	pollAuthStatus(url, "client123", "ABCD-EFGH")
	if i > 2 {
		// assert.Error(t, err)
		// } else {
		// assert.NoError(t, err)
		assert.Equal(t, "q1w2e3r4t5y6", viper.GetString("api_key"))
		assert.Equal(t, "r5t6y7u8i9o0", viper.GetString("refresh_token"))
	}
}

func TestRefreshToken(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name      string
		expect    string
		responder httpmock.Responder
	}{
		{
			name:   "success",
			expect: "",
			responder: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"access_token":       "q1w2e3r4t5y6",
					"refresh_expires_in": 300,
					"refresh_token":      "r5t6y7u8i9o0",
				})
			},
		},
		{
			name:   "unauthorized",
			expect: "nono",
			responder: func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(401, map[string]interface{}{
					"error": "Unauthorized",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			err := RefreshToken(ctx)
			SetupAuthConfig(dir)

			if tt.expect != "" {
				assert.Error(t, err)
			}

		})
	}
}

func TestClientCredentialLogin(t *testing.T) {
	dir := t.TempDir()
	SetupAuthConfig(dir)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	apiKey := "12ad54kl"
	refreshToken := "34fg78hj"
	httpmock.RegisterResponder("POST", tokenUrl, func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "client123", r.FormValue("client_id"))
		assert.Equal(t, "$3cr3t", r.FormValue("client_secret"))
		assert.Equal(t, "client_credentials", r.FormValue("grant_type"))

		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"access_token":       apiKey,
			"refresh_token":      refreshToken,
			"refresh_expires_in": 0,
		})
	})
	ClientCredentialLogin(tokenUrl, "client123", "$3cr3t")

	assert.Equal(t, viper.GetString("api_key"), apiKey)
	assert.Equal(t, viper.GetString("refresh_token"), refreshToken)
}
