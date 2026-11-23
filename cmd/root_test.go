package cmd

import (
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

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
			for _, key := range []string{"gateway", "scheme"} {
				k := key
				if v, ok := tt.config[k]; ok {
					viper.Set(k, v)
				}
				// A nil override is skipped by viper, which clears it
				t.Cleanup(func() { viper.Set(k, nil) })
			}

			ctx := &pkg.AppContext{}
			rootCmd := NewRootCommand(ctx)
			// Stand-in subcommand so the global flags are parsed the way
			// cobra parses them for real commands, without network calls.
			probe := &cobra.Command{
				Use:                "probe",
				PersistentPostRunE: func(_ *cobra.Command, _ []string) error { return nil },
				RunE: func(_ *cobra.Command, _ []string) error {
					applyConfigToContext(ctx, rootCmd)
					return nil
				},
			}
			rootCmd.AddCommand(probe)
			rootCmd.SetArgs(append([]string{"probe"}, tt.args...))

			assert.NoError(t, rootCmd.Execute())
			assert.Equal(t, tt.expectGateway, ctx.Gateway)
			assert.Equal(t, tt.expectScheme, ctx.Scheme)
		})
	}
}
