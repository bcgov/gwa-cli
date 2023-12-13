package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var boldText = lipgloss.NewStyle().Bold(true)

// DeviceLogin handles the default authentication flow
func DeviceLogin(ctx *AppContext) error {
	openApiPathname, err := fetchConfigUrl(ctx)
	if err != nil {
		return err
	}

	authTokenUrl, err := fetchOpenApiConfig(ctx, openApiPathname)
	if err != nil {
		return err
	}

	wellKnownConfig, err := fetchWellKnown(authTokenUrl)
	if err != nil {
		return err
	}
	viper.Set("token_endpoint", wellKnownConfig.TokenEndpoint)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	err = deviceLogin(wellKnownConfig, ctx.ClientId, 8)
	if err != nil {
		return err
	}

	return nil
}

// ClientCredentialsLogin handles where the user provides clientId and clientSecret
func ClientCredentialsLogin(ctx *AppContext, clientId string, clientSecret string) error {
	openApiPathname, err := fetchConfigUrl(ctx)
	if err != nil {
		return err
	}

	authTokenUrl, err := fetchOpenApiConfig(ctx, openApiPathname)
	if err != nil {
		return err
	}

	wellKnownConfig, err := fetchWellKnown(authTokenUrl)
	if err != nil {
		return err
	}
	viper.Set("token_endpoint", wellKnownConfig.TokenEndpoint)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	err = ClientCredentialLogin(wellKnownConfig.TokenEndpoint, clientId, clientSecret)
	if err != nil {
		return err
	}

	return nil
}

// fetchConfigUrl is the first step in an authentication request, it retreives the URL from the
// portal's link header. This prevents the need to embed the config URL in the binary.
func fetchConfigUrl(ctx *AppContext) (string, error) {
	client := http.Client{}
	URL, _ := ctx.CreateUrl("/ds/api", nil)
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusNoContent {
		linkHeader := response.Header.Get("Link")
		result, err := parseLinkHeader(linkHeader)
		if err != nil {
			return "", err
		}

		return result, nil
	}

	return "", fmt.Errorf("host is not configured correctly")
}

// parseLinkHeader takes a link header string and returns the embedded URL
func parseLinkHeader(link string) (string, error) {
	var result string
	links := strings.Split(link, ",")

	for _, link := range links {
		segments := strings.Split(link, ";")
		if len(segments) < 2 {
			continue
		}

		urlSegment := strings.TrimSpace(segments[0])
		url := strings.Trim(urlSegment, "<>")
		result = url
	}

	if result == "" {
		return result, fmt.Errorf("unable to find OpenAPI configuration")
	}

	return result, nil
}

