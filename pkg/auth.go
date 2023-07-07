package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

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

	err = deviceLogin(wellKnownConfig, ctx.ClientId)
	if err != nil {
		return err
	}

	return nil
}

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

	err = clientCredentialLogin(wellKnownConfig, clientId, clientSecret)
	if err != nil {
		return err
	}

	return nil
}

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
	linkHeader := response.Header.Get("Link")
	result, err := parseLinkHeader(linkHeader)
	if err != nil {
		return "", err
	}

	return result, nil
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
		return result, errors.New("unable to find OpenAPI")
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

func deviceLogin(wellKnownConfig WellKnownConfig, clientId string) error {
	data := url.Values{}
	data.Set("client_id", clientId)
	URL := wellKnownConfig.DeviceAuthorizationEndpoint
	request, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var device DeviceData
		json.Unmarshal(b, &device)
		urlLine := fmt.Sprintf("\n\nPlease sign in at %s", device.VerificationUri)
		fmt.Println(urlLine)
		fmt.Println("Input the following code", device.UserCode)
		fmt.Println("\nWaiting for authentication handshake...")

		for i := 0; i < 60; i++ {
			auth, err := pollAuthStatus(wellKnownConfig.TokenEndpoint, clientId, device.DeviceCode)
			if err == nil {
				viper.Set("api_key", auth.AccessToken)
				viper.Set("refresh_token", auth.RefreshToken)
				viper.Set("refresh_expires_in", auth.RefreshExpiresIn)
				viper.WriteConfig()
				return nil
			}
			time.Sleep(time.Second * 7)
		}
		return errors.New("login request timed out")
	} else {
		return errors.New(response.Status)
	}
}

type WellKnownConfig struct {
	ClientCredentials           string `json:"client_credentials"`
	TokenEndpoint               string `json:"token_endpoint"`
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`
}

func fetchWellKnown(url string) (WellKnownConfig, error) {
	var wellKnownConfig WellKnownConfig
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return wellKnownConfig, err
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return wellKnownConfig, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return wellKnownConfig, err
	}

	json.Unmarshal(body, &wellKnownConfig)

	return wellKnownConfig, err
}

func pollAuthStatus(URL string, clientId string, deviceCode string) (TokenResponse, error) {
	var tokenData TokenResponse
	data := url.Values{}
	data.Set("device_code", deviceCode)
	data.Set("client_id", clientId)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	request, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accepts", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return tokenData, err
		}

		var errorResult TokenRequestError
		json.Unmarshal(b, &errorResult)
		return tokenData, errors.New(errorResult.ErrorDescription)
	}

	if response.StatusCode == http.StatusOK {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(b, &tokenData)
		return tokenData, nil
	}
	return tokenData, nil
}

type CredentialError struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func clientCredentialLogin(wellKnownConfig WellKnownConfig, clientId string, clientSecret string) error {
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secretn", clientSecret)
	data.Set("grant_type", "client_credentials")
	request, err := http.NewRequest(http.MethodPost, wellKnownConfig.TokenEndpoint, strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var auth TokenResponse
		json.Unmarshal(b, &auth)
		viper.Set("api_key", auth.AccessToken)
		viper.Set("refresh_token", auth.RefreshToken)
		viper.Set("refresh_expires_in", auth.RefreshExpiresIn)
		viper.WriteConfig()
		return nil
	}
	var errorMessage CredentialError
	b, err := io.ReadAll(response.Body)
	json.Unmarshal(b, &errorMessage)

	if err != nil {
		return err
	}
	return errors.New(errorMessage.Description)
}
