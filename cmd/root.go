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

var rootCmd = &cobra.Command{
	Use:     "gwa <command> <subcommand> [flags]",
	Short:   "CLI tool supported by the APS team",
	Long:    `GWA CLI is a tool for composing, validating and generating Kong Gateway configuration files from OpenAPI (aka Swagger) specs and managing Kong Plugins.`,
	Version: "2.0.0-beta",
}

func Execute(ctx *pkg.AppContext) {
	cobra.OnInitialize(initConfig, func() {
		ctx.ApiKey = viper.GetString("api_key")
	})
	rootCmd.AddCommand(NewInit(ctx))
	rootCmd.AddCommand(NewPublishGatewayCmd(ctx))
	rootCmd.AddCommand(NewLoginCmd(ctx))
	rootCmd.AddCommand(NewNamespaceCmd(ctx))
	// Disable these for now since they don't do anything
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gwa-confg.yaml)")
	// rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "only print results, ideal for CI/CD")
	rootCmd.PersistentFlags().StringVar(&ctx.Host, "host", "", "Set the default host to use for the API")
	rootCmd.PersistentFlags().StringVar(&ctx.Scheme, "scheme", "", "Use to override default https")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
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

		viper.SafeWriteConfig()
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
