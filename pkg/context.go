package pkg

import (
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

// Global application context properties
type AppContext struct {
	ApiHost    string      // Stores the API host, which is embedded at compile time.
	ApiKey     string      // API key stored in config from `login` command.
	Auth       AuthDetails // Embedded authentication details
	ApiVersion string      // Stores the currently supported API version, emdedded at compile time.
	ClientId   string      // Stores client ID, embedded at compile time.
	Cwd        string      // Convenience property to access the current working directory at runtime.
	Host       string      // Host is variable and extracted from config.
	Namespace  string      // Namespace is variable and extracted from config.
	Scheme     string      // Scheme, defaults to https, extracted from config.
	Version    string      // Version, embedded at compile time.
}

// CreateUrl is a standardized method to compose APS specific API URLs.
// Example URL:
//
//	params := struct {
//	   Hello string `url:"hello"`
//	}{
//	  Hello: "world",
//	}
//	url, err := tt.ctx.CreateUrl("/status", params)
//
// which would return `https://api.gov.bc.ca/status?hello=world`
func (a *AppContext) CreateUrl(path string, params interface{}) (string, error) {
	q, err := query.Values(params)
	if err != nil {
		return "", err
	}
	queryString := q.Encode()

	if a.ApiHost == "" && a.Host == "" {
		return "", fmt.Errorf("no host set")
	}

	host := a.ApiHost
	if a.Host != "" {
		host = a.Host
	}

	scheme := "https"

	if a.Scheme != "" {
		scheme = a.Scheme
	}

	url := url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: queryString,
	}

	return url.String(), nil
}
