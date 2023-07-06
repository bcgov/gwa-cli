package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

type publishOptions struct {
	namespace  string
	dryRun     bool
	configFile string
}

func NewPublishGatewayCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &publishOptions{}
	var publishGatewayCmd = &cobra.Command{
		Use:   "publish-gateway [configFile]",
		Short: "Publish your gateway config",
		Long:  `Publishing content to come`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			opts.configFile = args[0]

			err := Publish(ctx, opts)
			if err != nil {
				return err
			}

			fmt.Println("Gateway config published")

			return nil
		},
	}

	publishGatewayCmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Dry run your API changes before committing to them")
	publishGatewayCmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "", "Publish your API to a specific namespace")

	return publishGatewayCmd
}

func Publish(ctx *pkg.AppContext, opts *publishOptions) error {
	filePath := filepath.Join(ctx.Cwd, opts.configFile)
	f, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := map[string]interface{}{
		"configFile": map[string]interface{}{
			"value": string(f),
			"options": map[string]interface{}{
				"filename": opts.configFile,
			},
		},
		"dryRun": opts.dryRun,
	}
	jsonBody, err := json.Marshal(content)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	// Request
	pathname := fmt.Sprintf("/namespaces/%s/gateway", opts.namespace)
	URL, _ := ctx.CreateUrl(pathname, nil)
	request, err := http.NewRequest(http.MethodPut, URL, bodyReader)
	if err != nil {
		return err
	}
	bearer := fmt.Sprintf("bearer %s", ctx.ApiKey)
	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			return err
		}
		// TODO: Return a message from the error JSON if available
		return errors.New(string(body))
	}

	return nil
}
