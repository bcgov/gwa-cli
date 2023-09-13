package pkg

import (
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

type AppContext struct {
	ApiHost    string
	ApiKey     string
	Auth       AuthDetails
	AppVersion int
	ClientId   string
	Cwd        string
	Host       string
	Namespace  string
	Scheme     string
	Version    string
}

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
