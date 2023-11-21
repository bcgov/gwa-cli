package cmd

import (
	"embed"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	tea "github.com/charmbracelet/bubbletea"
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

func (o *GenerateConfigOptions) IsEmpty() bool {
	return o.Template == "" && o.Service == "" && o.Upstream == ""
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

func (o *GenerateConfigOptions) ImportFromForm(m pkg.GenerateModel) tea.Cmd {
	return func() tea.Msg {
		o.Service = m.Prompts[service].Value
		o.Template = m.Prompts[template].Value
		o.Upstream = m.Prompts[upstream].Value
		o.Organization = m.Prompts[organization].Value
		o.OrganizationUnit = m.Prompts[orgUnit].Value
		o.Out = m.Prompts[outfile].Value
		return pkg.PromptCompleteEvent("")
	}

}

func NewGenerateConfigCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &GenerateConfigOptions{}
	var generateConfigCmd = &cobra.Command{
		Use:   "generate-config",
		Short: "Generate gateway resources based on pre-defined templates",
		Args:  cobra.OnlyValidArgs,
		Example: heredoc.Doc(`
$ gwa generate-config --template kong-httpbin \
    --service my-service \
	--upstream https://httpbin.org

$ gwa generate-config --template client-credentials-shared-idp \
    --service my-service \
	--upstream https://httpbin.org
    `),
		PreRun: func(cmd *cobra.Command, _ []string) {
			if !opts.IsEmpty() {
				cmd.MarkFlagRequired("template")
				cmd.MarkFlagRequired("service")
				cmd.MarkFlagRequired("upstream")
			}
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Namespace == "" {
				fmt.Println(heredoc.Doc(`
          A namespace must be set via the config command

          Example:
              $ gwa config set namespace YOUR_NAMESPACE_NAME
          `),
				)
				return fmt.Errorf("No namespace has been set")
			}
			opts.Namespace = ctx.Namespace
			if opts.IsEmpty() {
				model := initGenerateModel(ctx, opts)
				if _, err := tea.NewProgram(model).Run(); err != nil {
					return err
				}
			}
			err := opts.Exec()
			if err != nil {
				return err
			}

			err = GenerateConfig(ctx, opts)
			if err != nil {
				return err
			}

			output := fmt.Sprintf("\n%s File %s created", pkg.Checkmark(), opts.Out)
			fmt.Println(output)

			return nil
		},
	}

	generateConfigCmd.Flags().StringVarP(&opts.Template, "template", "t", "", "Name of a pre-defined template (kong-httpbin, client-credentials-shared-idp)")
	generateConfigCmd.Flags().StringVarP(&opts.Service, "service", "s", "", "A unique service subdomain for your vanity url: https://<service>.api.gov.bc.ca")
	generateConfigCmd.Flags().StringVarP(&opts.Upstream, "upstream", "u", "", "The upstream implementation of the API")
	generateConfigCmd.Flags().StringVar(&opts.Organization, "org", "ministry-of-citizens-services", "Set the organization")
	generateConfigCmd.Flags().StringVar(&opts.OrganizationUnit, "org-unit", "databc", "Set the organization unit")
	generateConfigCmd.Flags().StringVarP(&opts.Out, "out", "o", "gw-config.yml", "The file to output the generate config to")

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

// Prompt Code
const (
	service = iota
	template
	upstream
	organization
	orgUnit
	outfile
)

func initGenerateModel(ctx *pkg.AppContext, opts *GenerateConfigOptions) pkg.GenerateModel {
	var prompts = make([]pkg.PromptField, 6)

	prompts[service] = pkg.NewTextInput("Service", "", true)
	prompts[service].TextInput.Focus()

	prompts[template] = pkg.NewList("Template", []string{
		"kong-httpbin",
		"client-credentials-shared-idp",
	})

	prompts[upstream] = pkg.NewTextInput("Upstream", "", true)
	prompts[upstream].Validator = func(input string) error {
		_, err := url.ParseRequestURI(input)
		return err
	}

	prompts[organization] = pkg.NewTextInput("Organization", "", false)
	prompts[orgUnit] = pkg.NewTextInput("Org Unit", "", false)
	prompts[outfile] = pkg.NewTextInput("Filename", "Must be a YAML file", true)
	prompts[outfile].Validator = func(input string) error {
		if !strings.HasSuffix(input, ".yaml") || !strings.HasSuffix(input, ".yaml") {
			return fmt.Errorf("Must be a yaml filetype")
		}
		return nil
	}

	model := pkg.GenerateModel{
		Action:  opts.ImportFromForm,
		Ctx:     ctx,
		Prompts: prompts,
	}
	return model
}
