package cmd

import (
	"embed"
	"fmt"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"os"
)

//go:embed templates/*.go.tmpl
var templates embed.FS

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type GenerateConfigOptions struct {
	template  string
	service   string
	upstream  string
	out       string
}

func NewGenerateConfigCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &GenerateConfigOptions{}
	var generateConfigCmd = &cobra.Command{
		Use:       "generate-config",
		Short:     "Generate gateway configuration based on pre-defined templates",
		Args:      cobra.OnlyValidArgs,
		Example: `
$ gwa generate-config --template kong-httpbin --service my-service --upstream https://httpbin.org
$ gwa generate-config --template client-credentials-shared-idp --service my-service --upstream https://httpbin.org
    `,
		RunE: func(_ *cobra.Command, args []string) error {
			data := struct {
				Namespace   string
				Service     string
				UpstreamUrl    string
			}{
				Namespace:   ctx.Namespace,
				Service:     opts.service,
				UpstreamUrl: opts.upstream,
			}
		
			tmpl := pkg.NewTemplate()

			tplContent, err := templates.ReadFile("templates/" + opts.template + ".go.tmpl")
			check(err)

			tmpl, err = tmpl.Parse(string(tplContent))
			check(err)

			file, err := os.Create(opts.out)
			check(err)
			defer file.Close()

			// Execute the template with the data.
			err = tmpl.Execute(file, data)
			check(err)

			fmt.Println("File " + opts.out + " created")

			return nil
		},
	}

	generateConfigCmd.Flags().StringVarP(&opts.template, "template", "t", "", "Name of a pre-defined template (kong-httpbin, client-credentials-shared-idp)")
	generateConfigCmd.Flags().StringVarP(&opts.service, "service", "s", "", "A service subdomain for your vanity url: https://<service>.api.gov.bc.ca")
	generateConfigCmd.Flags().StringVarP(&opts.upstream, "upstream", "u", "", "The upstream implementation of the API")
	generateConfigCmd.Flags().StringVarP(&opts.out, "out", "o", "gw-config.yml", "The file to output the generate config to")
	generateConfigCmd.MarkFlagRequired("template")
	generateConfigCmd.MarkFlagRequired("service")
	generateConfigCmd.MarkFlagRequired("upstream")

	return generateConfigCmd
}

