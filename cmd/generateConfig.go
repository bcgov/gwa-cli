package cmd

import (
	"embed"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

//go:embed templates/*.go.tmpl
var templates embed.FS

type GenerateConfigOptions struct {
	Namespace        string
	Template         string
	Service          string
	Upstream         string
	UpstreamUrl      *url.URL
	UpstreamPort     string
	Organization     string
	OrganizationUnit string
	Out              string
}

func (o *GenerateConfigOptions) ValidateTemplate() error {
	if o.Template == "kong-httpbin" || o.Template == "client-credentials-shared-idp" {
		return nil
	}
	return fmt.Errorf("%s is not a valid template", o.Template)
}

func (o *GenerateConfigOptions) Exec() error {
	err := o.ValidateTemplate()
	if err != nil {
		return err
	}
	err = o.ParseUpstream()
	if err != nil {
		return err
	}
	return nil
}

func (o *GenerateConfigOptions) ParseUpstream() error {
	upstreamUrl, err := url.Parse(o.Upstream)
	if err != nil {
		return err
	}
	o.UpstreamUrl = upstreamUrl
	if upstreamUrl.Port() == "" {
		if upstreamUrl.Scheme == "https" {
			o.UpstreamPort = "443"
		} else {
			o.UpstreamPort = "80"
		}
	} else {
		o.UpstreamPort = upstreamUrl.Port()
	}
	return nil
}

func NewGenerateConfigCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &GenerateConfigOptions{}
	var generateConfigCmd = &cobra.Command{
		Use:   "generate-config",
		Short: "Generate gateway resources based on pre-defined templates",
		Args:  cobra.OnlyValidArgs,
		Example: `
$ gwa generate-config --template kong-httpbin --service my-service --upstream https://httpbin.org
$ gwa generate-config --template client-credentials-shared-idp --service my-service --upstream https://www.boredapi.com/api/activity
    `,
		RunE: func(_ *cobra.Command, _ []string) error {
			opts.Namespace = ctx.Namespace
			err := opts.Exec()
			if err != nil {
				return err
			}

			err = GenerateConfig(ctx, opts)
			if err != nil {
				return err
			}

			fmt.Println("File " + opts.Out + " created")

			return nil
		},
	}

	generateConfigCmd.Flags().StringVarP(&opts.Template, "template", "t", "", "Name of a pre-defined template (kong-httpbin, client-credentials-shared-idp)")
	generateConfigCmd.Flags().StringVarP(&opts.Service, "service", "s", "", "A unique service subdomain for your vanity url: https://<service>.api.gov.bc.ca")
	generateConfigCmd.Flags().StringVarP(&opts.Upstream, "upstream", "u", "", "The upstream implementation of the API")
	generateConfigCmd.Flags().StringVar(&opts.Organization, "org", "ministry-of-citizens-services", "Set the organization")
	generateConfigCmd.Flags().StringVar(&opts.OrganizationUnit, "org-unit", "databc", "Set the organization unit")
	generateConfigCmd.Flags().StringVarP(&opts.Out, "out", "o", "gw-config.yml", "The file to output the generate config to")
	generateConfigCmd.MarkFlagRequired("template")
	generateConfigCmd.MarkFlagRequired("service")
	generateConfigCmd.MarkFlagRequired("upstream")

	return generateConfigCmd
}

func GenerateConfig(ctx *pkg.AppContext, opts *GenerateConfigOptions) error {
	tmpl := pkg.NewTemplate()

	tplContent, err := templates.ReadFile("templates/" + opts.Template + ".go.tmpl")
	if err != nil {
		return err
	}

	tmpl, err = tmpl.Parse(string(tplContent))
	if err != nil {
		return err
	}

	file, err := os.Create(path.Join(ctx.Cwd, opts.Out))
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute the template with the data.
	err = tmpl.Execute(file, opts)
	if err != nil {
		return err
	}
	return nil
}
