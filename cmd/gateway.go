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

func NewGatewayCmd(ctx *pkg.AppContext) *cobra.Command {
	gatewayCmd := &cobra.Command{
		Use:   "gateway",
		Short: "Manage your gateways",
		Long:  `Gateways are used to organize your services.`,
	}
	gatewayCmd.AddCommand(GatewayListCmd(ctx))
	gatewayCmd.AddCommand(GatewayCreateCmd(ctx))
	gatewayCmd.AddCommand(GatewayDestroyCmd(ctx))
	gatewayCmd.AddCommand(GatewayCurrentCmd(ctx))
	return gatewayCmd
}

type GatewayFormData struct {
	GatewayId   string `json:"gatewayId,omitempty"        url:"gatewayId,omitempty"`
	DisplayName string `json:"displayName,omitempty" 	  url:"displayName,omitempty"`
}

func (n *GatewayFormData) IsEmpty() bool {
	return n.DisplayName == "" && n.GatewayId == ""
}

func GatewayListCmd(ctx *pkg.AppContext) *cobra.Command {
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List all your managed gateways",
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, _ []string) error {
			path := fmt.Sprintf("/ds/api/%s/gateways", ctx.ApiVersion)
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
		data := &GatewayFormData{}
		data.DisplayName = m.Prompts[displayName].TextInput.Value()

		gw, err := createGateway(m.Ctx, data)
		if err != nil {
			return pkg.PromptOutputErr{Err: err}
		}
		fmt.Println(gw)

		err = setCurrentGateway(gw)
		if err != nil {
			return pkg.PromptOutputErr{Err: err}
		}

		return pkg.PromptCompleteEvent(gw)
	}
}

const (
	displayName = iota
)

func initialModel(ctx *pkg.AppContext) pkg.GenerateModel {
	var prompts = make([]pkg.PromptField, 1)

	prompts[displayName] = pkg.NewTextInput("Display name", "A short, human-readable name", false)
	prompts[displayName].TextInput.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot

	m := pkg.GenerateModel{
		Action: runCreateRequest,
		Ctx:    ctx,
		Header: heredoc.Doc(`
Create Gateway

Hit enter to accept the default display name (<IDIR>'s gateway) or provide a name below.

Display names may consist of:
- Letters, numbers, spaces or the special characters -()_
- No more than 30 characters

`),
		Prompts: prompts,
		Spinner: s,
	}

	return m
}

func GatewayCreateCmd(ctx *pkg.AppContext) *cobra.Command {
	generate := false
	var gatewayFormData GatewayFormData
	createCommand := &cobra.Command{
		Use:   "create",
		Short: "Create a new gateway",
		Example: heredoc.Doc(`
    $ gwa gateway create --generate
    $ gwa gateway create --gateway-id my-gateway --display-name="This is my gateway"
    `),
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, _ []string) error {
			if gatewayFormData.IsEmpty() && generate == false {
				model := initialModel(ctx)
				if _, err := tea.NewProgram(model).Run(); err != nil {
					return err
				}
				return nil
			}

			gateway, err := createGateway(ctx, &gatewayFormData)
			if err != nil {
				return err
			}

			pkg.Info("Setting gateway to " + gateway)

			err = setCurrentGateway(gateway)
			if err != nil {
				return err
			}

			return nil
		}),
	}
	createCommand.Flags().
		BoolVarP(&generate, "generate", "g", false, "generates a unique gateway with the default display name")
	createCommand.Flags().
		StringVarP(&gatewayFormData.GatewayId, "gateway-id", "i", "", "optionally specify the gateway ID")
	createCommand.Flags().
		StringVarP(&gatewayFormData.DisplayName, "display-name", "d", "", "optionally set the gateway display name")

	return createCommand
}

type GatewayResult struct {
	GatewayId   string `json:"gatewayId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func createGateway(ctx *pkg.AppContext, data *GatewayFormData) (string, error) {
	path := fmt.Sprintf("/ds/api/%s/gateways", ctx.ApiVersion)
	URL, err := ctx.CreateUrl(path, nil)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	r, err := pkg.NewApiPost[GatewayResult](ctx, URL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	response, err := r.Do()
	if err != nil {
		return "", err
	}

	if response.Data.DisplayName != "" {
		fmt.Printf("Gateway created. Gateway ID: %s, display name: %s\n", response.Data.GatewayId, response.Data.DisplayName)
	} else {
		fmt.Printf("Gateway created. Gateway ID: %s\n", response.Data.GatewayId)
	}

	return response.Data.GatewayId, nil
}

func GatewayCurrentCmd(ctx *pkg.AppContext) *cobra.Command {
	currentCmd := &cobra.Command{
		Use:   "current",
		Short: "Display the current gateway",
		RunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Gateway == "" {
				fmt.Println(heredoc.Doc(`
You can create a gateway by running:
    $ gwa gateway create
`))
				return fmt.Errorf("no gateway has been defined")
			}

			fmt.Println(ctx.Gateway)
			return nil
		},
	}
	return currentCmd
}

func setCurrentGateway(gw string) error {
	viper.Set("gateway", gw)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

type GatewayDestroyOptions struct {
	Force bool `url:"force"`
}

func GatewayDestroyCmd(ctx *pkg.AppContext) *cobra.Command {
	var destroyOptions GatewayDestroyOptions
	destroyCommand := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the current gateway",
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

			loader := pkg.NewSpinner()
			loader.Start()
			err := destroyGateway(ctx, &destroyOptions)
			loader.Stop()
			if err != nil {
				return err
			}

			err = setCurrentGateway("")
			if err != nil {
				return err
			}

			fmt.Println("Gateway destroyed:", ctx.Gateway)
			return nil
		}),
	}

	destroyCommand.Flags().BoolVar(&destroyOptions.Force, "force", false, "force deletion")

	return destroyCommand
}

func destroyGateway(ctx *pkg.AppContext, destroyOptions *GatewayDestroyOptions) error {
	pathname := fmt.Sprintf("/ds/api/%s/gateways/%s", ctx.ApiVersion, ctx.Gateway)
	URL, err := ctx.CreateUrl(pathname, destroyOptions)
	if err != nil {
		return err
	}
	r, err := pkg.NewApiDelete[GatewayResult](ctx, URL)
	if err != nil {
		return err
	}

	_, err = r.Do()
	if err != nil {
		return err
	}

	return nil
}

func validateGateway(input string) error {
	pattern := `^[a-zA-Z0-9\-]{3,15}$`
	r := regexp.MustCompile(pattern)

	if !r.MatchString(input) {
		err := fmt.Errorf("%s is an invalid gateway", pkg.BoldStyle.Copy().Underline(true).Render(input))
		return pkg.PromptValidationErr{Err: err}
	}

	return nil
}
