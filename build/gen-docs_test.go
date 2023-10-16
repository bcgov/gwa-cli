package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestWriteDocument(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:  "gwa",
		Long: "This is a long title",
	}
	configCmd := &cobra.Command{
		Use: "config",
	}
	configGetCmd := &cobra.Command{
		Use: "get",
	}
	configCmd.AddCommand(configGetCmd)
	rootCmd.AddCommand(configCmd)
	output := writeDocument(rootCmd)
	expect := "# GWA CLI Commands\n\nThis is a long title\n\n## config\n\n**Usage:** `gwa config`\n\n\n### config.get\n\n**Usage:** `gwa config get`\n\n"
	assert.Equal(t, expect, output)
}

func TestCommandNoParentNoFlags(t *testing.T) {
	tests := []struct {
		name   string
		cmd    func() *cobra.Command
		expect string
	}{
		{
			name: "no flags",
			cmd: func() *cobra.Command {
				return &cobra.Command{
					Use:  "ping",
					Long: "This command pings a server",
				}
			},
			expect: `
## ping

**Usage:** ` + "`ping`" + `

This command pings a server

`,
		},
		{
			name: "with flags",
			cmd: func() *cobra.Command {
				cmd := &cobra.Command{
					Use:  "ping",
					Long: "This command pings a server",
				}
				cmd.Flags().String("port", "", "set a port")
				return cmd
			},
			expect: fmt.Sprintf(`
## ping

**Usage:** %s

This command pings a server

**Flags**

| Flag | Description |
| ----- | ------ |
| %s | set a port |


`, "`ping [flags]`", "`--port string`"),
		},
		{
			name: "with examples",
			cmd: func() *cobra.Command {
				cmd := &cobra.Command{
					Use:  "ping",
					Long: "This command pings a server",
					Example: heredoc.Doc(`
          $ ping
          $ ping --port 1234
          `),
				}
				cmd.Flags().String("port", "", "set a port")
				return cmd
			},
			expect: fmt.Sprintf(`
## ping

**Usage:** %s

This command pings a server

**Flags**

| Flag | Description |
| ----- | ------ |
| %s | set a port |


**Examples**

%s

`, "`ping [flags]`", "`--port string`", "```shell\n$ ping\n$ ping --port 1234\n```"),
		},
	}
	for _, tt := range tests {
		var output strings.Builder
		renderCommand(tt.cmd(), &output)
		assert.Equal(t, tt.expect, output.String())
	}
}
