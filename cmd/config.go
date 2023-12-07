package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigCmd(ctx *pkg.AppContext) *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Configuration commands",
	}
	configCmd.AddCommand(NewConfigSetCmd(ctx))
	configCmd.AddCommand(NewConfigGetCmd(ctx))
	return configCmd
}

func NewConfigGetCmd(_ *pkg.AppContext) *cobra.Command {
	args := []string{"api_key", "host", "namespace"}
	argsSentence := pkg.ArgumentsSliceToString(args, "and")

	var configGetCmd = &cobra.Command{
		Use:   "get [key]",
		Short: fmt.Sprintf("Look what value is set for %s", argsSentence),
		Long: heredoc.Docf(`
      This is a convenience getter to print out the currently stored global setting for the following arguments

      - api_key
      - host
      - namespace
    `),
		ValidArgs: args,
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(_ *cobra.Command, args []string) error {
			result := viper.Get(args[0])
			if result == "" {
				return nil
			}
			fmt.Println(result)
			return nil
		},
	}

	return configGetCmd
}

func NewConfigSetCmd(ctx *pkg.AppContext) *cobra.Command {
	var configSetCmd = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Write a specific global setting",
		Long: heredoc.Docf(`
Exposes some specific config values that can be defined by the user.

%s
  namespace:       The default namespace used
  token:           Use only if you have a token you know is authenticated
  host:            The API host you wish to communicate with
  scheme:          http or https

    `, lipgloss.NewStyle().Bold(true).Render("Configurable Settings:")),
		Example: heredoc.Doc(`
$ gwa config set namespace ns-sampler
$ gwa config set --namespace ns-sampler
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			totalArgs := len(args)
			if totalArgs > 1 {
				switch args[0] {
				case "token":
					viper.Set("api_key", args[1])
				case "namespace":
					viper.Set("namespace", args[1])
				case "host":
					viper.Set("host", args[1])
				case "scheme":
					viper.Set("scheme", args[1])
				default:
					return fmt.Errorf("The key <%s> is not allowed to be set", args[0])
				}
			}

			if totalArgs == 1 {
				return fmt.Errorf("No value was set for %s", args)
			}

			if totalArgs == 0 && !cmd.HasFlags() {
				model := initialSetModel(ctx)
				if _, err := tea.NewProgram(model).Run(); err != nil {
					return err
				}
				return nil
			}

			err := viper.WriteConfig()
			if err != nil {
				return err
			}
			fmt.Println(pkg.Checkmark(), pkg.PrintSuccess("Config settings saved"))
			return nil
		},
	}

	configSetCmd.Flags().String("token", "", "set the namespace")
	viper.BindPFlag("api_key", configSetCmd.Flags().Lookup("token"))
	configSetCmd.Flags().String("namespace", "", "set the namespace")
	viper.BindPFlag("namespace", configSetCmd.Flags().Lookup("namespace"))
	configSetCmd.Flags().String("host", "", "set the host")
	viper.BindPFlag("host", configSetCmd.Flags().Lookup("host"))
	configSetCmd.Flags().String("scheme", "", "set the scheme")
	viper.BindPFlag("scheme", configSetCmd.Flags().Lookup("scheme"))

	return configSetCmd
}

const (
	key = iota
	value
)

func initialSetModel(ctx *pkg.AppContext) pkg.GenerateModel {
	var prompts = make([]pkg.PromptField, 2)

	prompts[key] = pkg.NewList("Select a config key to set", []string{"host", "namespace", "scheme", "token"})
	prompts[value] = pkg.NewTextInput("Value", "", true)

	m := pkg.GenerateModel{
		Action:  setConfigValue,
		Ctx:     ctx,
		Prompts: prompts,
		Spinner: spinner.New(),
	}

	return m
}

func setConfigValue(m pkg.GenerateModel) tea.Cmd {
	return func() tea.Msg {
		key := m.Prompts[key].Value
		value := m.Prompts[value].TextInput.Value()

		viper.Set(key, value)
		err := viper.WriteConfig()
		if err != nil {
			return pkg.PromptOutputErr{Err: err}
		}
		return pkg.PromptCompleteEvent("")
	}
}