type TokenRequestError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int16  `json:"expires_in"`
	RefreshExpiresIn int32  `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type OpenApi struct {
	Components struct {
		SecuritySchemes struct {
			OpenId struct {
				OpenIdConnectUrl string `yaml:"openIdConnectUrl"`
			} `yaml:"openid"`
			Oauth2 struct {
				Flows struct {
					ClientCredentials struct {
						TokenUrl string `yaml:"tokenUrl"`
					} `yaml:"clientCredentials"`
				} `yaml:"flows"`
			} `yaml:"oauth2"`
		} `yaml:"securitySchemes"`
	} `yaml:"components"`
}

// fetchOpenApiConfig is the second request in the authentication toolchain. It parses the
// config YAML doc and returns the necessary properties to continue to step 3.
func fetchOpenApiConfig(ctx *AppContext, openApiPathname string) (string, error) {
	client := http.Client{}
	URL, _ := ctx.CreateUrl(openApiPathname, nil)
	request, err := http.NewRequest(http.MethodGet, URL, nil)

	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "text/yaml")
	response, err := client.Do(request)

	if err != nil {
		return "", err
	}

	var openApiConfig OpenApi
	body, err := io.ReadAll(response.Body)
	yaml.Unmarshal(body, &openApiConfig)
	if err != nil {
		return "", err
	}

	return openApiConfig.Components.SecuritySchemes.OpenId.OpenIdConnectUrl, nil
}

type AuthDetails struct {
	Token string
}

type DeviceData struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationUri         string `json:"verification_uri"`
	VerificationUriComplete string `json:"verification_uri_complete"`
}

// deviceLogin is the final step in the device authentication flow
func deviceLogin(wellKnownConfig WellKnownConfig, clientId string, timeout time.Duration) error {
	data := url.Values{}
	data.Set("client_id", clientId)
	URL := wellKnownConfig.DeviceAuthorizationEndpoint
	request, err := NewApiPost[DeviceData](&AppContext{}, URL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	request.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := request.Do()
	if err != nil {
		return err
	}

	urlLine := fmt.Sprintf("\n\nPlease sign in at %s", response.Data.VerificationUri)
	fmt.Println(urlLine)
	fmt.Println("Input the following code", boldText.Render(response.Data.UserCode))
	fmt.Print("\nWaiting for authentication handshake...")

	for i := 0; i < 60; i++ {
		err := pollAuthStatus(wellKnownConfig.TokenEndpoint, clientId, response.Data.DeviceCode)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * timeout)
	}

	return fmt.Errorf("login request timed out")
}

type WellKnownConfig struct {
	ClientCredentials           string `json:"client_credentials"`
	TokenEndpoint               string `json:"token_endpoint"`
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`
}

// fetchWellKnown is the 3rd request in the auth flow
func fetchWellKnown(url string) (WellKnownConfig, error) {
	request, err := NewApiGet[WellKnownConfig](&AppContext{}, url)
	if err != nil {
		return WellKnownConfig{}, err
	}
	response, err := request.Do()
	if err != nil {
		return WellKnownConfig{}, err
	}

	return response.Data, nil
}

// pollAuthStatus can be called during a device-based login, waiting for the
// server to return a success code when the session is active.
func pollAuthStatus(URL string, clientId string, deviceCode string) error {
	data := url.Values{}
	data.Set("device_code", deviceCode)
	data.Set("client_id", clientId)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	request, err := NewApiPost[TokenResponse](&AppContext{}, URL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	request.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Request.Header.Set("Accepts", "application/json")

	response, err := request.Do()
	if err != nil {
		return err
	}

	SaveConfig(&response.Data)
	return nil
}

// RefreshToken handles all the logic to attempt a token refresh. Returns an error
// if the token has expired.
func RefreshToken(ctx *AppContext) error {
	tokenEndpoint := viper.GetString("token_endpoint")
	refreshToken := viper.GetString("refresh_token")

	if refreshToken == "" {
		return nil
	}

	data := make(url.Values)
	data.Set("client_id", ctx.ClientId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	request, err := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var data TokenResponse
		json.Unmarshal(body, &data)

		ctx.ApiKey = data.AccessToken
		viper.Set("api_key", data.AccessToken)
		viper.Set("refresh_token", data.RefreshToken)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}

		return nil
	}

	var errorResponse ApiErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return fmt.Errorf(string(body))
	}

	return errorResponse.GetError()
}

// ClientCredentialLogin runs the acutal client credential login. It's the final step in the
// authentication flow.
func ClientCredentialLogin(tokenEndpoint string, clientId string, clientSecret string) error {
	data := make(url.Values)
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")
	ctx := &AppContext{}
	r, err := NewApiPost[TokenResponse](ctx, tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	r.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := r.Do()
	if err != nil {
		return err
	}

	return SaveConfig(&response.Data)
}

// SaveConfig writes all the important config values locally.
func SaveConfig(data *TokenResponse) error {
	viper.Set("api_key", data.AccessToken)
	viper.Set("refresh_token", data.RefreshToken)
	viper.Set("refresh_expires_in", data.RefreshExpiresIn)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
