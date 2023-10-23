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
	rootCmd.PersistentFlags().StringVar(&ctx.ApiVersion, "api-version", ctx.ApiVersion, "Override the current API version")
	rootCmd.PersistentFlags().StringVar(&ctx.ApiHost, "host", ctx.ApiHost, "Set the default host to use for the API")
	rootCmd.PersistentFlags().StringVar(&ctx.Scheme, "scheme", "", "Use to override default https")
	rootCmd.PersistentFlags().StringVar(&ctx.Namespace, "namespace", "", "Assign the namespace you would like to use")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return rootCmd
}

func Execute(ctx *pkg.AppContext) *cobra.Command {
	rootCmd := NewRootCommand(ctx)
	viper.BindPFlag("namespace", rootCmd.Flags().Lookup("namespace"))
	cobra.OnInitialize(initConfig, func() {
		ctx.ApiKey = viper.GetString("api_key")
		ctx.Namespace = viper.GetString("namespace")
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

		viper.SafeWriteConfig()
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
