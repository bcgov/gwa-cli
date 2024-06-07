package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bcgov/gwa-cli/pkg"
)

func NewNamespaceCmd(ctx *pkg.AppContext) *cobra.Command {
	namespaceCmd := &cobra.Command{
		Use:   "gateway",
		Short: "Manage your gateways",
		Long:  `Gateways are used to organize your services.`,
	}
	namespaceCmd.AddCommand(NamespaceListCmd(ctx))
	namespaceCmd.AddCommand(NamespaceCreateCmd(ctx))
	namespaceCmd.AddCommand(NamespaceDestroyCmd(ctx))
	namespaceCmd.AddCommand(NamespaceCurrentCmd(ctx))
	return namespaceCmd
}

type NamespaceFormData struct {
	Name        string `json:"name,omitempty"        url:"name,omitempty"`
	Description string `json:"displayName,omitempty" url:"description,omitempty"`
}

func (n *NamespaceFormData) IsEmpty() bool {
	return n.Description == "" && n.Name == ""
}

func NamespaceListCmd(ctx *pkg.AppContext) *cobra.Command {
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List all your managed gateways",
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, _ []string) error {
			path := fmt.Sprintf("/ds/api/%s/namespaces", ctx.ApiVersion)
			URL, _ := ctx.CreateUrl(path, nil)
			r, err := pkg.NewApiGet[[]string](ctx, URL)
			if err != nil {
				return err
			}
			loader := pkg.NewSpinner()
			loader.Start()
			response, err := r.Do()
			if err != nil {
				if response.StatusCode == http.StatusUnauthorized {
					fmt.Println()
					fmt.Println(
						heredoc.Doc(`
              Next Steps:
              Run gwa login to obtain another auth token
            `),
					)
				}
				return err
			}
			loader.Stop()

			if len(response.Data) <= 0 {
				fmt.Println("You have no gateways")
			}

			for _, n := range response.Data {
				fmt.Println(n)
			}

			return nil
		}),
	}

	return listCommand
}

// Start Prompt Code
type statusMsg int

func runCreateRequest(m pkg.GenerateModel) tea.Cmd {
	return func() tea.Msg {
		data := &NamespaceFormData{}
		data.Name = m.Prompts[namespace].TextInput.Value()
		data.Description = m.Prompts[description].TextInput.Value()

		ns, err := createNamespace(m.Ctx, data)
		if err != nil {
			return pkg.PromptOutputErr{Err: err}
		}
		return pkg.PromptCompleteEvent(ns)
	}
}

const (
	namespace = iota
	description
)

func initialModel(ctx *pkg.AppContext) pkg.GenerateModel {
	var prompts = make([]pkg.PromptField, 2)

	prompts[namespace] = pkg.NewTextInput("Name", "Must be between 3-15 characters", true)
	prompts[namespace].TextInput.Focus()
	prompts[namespace].Validator = validateNamespace
	prompts[description] = pkg.NewTextInput("Description", "A short, human readable name", false)

	s := spinner.New()
	s.Spinner = spinner.Dot

	m := pkg.GenerateModel{
		Action: runCreateRequest,
		Ctx:    ctx,
		Header: heredoc.Doc(`
      Create Gateway

      Names must be:
      - Alphanumeric (letters, numbers and dashes only, no special characters)
      - Unique to all other gateways

    `),
		Prompts: prompts,
		Spinner: s,
	}

	return m
}

