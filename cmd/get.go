package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type OutputFlags struct {
	Json bool
	Yaml bool
}

func NewGetCmd(ctx *pkg.AppContext, buf *bytes.Buffer) *cobra.Command {
	var outputOptions = new(OutputFlags)
	var getCmd = &cobra.Command{
		Use:   "get [type] <flags>",
		Short: "Retrieve a table of a namespace's datasets, issuers and products",
		Example: `$ gwa get datasets
$ gwa get datasets --json
$ gwa get datasets --yaml`,
		ValidArgs: []string{"datasets", "issuers", "products"},
		Args:      cobra.OnlyValidArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Must provide an argument of datasets, issuers or products to get command")
			}
			if ctx.Namespace == "" {
				return fmt.Errorf("no namespace selected")
			}
			data, err := CreateAction(ctx, args[0])
			if err != nil {
				return err
			}
			// outputOptions.Print(data, buf)
			if outputOptions.Json {
				json, err := json.Marshal(data)
				if err != nil {
					return err
				}
				fmt.Println(string(json))
				return nil
			}
			if outputOptions.Yaml {
				yaml, err := yaml.Marshal(data)
				if err != nil {
					return err
				}
				fmt.Println(string(yaml))
				return nil
			}
			var tbl table.Table
			switch args[0] {
			case "datasets":
				tbl = table.New("Name", "Title")
				for _, dataset := range data {
					tbl.AddRow(dataset["name"], dataset["title"])
				}
				break
			case "issuers":
				tbl = table.New("Name", "Flow", "Mode", "Owner")
				for _, issuer := range data {
					tbl.AddRow(issuer["name"], issuer["flow"], issuer["mode"], issuer["owner"])
				}
				break
			case "products":
				tbl = table.New("Name", "App ID", "Environments")
				for _, product := range data {
					tbl.AddRow(product["name"], product["appId"], 1)
				}
				break
			}
			if buf != nil {
				tbl.WithWriter(buf)
			}
			tbl.Print()
			return nil
		},
	}

	getCmd.Flags().BoolVar(&outputOptions.Json, "json", false, "Return output as JSON")
	getCmd.Flags().BoolVar(&outputOptions.Yaml, "yaml", false, "Return output as YAML")
	getCmd.MarkFlagsMutuallyExclusive("json", "yaml")
	return getCmd
}

func CreateAction(ctx *pkg.AppContext, operator string) ([]map[string]interface{}, error) {
	var path = fmt.Sprintf("/ds/api/v2/namespaces/%s/%s", ctx.Namespace, operator)
	if operator == "datasets" {
		path = fmt.Sprintf("/ds/api/v2/namespaces/%s/directory", ctx.Namespace)
	}
	url, _ := ctx.CreateUrl(path, nil)

	data, err := GetRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetRequest(ctx *pkg.AppContext, url string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	bearer := fmt.Sprintf("Bearer %s", ctx.ApiKey)
	req.Header.Set("Authorization", bearer)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusOK {
		var result []map[string]interface{}
		json.Unmarshal(body, &result)
		return result, nil
	}
	var errorResponse pkg.ApiErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return nil, fmt.Errorf(string(body))
	}
	return nil, errorResponse.GetError()
}
