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
		Long:         `GWA command line interface (CLI) helps manage gateway resources in a declarative fashion.`,
		Version:      ctx.Version,
		PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Debug {
				pkg.Info("Gateway: " + ctx.Gateway)
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
	rootCmd.AddCommand(NewGatewayCmd(ctx, nil))
	rootCmd.AddCommand(GatewayPatternCmd(ctx))
	rootCmd.AddCommand(NewStatusCmd(ctx, nil))
	// Disable these for now since they don't do anything
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gwa-confg.yaml)")
	// rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "only print results, ideal for CI/CD")
	// rootCmd.PersistentFlags().StringVar(&ctx.ApiVersion, "api-version", ctx.ApiVersion, "Set the global API version")
	rootCmd.PersistentFlags().BoolVarP(&ctx.Debug, "debug", "D", false, "Print debug information to stdout when the command has exited")
	rootCmd.PersistentFlags().StringVar(&ctx.ApiHost, "host", ctx.ApiHost, "Set the default host to use for the API")
	rootCmd.PersistentFlags().StringVar(&ctx.Scheme, "scheme", "", "Use to override default https")
	rootCmd.PersistentFlags().StringVar(&ctx.Gateway, "gateway", "", "Assign the Gateway (ID) you would like to use")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return rootCmd
}

func Execute(ctx *pkg.AppContext) *cobra.Command {
	rootCmd := NewRootCommand(ctx)
	cobra.OnInitialize(initConfig, func() {
		applyConfigToContext(ctx, rootCmd)
	})
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	return rootCmd
}

// applyConfigToContext fills ctx from the config file. Cobra has already
// parsed the flags into ctx by the time this runs, so a global flag the user
// passed (--gateway, --scheme) wins over the config file value.
func applyConfigToContext(ctx *pkg.AppContext, rootCmd *cobra.Command) {
	flags := rootCmd.PersistentFlags()

	ctx.ApiKey = viper.GetString("api_key")
	if !flags.Changed("gateway") {
		ctx.Gateway = viper.GetString("gateway")
	}
	if !flags.Changed("scheme") {
		ctx.Scheme = viper.GetString("scheme")
	}

	f := flags.Lookup("host")
	hostValue, _ := flags.GetString("host")

	if viper.GetString("host") != "" && f.DefValue == hostValue {
		ctx.ApiHost = viper.GetString("host")
	}
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
		pkg.Warning(fmt.Sprintf("No config file exists, creating new file at %s/.gwa-config.yaml", home))

		viper.SafeWriteConfig()
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
