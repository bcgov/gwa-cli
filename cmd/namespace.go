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
		Use:   "namespace",
		Short: "Manage your namespaces",
		Long:  `Namespaces are used to organize your services.`,
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

// IsEmpty provides some simple validation on user input
func (n *NamespaceFormData) IsEmpty() bool {
	return n.Description == "" && n.Name == ""
}

func NamespaceListCmd(ctx *pkg.AppContext) *cobra.Command {
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List all your managed namespaces",
		RunE: func(_ *cobra.Command, _ []string) error {
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
				fmt.Println("You have no namespaces")
			}

			for _, n := range response.Data {
				fmt.Println(n)
			}

			return nil
		},
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

// @enum prompt form fields
//
//   - namespace
//   - description
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
      Create Namespace

      Names must be:
      - Alphanumeric (letters, numbers and dashes only, no special characters)
      - Unique to all other namespaces

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
		Short: "Create a new namespace",
		Example: heredoc.Doc(`
    $ gwa namespace create --generate
    $ gwa namespace create --name my-namespace --description="This is my namespace"
    `),
		RunE: func(_ *cobra.Command, _ []string) error {
			if namespaceFormData.IsEmpty() && generate == false {
				model := initialModel(ctx)
				if _, err := tea.NewProgram(model).Run(); err != nil {
					return err
				}
				return nil
			}

			namespace, err := createNamespace(ctx, &namespaceFormData)
			if err != nil {
				return err
			}

			err = setCurrentNamespace(namespace)
			if err != nil {
				return err
			}

			fmt.Println(namespace)
			return nil
		},
	}
	createCommand.Flags().
		BoolVarP(&generate, "generate", "g", false, "generates a random, unique namespace")
	createCommand.Flags().
		StringVarP(&namespaceFormData.Name, "name", "n", "", "optionally define your own namespace")
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
		Short: "Display the current namespace",
		RunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Namespace == "" {
				fmt.Println(heredoc.Doc(`
You can create a namespace by running:
    $ gwa namespace create
`))
				return fmt.Errorf("no namespace has been defined")
			}

			fmt.Println(ctx.Namespace)
			return nil
		},
	}
	return currentCmd
}

func setCurrentNamespace(ns string) error {
	viper.Set("namespace", ns)
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
		Short: "Destroy the current namespace",
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

			fmt.Println("Namespace destroyed:", ctx.Namespace)
			return nil
		},
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

// validateNamespace runs client validation
func validateNamespace(input string) error {
	pattern := `^[a-zA-Z0-9\-]{3,15}$`
	r := regexp.MustCompile(pattern)

	if !r.MatchString(input) {
		err := fmt.Errorf("%s is an invalid namespace", pkg.BoldStyle.Copy().Underline(true).Render(input))
		return pkg.PromptValidationErr{Err: err}
	}

	return nil
}
