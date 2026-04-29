package pkg

import (
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/spf13/cobra"
)

type AppContext struct {
	ApiHost        string
	ApiKey         string
	Auth           AuthDetails
	ApiVersion     string
	ClientId       string
	Cwd            string
	DefaultOrg     string
	DefaultOrgUnit string
	Debug          bool
	Host           string
	Gateway        string
	Scheme         string
	Version        string
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

// Cobra's error handling doesn't bubble to the root's lifecycle hooks
// so this HOF is used to catch errors and print them before exit
// Usage
// ```go
//
//	var cmd = &cobra.Command{
//	  Use: "cmd",
//	  RunE: pkg.WrapError(ctx, func (cmd *cobra.Command, args []string) error {
//	    ...
//	  }),
//	}
//
// ```
func WrapError(ctx *AppContext, handler func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := handler(cmd, args)
		if err != nil {
			if ctx.Debug {
				PrintLog()
			}
			return err
		}
		return nil
	}
}
