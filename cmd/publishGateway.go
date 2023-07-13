package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
		Long:    `Publishing content to come`,
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
			err := Publish(ctx, opts)
			if err != nil {
				return err
			}

			fmt.Println("Gateway config published")

			return nil
		},
	}

	publishGatewayCmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Dry run your API changes before committing to them")

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
	body := bytes.NewReader(jsonBody)

	// Request
	pathname := fmt.Sprintf("/namespaces/%s/gateway", ctx.Namespace)
	URL, _ := ctx.CreateUrl(pathname, nil)
	data, err := pkg.ApiPut[any](ctx, URL, body)
	if err != nil {
		return err
	}

	return nil
}
