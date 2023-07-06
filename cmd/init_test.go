package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func TestInitCommands(t *testing.T) {
	ctx := &pkg.AppContext{
		Cwd: t.TempDir(),
	}
	envFilePath := filepath.Join(ctx.Cwd, ".env")

	tests := []struct {
		name       string
		beforeEach func()
		expect     string
		args       []string
	}{
		{
			name:       "successful init command",
			beforeEach: nil,
			expect:     ".env created",
			args:       []string{"--namespace=ns-sampler", "--client-id=asdf", "--client-secret=2338"},
		},
		{
			name: ".env file already exists",
			beforeEach: func() {
				os.WriteFile(envFilePath, []byte("\t\n"), 0644)
			},
			expect: "Error: .env already exists",
			args:   []string{"--namespace=ns-sampler", "--client-id=asdf", "--client-secret=2342"},
		},
		{
			name:       "short namespace",
			beforeEach: nil,
			expect:     "namespace must be between 5 and 15 characters long",
			args:       []string{"--namespace=ns", "--client-id=asdf", "--client-secret=2342"},
		},
		{
			name:       "long namespace",
			beforeEach: nil,
			expect:     "namespace must be between 5 and 15 characters long",
			args:       []string{"--namespace=abcdefghijklmnopqrstuvwxyz", "--client-id=asdf", "--client-secret=2342"},
		}, {
			name:       "invalid characters",
			beforeEach: nil,
			expect:     "namespace can only contain alphanumeric characters and -",
			args:       []string{"--namespace=!@#$%^&*()", "--client-id=asdf", "--client-secret=2341"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.beforeEach != nil {
				tt.beforeEach()
			}
			args := append([]string{"init"}, tt.args...)
			mainCmd := &cobra.Command{
				Use: "gwa",
			}
			mainCmd.AddCommand(NewInit(ctx))
			mainCmd.SetArgs(args)

			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})

			assert.Contains(t, out, tt.expect, "Expect: %v\nActual: %v\n", tt.expect, out)
		})
	}
}

func Test_createConfig(t *testing.T) {
	cwd := t.TempDir()
	envFilePath := filepath.Join(cwd, ".env")
	tests := []struct {
		name   string
		hasEnv bool
	}{
		{
			name:   "empty env",
			hasEnv: false,
		},
		{
			name:   "empty env",
			hasEnv: true,
		},
	}
	for _, tt := range tests {
		if tt.hasEnv {
			os.WriteFile(envFilePath, []byte("\t\n"), 0644)
		}
		opts := &initOptions{
			namespace:    "ns.sampler",
			clientId:     "123",
			clientSecret: "456",
			dev:          true,
			prod:         false,
			test:         false,
			dataCenter:   "calgary",
			apiVersion:   2,
			cwd:          cwd,
		}
		err := createConfig(opts)

		if tt.hasEnv {
			assert.Error(t, err, "is error")
		}
	}
}
