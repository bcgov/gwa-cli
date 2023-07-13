package cmd

import (
	"fmt"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigCmd(ctx *pkg.AppContext) *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Configuration commands",
	}
	configCmd.AddCommand(NewConfigSetCmd(ctx))
	return configCmd
}

func NewConfigSetCmd(_ *pkg.AppContext) *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Write a specific global setting",
		Long: `
Exposes some specific config values that can be defined by the user.

Configurable Settings:
  namespace:       The default namespace used
  token:           Use only if you have a token you know is authenticated
  host:            The API host you wish to communicate with
  scheme:          http or https

`,
		Example: `
$ gwa config set namespace ns-sampler
$ gwa config set --namespace ns-sampler
    `,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 1 {
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

			if len(args) == 1 {
				return fmt.Errorf("No value was set for %s", args)
			}

			err := viper.WriteConfig()
			if err != nil {
				return err
			}
			fmt.Println("Config settings saved")
			return nil
		},
	}

	configCmd.Flags().String("token", "", "set the namespace")
	viper.BindPFlag("api_key", configCmd.Flags().Lookup("token"))
	configCmd.Flags().String("namespace", "", "set the namespace")
	viper.BindPFlag("namespace", configCmd.Flags().Lookup("namespace"))
	configCmd.Flags().String("host", "", "set the host")
	viper.BindPFlag("host", configCmd.Flags().Lookup("host"))
	configCmd.Flags().String("scheme", "", "set the scheme")
	viper.BindPFlag("scheme", configCmd.Flags().Lookup("scheme"))

	return configCmd
}
