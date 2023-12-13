package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

type PublishGatewayOptions struct {
	dryRun    bool
	qualifier string
	inputs    []string
}

func NewPublishGatewayCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &PublishGatewayOptions{}
	var publishGatewayCmd = &cobra.Command{
		Use:     "publish-gateway [inputs...]",
		Aliases: []string{"pg"},
		Short:   "Publish your Kong gateway configuration",
		Long: heredoc.Doc(`
    Once you have a gateway configuration file ready to publish, you can run the following command to reflect your changes in the gateway:

      $ gwa pg sample.yaml

    If you want to see the expected changes but not actually apply them, you can run:

      $ gwa pg --dry-run sample.yaml

    inputs accepts a wide variety of formats, for example:

      1. Empty, which means find all the possible YAML files in the current directory and publish them
      2. A space-separated list of specific YAML files in the current directory, or
      3. A directory relative to the current directory
    `),
		Example: heredoc.Doc(`
    $ gwa publish-gateway
    $ gwa publish-gateway path/to/config1.yaml other-path/to/config2.yaml
    $ gwa publish-gateway path/to/directory/containing-configs/
    $ gwa publish-gateway path/to/config.yaml --dry-run
    $ gwa publish-gateway path/to/config.yaml --qualifier dev
    `),
		RunE: func(_ *cobra.Command, args []string) error {
			if ctx.Namespace == "" {
				fmt.Println(heredoc.Doc(`
          A namespace must be set via the config command

          Example:
            $ gwa config set namespace YOUR_NAMESPACE_NAME
        `))
				return fmt.Errorf("No namespace has been set\n")
			}

			opts.inputs = args
			if len(args) == 0 {
				opts.inputs = []string{""}
			}
			config, err := PrepareConfigFile(ctx, opts)
			if err != nil {
				return err
			}

			result, err := PublishToGateway(ctx, opts, config)
			if err != nil {
				return err
			}

			fmt.Println(pkg.Checkmark(), "Gateway config published")
			fmt.Printf(`
Details:
   %s

%s
`, result.Message, result.Results)

			return nil
		},
	}

	publishGatewayCmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Dry run your API changes before committing to them")
	publishGatewayCmd.Flags().StringVar(&opts.qualifier, "qualifier", "", "Sets a tag qualifier, which specifies that the gateway configuration is a partial set of configuration")

	return publishGatewayCmd
}

type PublishGatewayResponse struct {
	Message string `json:"message"`
	Results string `json:"results"`
	Error   string `json:"error"`
}

// isYamlFile will use the Ext module to determine the extension
func isYamlFile(filename string) bool {
	ext := filepath.Ext(filename)

	if ext == ".yaml" || ext == ".yml" {
		return true
	}
	return false
}

func PrepareConfigFile(ctx *pkg.AppContext, opts *PublishGatewayOptions) (io.Reader, error) {
	var resultBuffer = []byte("")
	var validFiles = []string{}

	// validate all the inputs are YAML, if directory loop through
	for _, input := range opts.inputs {
		filePath := filepath.Join(ctx.Cwd, input)
		info, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			files, err := os.ReadDir(filePath)
			if err != nil {
				return nil, err
			}

			// Filter all the files in the dir
			for _, f := range files {
				filename := f.Name()
				if isYamlFile(filename) {
					fileInFolder := filepath.Join(ctx.Cwd, input, filename)
					validFiles = append(validFiles, fileInFolder)
				}
			}
		} else {
			if isYamlFile(input) {
				validFiles = append(validFiles, filePath)
			}
		}
	}

	if len(validFiles) == 0 {
		return nil, fmt.Errorf("This directory contains no yaml config files\n")
	}

	for i, file := range validFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		if i > 0 {
			resultBuffer = append(resultBuffer, []byte("\n---\n")...)
		}

		resultBuffer = append(resultBuffer, content...)
	}

	return bytes.NewReader(resultBuffer), nil
}

func PublishToGateway(ctx *pkg.AppContext, opts *PublishGatewayOptions, configFile io.Reader) (PublishGatewayResponse, error) {
	var result PublishGatewayResponse

	body := &bytes.Buffer{}
	fw := multipart.NewWriter(body)

	dryRunField, err := fw.CreateFormField("dryRun")
	if err != nil {
		return result, err
	}

	dryRunValue := strconv.FormatBool(opts.dryRun)
	dryRunField.Write([]byte(dryRunValue))

	if opts.qualifier != "" {
		qualifierField, err := fw.CreateFormField("qualifier")
		if err != nil {
			return result, err
		}
		qualifierField.Write([]byte(opts.qualifier))
	}

	fileField, err := fw.CreateFormFile("configFile", "config.yaml")
	if err != nil {
		return result, err
	}

	_, err = io.Copy(fileField, configFile)
	if err != nil {
		return result, err
	}

	err = fw.Close()
	if err != nil {
		return result, err
	}

	path := fmt.Sprintf("/gw/api/%s/namespaces/%s/gateway", ctx.ApiVersion, ctx.Namespace)
	URL, _ := ctx.CreateUrl(path, nil)
	r, err := pkg.NewApiPut[PublishGatewayResponse](ctx, URL, body)
	if err != nil {
		return result, err
	}
	r.Request.Header.Set("Content-Type", fw.FormDataContentType())
	contentLength := int64(body.Len())
	r.Request.ContentLength = contentLength

	response, err := r.Do()
	if err != nil {
		return result, err
	}

	result = response.Data

	return result, nil
}
