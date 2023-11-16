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
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
				return nil
			}
			opts.Namespace = ctx.Namespace
			err := opts.Exec()
			if err != nil {
				return err
			}

			err = GenerateConfig(ctx, opts)
			if err != nil {
				return err
			}

			output := fmt.Sprintf("File %s created", opts.Out)
			fmt.Println(pkg.PrintSuccess(output))

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

type validationErr struct{ error }

func (e validationErr) Error() string {
	return e.error.Error()
}

type outputErr struct{ error }

func (e outputErr) Error() string {
	return e.error.Error()
}

type generateModel struct {
	ctx        *pkg.AppContext
	data       *GenerateConfigOptions
	errorMsg   string
	focusIndex int
	prompts    []pkg.PromptField
	success    bool
}

func (m generateModel) Init() tea.Cmd {
	return nil
}

func (m generateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case validationErr:
		m.errorMsg = msg.Error()
		return m, nil
	case outputErr:
		fmt.Println(fmt.Sprintf("%s Unable to generate output\n%s", pkg.Times(), msg.Error()))
		return m, tea.Quit
	case validValue:
		m.prompts[m.focusIndex].Value = string(msg)
		m.focusIndex++

		if m.focusIndex < len(m.prompts) {
			promptType := m.prompts[m.focusIndex].PromptType
			if promptType == pkg.TextInput {
				return m, m.prompts[m.focusIndex].TextInput.Focus()
			}
		}

		return m, nil
	case success:
		if msg == 1 {
			m.success = true
			fmt.Println(fmt.Sprintf("\n\n%s Success, %s generated", pkg.Checkmark(), pkg.BoldStyle.Render(m.data.Out)))
			return m, tea.Quit
		}
	case tea.KeyMsg:
		m.errorMsg = ""
		key := msg.String()
		if key == "esc" || key == "ctrl+c" {
			return m, tea.Quit
		}

		switch key {
		case "enter":
			totalPrompts := len(m.prompts)
			if totalPrompts == m.focusIndex {
				m.data.Service = m.prompts[service].Value
				m.data.Template = m.prompts[template].Value
				m.data.Upstream = m.prompts[upstream].Value
				m.data.Organization = m.prompts[organization].Value
				m.data.OrganizationUnit = m.prompts[orgUnit].Value
				m.data.Out = m.prompts[outfile].Value

				return m, runGenerateConfig(m)
			}

			if m.focusIndex < totalPrompts {
				return m, validateField(m.prompts[m.focusIndex])
			}
		}
	}

	// Update the currently focused input
	if m.focusIndex < len(m.prompts) {
		current := m.prompts[m.focusIndex]
		switch current.PromptType {
		case pkg.TextInput:
			m.prompts[m.focusIndex].TextInput, cmd = current.TextInput.Update(msg)
		case pkg.ListInput:
			m.prompts[m.focusIndex].List, cmd = current.List.Update(msg)
		}
	}

	return m, cmd
}

func (m generateModel) View() string {
	var b strings.Builder

	for i, p := range m.prompts {
		if i > m.focusIndex {
			continue
		}

		// Render the actul input
		if i == m.focusIndex {
			if m.errorMsg != "" {
				b.WriteString(fmt.Sprintf("%s %s\n", pkg.Times(), m.errorMsg))
			}

			var s string
			switch p.PromptType {
			case pkg.TextInput:
				s = p.TextInput.View()
			case pkg.ListInput:
				b.WriteString(pkg.NewPromptLabel(p.Label))
				s = p.List.View()
			}
			b.WriteString(s)
		}

		// Render the entered value
		if i < m.focusIndex {
			switch p.PromptType {
			case pkg.TextInput:
				b.WriteString(p.TextInput.Prompt)
			case pkg.ListInput:
				b.WriteString(pkg.NewPromptLabel(p.Label))
			}
			b.WriteString(pkg.InputStyle.Render(p.Value))
		}
		b.WriteRune('\n')
	}

	if len(m.prompts) == m.focusIndex {
		buttonText := fmt.Sprintf("%s Submit", pkg.PromptBulletStyle)
		b.WriteString(buttonText)
	}

	return b.String()
}

type success int

func runGenerateConfig(m generateModel) tea.Cmd {
	return func() tea.Msg {
		err := m.data.Exec()
		if err != nil {
			return outputErr{err}
		}
		err = GenerateConfig(m.ctx, m.data)
		if err != nil {
			return outputErr{err}
		}
		return success(1)
	}
}

func initGenerateModel(ctx *pkg.AppContext, opts *GenerateConfigOptions) generateModel {
	var prompts = make([]pkg.PromptField, 6)

	serviceInput := textinput.New()
	serviceInput.Prompt = pkg.NewPromptLabel("Service")
	serviceInput.Focus()
	// serviceInput.TextStyle = pkg.InputStyle
	prompts[service] = pkg.PromptField{
		PromptType: pkg.TextInput,
		IsRequired: true,
		TextInput:  serviceInput,
	}

	templateOptions := []list.Item{
		pkg.ListItem("kong-httpbin"),
		pkg.ListItem("client-credentials-shared-idp"),
	}
	templateInput := pkg.NewList(templateOptions)
	prompts[template] = pkg.PromptField{
		PromptType: pkg.ListInput,
		List:       templateInput,
		Label:      "Template",
	}

	upstreamInput := textinput.New()
	upstreamInput.Prompt = pkg.NewPromptLabel("Upstream")
	prompts[upstream] = pkg.PromptField{
		PromptType: pkg.TextInput,
		IsRequired: true,
		TextInput:  upstreamInput,
		Validator: func(input string) error {
			_, err := url.ParseRequestURI(input)
			return err
		},
	}

	orgInput := textinput.New()
	orgInput.Prompt = pkg.NewPromptLabel("Organization")
	prompts[organization] = pkg.PromptField{
		PromptType: pkg.TextInput,
		TextInput:  orgInput,
	}

	orgUnitInput := textinput.New()
	orgUnitInput.Prompt = pkg.NewPromptLabel("Org Unit")
	prompts[orgUnit] = pkg.PromptField{
		PromptType: pkg.TextInput,
		TextInput:  orgUnitInput,
	}

	outInput := textinput.New()
	outInput.Prompt = pkg.NewPromptLabel("Filename")
	prompts[outfile] = pkg.PromptField{
		PromptType: pkg.TextInput,
		TextInput:  outInput,
		IsRequired: true,
	}

	model := generateModel{
		ctx:     ctx,
		data:    opts,
		prompts: prompts,
	}
	return model
}

type validValue string

func validateField(p pkg.PromptField) tea.Cmd {
	return func() tea.Msg {
		value := p.TextInput.Value()
		switch p.PromptType {
		case pkg.TextInput:
			if p.IsRequired && value == "" {
				return validationErr{fmt.Errorf("Field is required")}
			}
		case pkg.ListInput:
			v, ok := p.List.SelectedItem().(pkg.ListItem)
			if ok {
				return validValue(v)
			}
		}
		if p.Validator != nil {
			err := p.Validator(value)
			if err != nil {
				return validationErr{err}
			}
		}
		return validValue(value)
	}
}