func NamespaceCreateCmd(ctx *pkg.AppContext) *cobra.Command {
	generate := false
	var namespaceFormData NamespaceFormData
	createCommand := &cobra.Command{
		Use:   "create",
		Short: "Create a new gateway",
		Example: heredoc.Doc(`
    $ gwa gateway create --generate
    $ gwa gateway create --name my-gateway --description="This is my gateway"
    `),
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, _ []string) error {
			if namespaceFormData.IsEmpty() && generate == false {
				model := initialModel(ctx)
				if _, err := tea.NewProgram(model).Run(); err != nil {
					return err
				}
				return nil
			}

			gateway, err := createNamespace(ctx, &namespaceFormData)
			if err != nil {
				return err
			}

			pkg.Info("Setting gateway to " + gateway)

			err = setCurrentNamespace(gateway)
			if err != nil {
				return err
			}

			fmt.Println(gateway)
			return nil
		}),
	}
	createCommand.Flags().
		BoolVarP(&generate, "generate", "g", false, "generates a random, unique gateway")
	createCommand.Flags().
		StringVarP(&namespaceFormData.Name, "name", "n", "", "optionally define your own gateway")
	createCommand.Flags().
		StringVarP(&namespaceFormData.Description, "description", "d", "", "optionally add a description")

	return createCommand
}

type NamespaceResult struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func createNamespace(ctx *pkg.AppContext, data *NamespaceFormData) (string, error) {
	path := fmt.Sprintf("/ds/api/%s/namespaces", ctx.ApiVersion)
	URL, err := ctx.CreateUrl(path, nil)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	r, err := pkg.NewApiPost[NamespaceResult](ctx, URL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	response, err := r.Do()
	if err != nil {
		return "", err
	}

	return response.Data.Name, nil
}

func NamespaceCurrentCmd(ctx *pkg.AppContext) *cobra.Command {
	currentCmd := &cobra.Command{
		Use:   "current",
		Short: "Display the current gateway",
		RunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Namespace == "" {
				fmt.Println(heredoc.Doc(`
You can create a gateway by running:
    $ gwa gateway create
`))
				return fmt.Errorf("no gateway has been defined")
			}

			fmt.Println(ctx.Namespace)
			return nil
		},
	}
	return currentCmd
}

func setCurrentNamespace(ns string) error {
	viper.Set("gateway", ns)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

type NamespaceDestroyOptions struct {
	Force bool `url:"force"`
}

func NamespaceDestroyCmd(ctx *pkg.AppContext) *cobra.Command {
	var destroyOptions NamespaceDestroyOptions
	destroyCommand := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the current gateway",
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, _ []string) error {
			if ctx.Namespace == "" {
				fmt.Println(heredoc.Doc(`
          A gateway must be set via the config command

          Example:
              $ gwa config set gateway YOUR_GATEWAY_NAME
          `),
				)
				return fmt.Errorf("No gateway has been set")
			}

			loader := pkg.NewSpinner()
			loader.Start()
			err := destroyNamespace(ctx, &destroyOptions)
			loader.Stop()
			if err != nil {
				return err
			}

			err = setCurrentNamespace("")
			if err != nil {
				return err
			}

			fmt.Println("Gateway destroyed:", ctx.Namespace)
			return nil
		}),
	}

	destroyCommand.Flags().BoolVar(&destroyOptions.Force, "force", false, "force deletion")

	return destroyCommand
}

func destroyNamespace(ctx *pkg.AppContext, destroyOptions *NamespaceDestroyOptions) error {
	pathname := fmt.Sprintf("/ds/api/%s/namespaces/%s", ctx.ApiVersion, ctx.Namespace)
	URL, err := ctx.CreateUrl(pathname, destroyOptions)
	if err != nil {
		return err
	}
	r, err := pkg.NewApiDelete[NamespaceResult](ctx, URL)
	if err != nil {
		return err
	}

	_, err = r.Do()
	if err != nil {
		return err
	}

	return nil
}

func validateNamespace(input string) error {
	pattern := `^[a-zA-Z0-9\-]{3,15}$`
	r := regexp.MustCompile(pattern)

	if !r.MatchString(input) {
		err := fmt.Errorf("%s is an invalid gateway", pkg.BoldStyle.Copy().Underline(true).Render(input))
		return pkg.PromptValidationErr{Err: err}
	}

	return nil
}
