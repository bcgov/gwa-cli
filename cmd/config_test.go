package cmd

import (
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func SetupConfig(dir string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName(".testing")
	viper.SetConfigType("yaml")
	viper.SafeWriteConfig()
	err := viper.ReadInConfig()
	return err
}

func TestSuccessfulConfigCommands(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name      string
		args      []string
		configKey string
	}{
		{
			name: "set host",
			args: []string{"host", "my.local.dev:8000"},
		},
		{
			name: "set namespace",
			args: []string{"namespace", "ns-sampler"},
		},
		{
			name: "set scheme",
			args: []string{"scheme", "http"},
		},
		{
			name:      "set token",
			args:      []string{"token", "q1w2e3r4t5y6"},
			configKey: "api_key",
		},
		{
			name:      "set host flag",
			args:      []string{"--host", "my.local.dev:8000"},
			configKey: "host",
		},
		{
			name:      "set namespace flag",
			args:      []string{"--namespace", "ns-sampler"},
			configKey: "namespace",
		},
		{
			name:      "set scheme flag",
			args:      []string{"--scheme", "http"},
			configKey: "scheme",
		},
		{
			name:      "set token flag",
			args:      []string{"--token", "q1w2e3r4t5y6"},
			configKey: "api_key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"config", "set"}, tt.args...)
			ctx := &pkg.AppContext{}
			mainCmd := &cobra.Command{
				Use: "gwa",
			}
			cobra.OnInitialize(func() {
				SetupConfig(dir)
			})
			mainCmd.PersistentFlags().StringVar(&ctx.Host, "host", "", "Set the default host to use for the API")
			mainCmd.PersistentFlags().StringVar(&ctx.Scheme, "scheme", "", "Use to override default https")
			mainCmd.PersistentFlags().StringVar(&ctx.Namespace, "namespace", "", "Assign the namespace you would like to use")
			mainCmd.AddCommand(NewConfigCmd(ctx))
			mainCmd.SetArgs(args)
			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})

			assert.Contains(t, out, "Config settings saved", "Expect: %v\nActual: %v\n")
			key := args[2]
			if tt.configKey != "" {
				key = tt.configKey
			}
			assert.Equal(t, tt.args[1], viper.GetString(key))
		})
	}
}
