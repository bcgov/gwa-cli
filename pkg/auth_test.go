package pkg

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var tokenUrl = "https://api.gov.bc.ca/auth/token"

func SetupConfig(dir string) {
	viper.AddConfigPath(dir)
	viper.SetConfigName(".testing")
	viper.SetConfigType("yaml")
	viper.SetDefault("namespace", "abc")
	viper.SafeWriteConfig()
}

func TestFetchConfigUrl(t *testing.T) {
	t.Skip()
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
	t.Skip()
}

func TestDeviceLogin(t *testing.T) {
	t.Skip()
}

func TestFetchWellKnown(t *testing.T) {
	t.Skip()
}

func TestPollAuthStatus(t *testing.T) {
	t.Skip()
}

func TestRefreshToken(t *testing.T) {
	t.Skip()
}

func TestClientCredentialLogin(t *testing.T) {
	dir := t.TempDir()
	SetupConfig(dir)
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
	err := clientCredentialLogin(tokenUrl, "client123", "$3cr3t")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, viper.GetString("api_key"), apiKey)
	assert.Equal(t, viper.GetString("refresh_token"), refreshToken)
}
