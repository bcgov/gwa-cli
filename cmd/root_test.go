package cmd

import (
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// fakeConfigValue stands in for a config file value for the rest of the test.
// "" means not set. It overrides any config another test loaded.
func fakeConfigValue(t *testing.T, key string, value string) {
	t.Helper()
	viper.Set(key, value)
	// A nil override is skipped by viper, which clears it
	t.Cleanup(func() { viper.Set(key, nil) })
}

// runWithGlobalFlags parses args with a stand-in subcommand, the way cobra
// parses them for real commands, then applies the config. No network calls.
func runWithGlobalFlags(t *testing.T, ctx *pkg.AppContext, args []string) {
	t.Helper()
	rootCmd := NewRootCommand(ctx)
	probe := &cobra.Command{
		Use:                "probe",
		PersistentPostRunE: func(_ *cobra.Command, _ []string) error { return nil },
		RunE: func(_ *cobra.Command, _ []string) error {
			applyConfigToContext(ctx, rootCmd)
			return nil
		},
	}
	rootCmd.AddCommand(probe)
	rootCmd.SetArgs(append([]string{"probe"}, args...))
	assert.NoError(t, rootCmd.Execute())
}

func TestGlobalFlagsOverrideConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        map[string]string
		args          []string
		expectGateway string
		expectScheme  string
	}{
		{
			name:          "config only",
			config:        map[string]string{"gateway": "cfg-gw", "scheme": "http"},
			args:          []string{},
			expectGateway: "cfg-gw",
			expectScheme:  "http",
		},
		{
			name:          "flags override config",
			config:        map[string]string{"gateway": "cfg-gw", "scheme": "http"},
			args:          []string{"--gateway", "flag-gw", "--scheme", "https"},
			expectGateway: "flag-gw",
			expectScheme:  "https",
		},
		{
			name:          "flags with no config",
			config:        map[string]string{},
			args:          []string{"--gateway=flag-gw", "--scheme=http"},
			expectGateway: "flag-gw",
			expectScheme:  "http",
		},
		{
			name:          "neither set",
			config:        map[string]string{},
			args:          []string{},
			expectGateway: "",
			expectScheme:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeConfigValue(t, "gateway", tt.config["gateway"])
			fakeConfigValue(t, "scheme", tt.config["scheme"])

			ctx := &pkg.AppContext{}
			runWithGlobalFlags(t, ctx, tt.args)
			assert.Equal(t, tt.expectGateway, ctx.Gateway)
			assert.Equal(t, tt.expectScheme, ctx.Scheme)
		})
	}
}

func TestHostFlagOverridesConfig(t *testing.T) {
	// Release builds compile in api.gov.bc.ca as the --host default
	compiledHost := "api.gov.bc.ca"
	tests := []struct {
		name       string
		configHost string
		args       []string
		expectHost string
	}{
		{
			name:       "config host replaces the compiled default",
			configHost: "dev.example.com",
			args:       []string{},
			expectHost: "dev.example.com",
		},
		{
			name:       "--host set to the compiled default still wins",
			configHost: "dev.example.com",
			args:       []string{"--host", compiledHost},
			expectHost: compiledHost,
		},
		{
			name:       "--host set to another host wins",
			configHost: "dev.example.com",
			args:       []string{"--host", "other.example.com"},
			expectHost: "other.example.com",
		},
		{
			name:       "compiled default with no config or flag",
			configHost: "",
			args:       []string{},
			expectHost: compiledHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeConfigValue(t, "host", tt.configHost)

			ctx := &pkg.AppContext{ApiHost: compiledHost}
			runWithGlobalFlags(t, ctx, tt.args)
			assert.Equal(t, tt.expectHost, ctx.ApiHost)
		})
	}
}
