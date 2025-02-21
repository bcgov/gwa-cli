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

func DeviceLogin(ctx *AppContext) error {
	Info("Auth: Device login selected")
	openApiPathname, err := fetchConfigUrl(ctx)
	if err != nil {
		return err
	}
	Info("OpenAPI Pathname received")

	authTokenUrl, err := fetchOpenApiConfig(ctx, openApiPathname)
	if err != nil {
		return err
	}
	Info("Auth token recieved")

	wellKnownConfig, err := fetchWellKnown(authTokenUrl)
	if err != nil {
		return err
	}
	Info("Well known config received")
	viper.Set("token_endpoint", wellKnownConfig.TokenEndpoint)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	Info("Config updated")

	err = deviceLogin(wellKnownConfig, ctx.ClientId, 8)
	if err != nil {
		return err
	}
	Info("Logged in")

	return nil
}

func ClientCredentialsLogin(ctx *AppContext, clientId string, clientSecret string) error {
	Info("Auth method: Client Credential")
	openApiPathname, err := fetchConfigUrl(ctx)
	if err != nil {
		return err
	}
	Info("OpenAPI Pathname received")

	authTokenUrl, err := fetchOpenApiConfig(ctx, openApiPathname)
	if err != nil {
		return err
	}
	Info("Auth token recieved")

	wellKnownConfig, err := fetchWellKnown(authTokenUrl)
	if err != nil {
		return err
	}
	Info("Well known config received")
	viper.Set("token_endpoint", wellKnownConfig.TokenEndpoint)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	Info("Config updated")

	err = ClientCredentialLogin(wellKnownConfig.TokenEndpoint, clientId, clientSecret)
	if err != nil {
		return err
	}
	Info("Logged in")

	return nil
}

func fetchConfigUrl(ctx *AppContext) (string, error) {
	client := http.Client{}
	URL, _ := ctx.CreateUrl("/ds/api", nil)
	Info(fmt.Sprintf("Config URL: %s", URL))
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

func fetchOpenApiConfig(ctx *AppContext, openApiPathname string) (string, error) {
	client := http.Client{}
	URL, _ := ctx.CreateUrl(openApiPathname, nil)
	Info(fmt.Sprintf("Config URL: %s", URL))
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
	Info(fmt.Sprintf("Config URL: %s", URL))

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

	fmt.Println("\nTo complete the login process, please follow these steps:")
	fmt.Printf("1. Open this URL in your web browser: %s\n", boldText.Render(response.Data.VerificationUri))
	fmt.Printf("2. Enter this code when prompted: %s\n", boldText.Render(response.Data.UserCode))
	fmt.Println("3. Complete the IDIR authentication process in your browser")
	fmt.Println("4. Return to this terminal window after successful authentication")
	fmt.Println("\nWaiting for you to complete the login process...")

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

func RefreshToken(ctx *AppContext) error {
	Info("Auth: Refreshing token")
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
