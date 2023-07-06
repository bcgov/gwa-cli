package pkg

import (
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
}

func (a *AppContext) CreateUrl(path string, params interface{}) (string, error) {
	q, err := query.Values(params)
	if err != nil {
		return "", err
	}
	queryString := q.Encode()

	url := url.URL{
		Scheme:   "https",
		Host:     a.ApiHost,
		Path:     path,
		RawQuery: queryString,
	}

	return url.String(), nil
}
