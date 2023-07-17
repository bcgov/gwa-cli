package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

type publishOptions struct {
	dryRun     bool
	configFile string
}

func NewPublishGatewayCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &publishOptions{}
	var publishGatewayCmd = &cobra.Command{
		Use:     "publish-gateway [configFile]",
		Aliases: []string{"pg"},
		Short:   "Publish your gateway config",
		Args:    cobra.MinimumNArgs(1),
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

			fmt.Println("Gateway config published")
			fmt.Printf(`
Details:
  %s

%s
`, result.Message, result.Results)

			return nil
		},
	}

	publishGatewayCmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Dry run your API changes before committing to them")

	return publishGatewayCmd
}

type PublishGatewayResponse struct {
	Message string `json:"message"`
	Results string `json:"results"`
	Error   string `json:"error"`
}

func PublishGateway(ctx *pkg.AppContext, opts *publishOptions) (PublishGatewayResponse, error) {
	var result PublishGatewayResponse
	// Open the file
	filePath := filepath.Join(ctx.Cwd, opts.configFile)
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	fw := multipart.NewWriter(body)

	dryRunField, err := fw.CreateFormField("dryRun")
	if err != nil {
		return result, err
	}

	dryRunValue := strconv.FormatBool(opts.dryRun)
	dryRunField.Write([]byte(dryRunValue))

	fileField, err := fw.CreateFormFile("configFile", file.Name())
	if err != nil {
		return result, err
	}

	_, err = io.Copy(fileField, file)
	if err != nil {
		return result, err
	}

	err = fw.Close()
	if err != nil {
		return result, err
	}

	pathname := fmt.Sprintf("/gw/api/namespaces/%s/gateway", ctx.Namespace)
	URL, _ := ctx.CreateUrl(pathname, nil)
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
