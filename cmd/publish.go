package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type PublishOptions struct {
	input string
}

func (o *PublishOptions) ParseInput(ctx *pkg.AppContext) ([]byte, error) {
	filePath := filepath.Join(ctx.Cwd, o.input)
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var configYaml map[string]interface{}
	err = yaml.Unmarshal(file, &configYaml)
	if err != nil {
		return nil, err
	}
	j, err := json.Marshal(configYaml)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func NewPublishCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &PublishOptions{}
	var publishCmd = &cobra.Command{
		Use:       "publish <type>",
		Short:     "Publish to DS API. Available commands are dataset, issuer and product",
		ValidArgs: []string{"dataset", "product", "issuer"},
		Args:      cobra.OnlyValidArgs,
		Example: `
$ gwa publish dataset --input content.yaml
$ gwa publish product --input content.yaml
$ gwa publish issuer --input content.yaml
    `,
		RunE: func(_ *cobra.Command, args []string) error {
			body, err := opts.ParseInput(ctx)
			if err != nil {
				return err
			}

			err = Publish(ctx, body, args[0])
			if err != nil {
				return err
			}

			output := fmt.Sprintf("%s successfully published", args[0])
			fmt.Println(pkg.Checkmark(), strings.ToUpper(output[0:1])+output[1:])

			return nil
		},
	}

	publishCmd.Flags().StringVarP(&opts.input, "input", "i", "", "YAML file to convert to JSON")
	publishCmd.MarkFlagRequired("input")

	return publishCmd
}

type PublishResponse struct {
	Status       int
	Result       string
	Reason       string
	Id           string
	OwnedBy      string
	ChildResults string
}

func Publish(ctx *pkg.AppContext, body []byte, arg string) error {
	route := fmt.Sprintf("/ds/api/v2/namespaces/%s/%ss", ctx.Namespace, arg)
	URL, _ := ctx.CreateUrl(route, nil)
	request, err := pkg.NewApiPut[PublishResponse](ctx, URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = request.Do()
	if err != nil {
		return err
	}

	return nil
}
