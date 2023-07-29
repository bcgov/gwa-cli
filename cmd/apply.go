package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ApplyOptions struct {
	input string
}

func (o *ApplyOptions) ParseInput(ctx *pkg.AppContext) ([][]byte, error) {
	filePath := filepath.Join(ctx.Cwd, o.input)
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	splitDocs, err := pkg.SplitYAML(file)
	if err != nil {
		return nil, err
	}

	return splitDocs, nil
}

func NewApplyCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &ApplyOptions{}
	var applyCmd = &cobra.Command{
		Use:   "apply <type>",
		Short: "Apply gateway resources",
		Args:  cobra.OnlyValidArgs,
		Example: `
$ gwa apply --input gw-config.yaml
    `,
		RunE: func(_ *cobra.Command, args []string) error {
			kindMapper := map[string]string{
				"CredentialIssuer": "issuer",
				"DraftDataset":     "dataset",
				"Product":          "product",
				"Environment":      "environment",
			}

			yamlDocs, err := opts.ParseInput(ctx)
			if err != nil {
				return err
			}

			for _, v := range yamlDocs {
				var configYaml map[string]interface{}
				err = yaml.Unmarshal(v, &configYaml)
				if err != nil {
					return err
				}

				var kind = configYaml["kind"].(string)

				delete(configYaml, "kind")
				body, err := json.Marshal(configYaml)
				if err != nil {
					return err
				}

				if apiPath, ok := kindMapper[kind]; ok {
					data, err := Put(ctx, body, apiPath)
					if err != nil {
						return err
					}
					fmt.Printf("%-20s %-40s %s\n", kind, configYaml["name"], data.Result)
				} else if kind == "GatewayService" {
					fmt.Printf("%-20s %-40s publishing...", kind, configYaml["name"])

					var kongConfig = struct {
						Services []map[string]interface{} `json:"services"`
					}{}

					kongConfig.Services = append([]map[string]interface{}{}, configYaml)

					body, err := json.Marshal(kongConfig)
					if err != nil {
						return err
					}

					_, err = PublishToGateway(ctx, false, bytes.NewReader(body))
					if err != nil {
						return err
					}
					fmt.Printf("\r")
					fmt.Printf("%-20s %-40s published    \n", kind, configYaml["name"])
					// fmt.Printf(`
					// 	Details:
					// 	%s
					// 	%s
					// 	`, result.Message, result.Results)

				} else {
					fmt.Printf("%-20s %-40s skipped\n", kind, configYaml["name"])
				}
			}
			return nil
		},
	}

	applyCmd.Flags().StringVarP(&opts.input, "input", "i", "gw-config.yml", "YAML file containing your configuration")

	return applyCmd
}

type PutResponse struct {
	Status       int
	Result       string
	Reason       string
	Id           string
	OwnedBy      string
	ChildResults string
}

func Put(ctx *pkg.AppContext, body []byte, arg string) (*PutResponse, error) {
	route := fmt.Sprintf("/ds/api/v2/namespaces/%s/%ss", ctx.Namespace, arg)
	URL, _ := ctx.CreateUrl(route, nil)
	request, err := pkg.NewApiPut[PutResponse](ctx, URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	res, err := request.Do()
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}
