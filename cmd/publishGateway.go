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
	dryRun     bool
	qualifier  string
	configFile string
}

func NewPublishGatewayCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &PublishGatewayOptions{}
	var publishGatewayCmd = &cobra.Command{
		Use:     "publish-gateway [configFile]",
		Aliases: []string{"pg"},
		Short:   "Publish your gateway config",
		Example: heredoc.Doc(`
    $ gwa publish-gateway path/to/config.yaml
    $ gwa publish-gateway path/to/config.yaml --dry-run
    `),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if ctx.Namespace == "" {
				cmd.SetUsageTemplate(`
A namespace must be set via the config command

Example:
    $ gwa config set namespace YOUR_NAMESPACE_NAME
`)
				return fmt.Errorf("No namespace has been set")
			}
			opts.configFile = args[0]
			result, err := PublishGateway(ctx, opts)
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

func PublishGateway(ctx *pkg.AppContext, opts *PublishGatewayOptions) (PublishGatewayResponse, error) {
	var result PublishGatewayResponse
	// Open the file
	filePath := filepath.Join(ctx.Cwd, opts.configFile)
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	return PublishToGateway(ctx, opts, file)
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
