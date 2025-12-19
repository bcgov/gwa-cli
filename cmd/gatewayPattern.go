package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type GatewayPatternOptions struct {
	input string
}

func GatewayPatternCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &GatewayPatternOptions{}
	var gatewayPatternCmd = &cobra.Command{
		Use:     "gateway-pattern input",
		Aliases: []string{"p"},
		Short:   "Generate gateway configuration based on pattern",
		Long: heredoc.Doc(`
    `),
		Example: heredoc.Doc(`
    $ gwa gateway-pattern path/to/config1.yaml
    `),
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, args []string) error {
			if ctx.Gateway == "" {
				fmt.Println(heredoc.Doc(`
          A gateway must be set via the config command

          Example:
            $ gwa config set gateway YOUR_GATEWAY_NAME
        `))
				return fmt.Errorf("no gateway has been set")
			}

			if len(args) == 0 {
				return fmt.Errorf("a pattern input file is required")
			}

			opts.input = args[0]
			config, err := PreparePatternFile(ctx, opts)
			if err != nil {
				return err
			}

			result, err := GatewayPattern(ctx, opts, config)
			if err != nil {
				return err
			}

			yamlContent, err := yaml.Marshal(result.Documents[0])
			if err != nil {
				return err
			}

			fmt.Printf(`%s`, yamlContent)

			return nil
		}),
	}
	return gatewayPatternCmd
}

type GatewayPatternResponse struct {
	Documents []interface{} `json:"documents"`
}

func PreparePatternFile(ctx *pkg.AppContext, opts *GatewayPatternOptions) (io.Reader, error) {
	var validFiles = []string{}

	// validate all the inputs are YAML, if directory loop through
	var input = opts.input
	if input == "-" {
		// read from stdin
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		jsonContent, err := yamlToJson(content)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(jsonContent), nil
	}

	filePath := filepath.Join(ctx.Cwd, input)
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("must be a file")
	} else {
		if isYamlFile(input) {
			validFiles = append(validFiles, filePath)
		}
	}

	if len(validFiles) == 0 {
		return nil, fmt.Errorf("this directory contains no yaml config files")
	}

	// read yaml and convert to json
	content, err := os.ReadFile(validFiles[0])
	if err != nil {
		return nil, err
	}

	jsonContent, err := yamlToJson(content)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonContent), nil
}

func yamlToJson(content []byte) ([]byte, error) {
	var data interface{}

	var err = yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	jsonContent, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}

func GatewayPattern(ctx *pkg.AppContext, opts *GatewayPatternOptions, configFile io.Reader) (GatewayPatternResponse, error) {
	var result GatewayPatternResponse

	body := &bytes.Buffer{}

	var _, err = io.Copy(body, configFile)
	if err != nil {
		return result, err
	}

	path := fmt.Sprintf("/ds/api/%s/gateways/%s/pattern", ctx.ApiVersion, ctx.Gateway)
	URL, _ := ctx.CreateUrl(path, nil)
	r, err := pkg.NewApiPut[GatewayPatternResponse](ctx, URL, body)
	if err != nil {
		return result, err
	}
	r.Request.Header.Set("Content-Type", "application/json")
	contentLength := int64(body.Len())
	r.Request.ContentLength = contentLength

	response, err := r.Do()
	if err != nil {
		return result, err
	}

	result = response.Data

	return result, nil
}
