package cmd

import (
	"fmt"
	"os"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var quiet bool

func NewRootCommand(ctx *pkg.AppContext) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:          "gwa <command> <subcommand> [flags]",
		Short:        "CLI tool supported by the APS team",
		SilenceUsage: true,
		Long:         `GWA CLI helps manage gateway resources in a declarative fashion.`,
		Version:      ctx.Version,
		PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Debug {
				pkg.Info("Namespace: " + ctx.Namespace)
				pkg.Info("API Version: " + ctx.ApiVersion)
				pkg.PrintLog()
			}
			return nil
		},
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			pkg.CheckForVersion(ctx)
		},
	}
	rootCmd.AddCommand(NewConfigCmd(ctx))
	rootCmd.AddCommand(NewInit(ctx))
	rootCmd.AddCommand(NewPublishGatewayCmd(ctx))
	rootCmd.AddCommand(NewPublishCmd(ctx))
	rootCmd.AddCommand(NewGetCmd(ctx, nil))
	rootCmd.AddCommand(NewApplyCmd(ctx))
	rootCmd.AddCommand(NewGenerateConfigCmd(ctx))
	rootCmd.AddCommand(NewLoginCmd(ctx))
	rootCmd.AddCommand(NewNamespaceCmd(ctx))
	rootCmd.AddCommand(NewStatusCmd(ctx, nil))
	// Disable these for now since they don't do anything
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gwa-confg.yaml)")
	// rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "only print results, ideal for CI/CD")
	rootCmd.PersistentFlags().BoolVarP(&ctx.Debug, "debug", "D", false, "Print debug information to stdout when the command has exited")
	rootCmd.PersistentFlags().StringVar(&ctx.ApiVersion, "api-version", ctx.ApiVersion, "Set the global API version")
	rootCmd.PersistentFlags().StringVar(&ctx.ApiHost, "host", ctx.ApiHost, "Set the default host to use for the API")
	rootCmd.PersistentFlags().StringVar(&ctx.Scheme, "scheme", "", "Use to override default https")
	rootCmd.PersistentFlags().StringVar(&ctx.Namespace, "gateway", "", "Assign the gateway you would like to use")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return rootCmd
}

func Execute(ctx *pkg.AppContext) *cobra.Command {
	rootCmd := NewRootCommand(ctx)
	viper.BindPFlag("gateway", rootCmd.Flags().Lookup("gateway"))
	cobra.OnInitialize(initConfig, func() {
		ctx.ApiKey = viper.GetString("api_key")
		ctx.Namespace = viper.GetString("gateway")
		ctx.Scheme = viper.GetString("scheme")

		f := rootCmd.Flags().Lookup("host")
		hostValue, _ := rootCmd.Flags().GetString("host")

		if viper.GetString("host") != "" && f.DefValue == hostValue {
			ctx.ApiHost = viper.GetString("host")
		}
	})
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	return rootCmd
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gwa-config")
		viper.SetDefault("scheme", "https")
		pkg.Warning(fmt.Sprintf("No config file exists, creating new file at %s/.gwa-config.yml", home))

		viper.SafeWriteConfig()
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
