package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type payload map[string]interface{}

type OutputFlags struct {
	Json bool
	Yaml bool
}

func NewGetCmd(ctx *pkg.AppContext, buf *bytes.Buffer) *cobra.Command {
	var outputOptions = new(OutputFlags)
	var validArgs = []string{"datasets", "issuers", "organizations", "products"}
	var getCmd = &cobra.Command{
		Use:   "get [type] <flags>",
		Short: fmt.Sprintf("Retrieve a table of a namespace's %s", pkg.ArgumentsSliceToString(validArgs, "or")),
		Example: `$ gwa get datasets
$ gwa get datasets --json
$ gwa get datasets --yaml`,
		ValidArgs: validArgs,
		Args:      cobra.OnlyValidArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Must provide an argument of %s to get command", pkg.ArgumentsSliceToString(validArgs, "or"))
			}
			if ctx.Namespace == "" {
				return fmt.Errorf("no namespace selected")
			}
			data, err := CreateAction(ctx, args[0])
			if err != nil {
				return err
			}

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
				fallthrough
			case "organizations":
				tbl = table.New("Name", "Title")
				for _, item := range data {
					tbl.AddRow(item["name"], item["title"])
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
					var totalEnvironments int
					if envs, ok := product["environments"].([]interface{}); ok {
						totalEnvironments = len(envs)
					} else {
						totalEnvironments = 0
					}
					tbl.AddRow(product["name"], product["appId"], totalEnvironments)
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

func CreateAction(ctx *pkg.AppContext, operator string) ([]payload, error) {
	var path string
	switch operator {
	case "datasets":
		path = fmt.Sprintf("/ds/api/v2/namespaces/%s/directory", ctx.Namespace)
		break
	case "organizations":
		path = "/ds/api/v2/organizations"
		break
	default:
		path = fmt.Sprintf("/ds/api/v2/namespaces/%s/%s", ctx.Namespace, operator)
	}
	url, _ := ctx.CreateUrl(path, nil)

	err := pkg.RefreshToken(ctx)
	if err != nil {
		return nil, err
	}
	req, err := pkg.NewApiGet[[]payload](ctx, url)
	if err != nil {
		return nil, err
	}
	data, err := req.Do()
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}
