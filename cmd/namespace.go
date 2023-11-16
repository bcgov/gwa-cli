package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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

var errorHeaderText = fmt.Sprintf("%s Namespace create failed:", pkg.Times())

type namespaceCreateModel struct {
	ctx          *pkg.AppContext
	data         *NamespaceFormData
	err          error
	focusIndex   int
	inputs       []textinput.Model
	isRequesting bool
	invalid      error
	spinner      spinner.Model
	status       statusMsg
}

func (m namespaceCreateModel) startSpinner() {
	m.spinner = spinner.New()
}

type nsValueErrMsg struct {
	err error
}

func (e nsValueErrMsg) Error() string {
	return e.err.Error()
}

type requestErrMsg struct {
	err error
}

func (e requestErrMsg) Error() string {
	return e.err.Error()
}

func newNamespaceCreateRequest(ctx *pkg.AppContext, data *NamespaceFormData) tea.Cmd {
	return func() tea.Msg {
		_, err := createNamespace(ctx, data)
		if err != nil {
			return requestErrMsg{err}
		}
		return statusMsg(200)
	}
}

func initialModel(ctx *pkg.AppContext) namespaceCreateModel {
	m := namespaceCreateModel{
		ctx:     ctx,
		data:    &NamespaceFormData{},
		inputs:  make([]textinput.Model, 2),
		invalid: nsValueErrMsg{fmt.Errorf("empty")},
		spinner: spinner.New(),
	}

	for i := range m.inputs {
		var t textinput.Model
		t = textinput.New()
		t.PromptStyle = pkg.BoldStyle
		switch i {
		case 0:
			t.Prompt = "Name: "
			t.Placeholder = "Max-length: 15 characters"
			t.CharLimit = 15
			t.Focus()
		case 1:
			t.Prompt = "Description: "
			t.Placeholder = "(optional)"
		}
		m.inputs[i] = t
	}

	return m
}

func (m namespaceCreateModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m namespaceCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Request handlers
	case statusMsg:
		if msg == 200 {
			m.isRequesting = false
			m.status = msg
			return m, tea.Quit
		} else {
			var cmd tea.Cmd
			m.invalid = nil
			return m, cmd
		}
	case requestErrMsg:
		m.isRequesting = false
		m.err = msg
		return m, tea.Quit
	case nsValueErrMsg:
		var cmd tea.Cmd
		m.invalid = msg
		return m, cmd
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	// Keyboard event handlers
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter", "tab", "shift+tab", "up", "down":
			if key == "enter" && m.focusIndex == len(m.inputs) {
				m.data.Name = m.inputs[0].Value()
				m.data.Description = m.inputs[1].Value()
				m.focusIndex = len(m.inputs) + 1
				m.status = 0
				m.isRequesting = true
				m.startSpinner()
				return m, newNamespaceCreateRequest(m.ctx, m.data)
			}

			// Namespace string validation
			if m.invalid != nil {
				return m, validateNamespace(m.inputs[0].Value())
			}

			// Movement
			if key == "up" || key == "shift+up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			// Update state of inputs
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m namespaceCreateModel) View() string {
	var b strings.Builder

	b.WriteString(heredoc.Doc(`
    Create Namespace

    Names must be:
    - Alphanumeric (letters, numbers and dashes only, no special characters)
    - Unique to all other namespaces

  `))

	for i := range m.inputs {
		if i == m.focusIndex {
			if m.invalid != nil {
				b.WriteString(pkg.PromptErrorStyle.Render("! "))
			} else {
				b.WriteString(pkg.PromptSymbolStyle.Render("> "))
			}
		} else {
			b.WriteString(pkg.PromptSymbolStyle.Render("? "))
		}

		b.WriteString(m.inputs[i].View())

		if i < len(m.inputs) {
			b.WriteRune('\n')
		}
	}

	// Request results
	if m.isRequesting {
		s := fmt.Sprintf("\n%s Creating namespace...", m.spinner.View())
		b.WriteString(s)
		if m.err != nil {
			b.WriteRune('\n')
			b.WriteString(pkg.ErrorStyle.Render(errorHeaderText))
			b.WriteRune('\n')
			b.WriteString(m.err.Error())
			b.WriteRune('\n')
		}
	}

	if m.status != 0 {
		b.WriteString(fmt.Sprintf("\n%s %s created", pkg.Checkmark(), m.data.Name))
		// NOTE: We could add some next steps, nicely formatted here
	}

	return b.String()
}

func (m namespaceCreateModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
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

func validateNamespace(input string) tea.Cmd {
	return func() tea.Msg {
		pattern := `^[a-zA-Z0-9\-]{3,15}$`
		r := regexp.MustCompile(pattern)

		if !r.MatchString(input) {
			err := fmt.Errorf("%s is an invalid namespace", input)
			return nsValueErrMsg{err}
		}

		return statusMsg(1)
	}
}
