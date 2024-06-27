package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func SetupConfig(dir string) error {
	fileName := ".gwa-config.yaml"
	path := path.Join(dir, fileName)
	configFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer configFile.Close()
	viper.AddConfigPath(dir)
	viper.SetConfigFile(path)
	return nil
}

func TestSuccessfulConfigCommands(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name      string
		args      []string
		configKey string
		expect    string
	}{
		{
			name: "set host",
			args: []string{"set", "host", "my.local.dev:8000"},
		},
		{
			name: "set gateway",
			args: []string{"set", "gateway", "ns-sampler"},
		},
		{
			name: "set scheme",
			args: []string{"set", "scheme", "http"},
		},
		{
			name:      "set token",
			args:      []string{"set", "token", "q1w2e3r4t5y6"},
			configKey: "api_key",
		},
		{
			name:      "set host flag",
			args:      []string{"set", "--host", "my.local.dev:8000"},
			configKey: "host",
		},
		{
			name:      "set gateway flag",
			args:      []string{"set", "--gateway", "ns-sampler"},
			configKey: "gateway",
		},
		{
			name:      "set scheme flag",
			args:      []string{"set", "--scheme", "http"},
			configKey: "scheme",
		},
		{
			name:      "set token flag",
			args:      []string{"set", "--token", "q1w2e3r4t5y6"},
			configKey: "api_key",
		},
		{
			name: "get gateway",
			args: []string{"get", "gateway"},
			expect: `ns-sampler
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"config"}, tt.args...)
			ctx := &pkg.AppContext{}
			mainCmd := &cobra.Command{
				Use: "gwa",
			}
			cobra.OnInitialize(func() {
				SetupConfig(dir)
			})
			mainCmd.AddCommand(NewConfigCmd(ctx))
			mainCmd.SetArgs(args)
			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})

			if tt.args[1] == "set" {
				assert.Contains(t, out, "Config settings saved", "Expect: %v\nActual: %v\n")
				key := args[3]
				if tt.configKey != "" {
					key = tt.configKey
				}
				assert.Equal(t, tt.args[2], viper.GetString(key))
			}
			if tt.args[0] == "get" {
				assert.Equal(t, tt.expect, out)
			}
		})
	}
}

func TestErrorConfigCommands(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name   string
		args   []string
		expect string
	}{
		{
			name:   "key doesn't exist",
			args:   []string{"set", "random", "akasjfowej"},
			expect: "The key <random> is not allowed to be set",
		},
		{
			name:   "no value set",
			args:   []string{"set", "host"},
			expect: "No value was set for [host]",
		},
		{
			name:   "no flag value set",
			args:   []string{"set", "--gateway"},
			expect: "flag needs an argument: --gateway",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"config"}, tt.args...)
			ctx := &pkg.AppContext{}
			mainCmd := &cobra.Command{
				Use: "gwa",
			}
			cobra.OnInitialize(func() {
				SetupConfig(dir)
			})
			mainCmd.AddCommand(NewConfigCmd(ctx))
			mainCmd.SetArgs(args)
			out := capturer.CaptureStderr(func() {
				mainCmd.Execute()
			})

			assert.Contains(t, out, tt.expect)
		})
	}
}
