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
	Gateway          string
	Template         string
	Service          string
	Upstream         string
	UpstreamUrl      *url.URL
	UpstreamPort     string
	Organization     string
	OrganizationUnit string
	Out              string
}

type Response struct {
	Available  bool       `json:"available"`
	Suggestion Suggestion `json:"suggestion"`
}

type Suggestion struct {
	ServiceName string   `json:"serviceName"`
	Names       []string `json:"names"`
	Hosts       []string `json:"hosts"`
}

func (o *GenerateConfigOptions) IsEmpty() bool {
	return o.Template == "" && o.Service == "" && o.Upstream == ""
}

func (o *GenerateConfigOptions) ValidateTemplate() error {
	if o.Template == "kong-httpbin" || o.Template == "client-credentials-shared-idp" || o.Template == "quick-start" {
		return nil
	}
	return fmt.Errorf("%s is not a valid template", o.Template)
}

func (o *GenerateConfigOptions) ValidateService(ctx *pkg.AppContext, service string) error {
	path := fmt.Sprintf("/ds/api/%s/routes/availability?gatewayId=%s&serviceName=%s", ctx.ApiVersion, ctx.Gateway, service)
	URL, _ := ctx.CreateUrl(path, nil)
	decodedURL, err := url.QueryUnescape(URL)
	if err != nil {
		return err
	}
	request, err := pkg.NewApiGet[Response](ctx, decodedURL)
	if err != nil {
		return err
	}

	loader := pkg.NewSpinner()
	loader.Suffix = " Checking service availability"
	loader.Start()
	response, err := request.Do()
	if err != nil {
		return err
	}
	loader.Stop()

	if !response.Data.Available {
		return fmt.Errorf("Checking service availability: Service %s is already in use. Suggestion: %s", service, response.Data.Suggestion.ServiceName)
	}

	return nil
}

func (o *GenerateConfigOptions) Exec(ctx *pkg.AppContext) error {
	err := o.ValidateTemplate()
	if err != nil {
		return err
	}
	err = o.ValidateService(ctx, o.Service)
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
$ gwa generate-config --template quick-start \
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
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, _ []string) error {
			if ctx.Gateway == "" {
				fmt.Println(heredoc.Doc(`
          A gateway must be set via the config command

          Example:
              $ gwa config set gateway YOUR_GATEWAY_NAME
          `),
				)
				return fmt.Errorf("No gateway has been set")
			}

			opts.Gateway = ctx.Gateway
			pkg.Info(fmt.Sprintf("Options received %v", opts))

			if opts.IsEmpty() {
				model := initGenerateModel(ctx, opts)
				if _, err := tea.NewProgram(model).Run(); err != nil {
					return err
				}
			}
			err := opts.Exec(ctx)
			if err != nil {
				return err
			}
			pkg.Info("Options executed")

			err = GenerateConfig(ctx, opts)
			if err != nil {
				return err
			}

			output := fmt.Sprintf("\n%s File %s created", pkg.Checkmark(), opts.Out)
			fmt.Println(output)

			return nil
		}),
	}

	generateConfigCmd.Flags().StringVarP(&opts.Template, "template", "t", "", "Name of a pre-defined template (quick-start, client-credentials-shared-idp, kong-httpbin)")
	generateConfigCmd.Flags().StringVarP(&opts.Service, "service", "s", "", "A unique service subdomain for your vanity url: https://<service>.api.gov.bc.ca")
	generateConfigCmd.Flags().StringVarP(&opts.Upstream, "upstream", "u", "", "The upstream implementation of the API")
	generateConfigCmd.Flags().StringVar(&opts.Organization, "org", ctx.DefaultOrg, "Set the organization")
	generateConfigCmd.Flags().StringVar(&opts.OrganizationUnit, "org-unit", ctx.DefaultOrgUnit, "Set the organization unit")
	generateConfigCmd.Flags().StringVarP(&opts.Out, "out", "o", "gw-config.yaml", "The file to output the generate config to")

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
	pkg.Info(fmt.Sprintf("%s template parsed", opts.Template))

	file, err := os.Create(path.Join(ctx.Cwd, opts.Out))
	if err != nil {
		return err
	}
	defer file.Close()
	pkg.Info("File created")

	// Execute the template with the data.
	err = tmpl.Execute(file, opts)
	if err != nil {
		return err
	}
	pkg.Info("Template successfully parsed")
	return nil
}

// Prompt Code
const (
	service = iota
	template
	upstream
	outfile
)

func initGenerateModel(ctx *pkg.AppContext, opts *GenerateConfigOptions) pkg.GenerateModel {
	var prompts = make([]pkg.PromptField, 4)

	prompts[service] = pkg.NewTextInput("Service", "", true)
	prompts[service].TextInput.Focus()
	prompts[service].Validator = func(input string) error {
		err := opts.ValidateService(ctx, input)
		if err != nil {
			return err
		}
		return nil
	}

	prompts[template] = pkg.NewList("Template", []string{
		"client-credentials-shared-idp",
		"kong-httpbin",
		"quick-start",
	})

	prompts[upstream] = pkg.NewTextInput("Upstream (URL)", "", true)
	prompts[upstream].Validator = func(input string) error {
		_, err := url.ParseRequestURI(input)
		return err
	}

	opts.Organization = ctx.DefaultOrg
	opts.OrganizationUnit = ctx.DefaultOrgUnit

	prompts[outfile] = pkg.NewTextInput("Filename", "Must be a YAML file", true)
	prompts[outfile].TextInput.SetValue("gw-config.yaml")
	prompts[outfile].Validator = func(input string) error {
		if strings.HasSuffix(input, ".yml") || strings.HasSuffix(input, ".yaml") {
			return nil
		}
		return fmt.Errorf("Filename %s is invalid. Only YAML files are accepted.", pkg.BoldStyle.Underline(true).Render(input))
	}

	model := pkg.GenerateModel{
		Action:  opts.ImportFromForm,
		Ctx:     ctx,
		Prompts: prompts,
	}
	return model
}
