package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bcgov/gwa-cli/cmd"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

func main() {
	ctx := &pkg.AppContext{}
	rootCmd := cmd.NewRootCommand(ctx)
	output := writeDocument(rootCmd)
	fmt.Println(output)
}

func writeDocument(rootCmd *cobra.Command) string {
	var output strings.Builder
	output.WriteString("# GWA CLI Commands\n\n")

	for _, cmd := range rootCmd.Commands() {
		renderCommand(cmd, &output)
		if cmd.HasSubCommands() {
			for _, subCmd := range cmd.Commands() {
				renderCommand(subCmd, &output)
			}
		}
	}

	return output.String()
}

func renderCommand(cmd *cobra.Command, output *strings.Builder) {
	title := cmd.Name()

	if cmd.HasParent() {
		parentName := cmd.Parent().Name()
		if parentName != "gwa" {
			title = fmt.Sprintf("%s.%s", cmd.Parent().Name(), title)
		}
	}

	output.WriteString(fmt.Sprintf("\n\n## %s\n\n", title))
	output.WriteString(fmt.Sprintf("%s\n\n", strings.ReplaceAll(cmd.Long, "\n", "  \n")))
	output.WriteString(fmt.Sprintf("**Usage:** `%s`\n\n", cmd.UseLine()))

	flagUsages := cmd.Flags().FlagUsages()
	if flagUsages != "" {
		output.WriteString("**Flags**\n\n")
		output.WriteString("| Flag | Description |\n")
		output.WriteString("| ----- | ------ |\n")

		flags := strings.Split(flagUsages, "\n")
		for _, f := range flags {
			trimmedString := strings.TrimSpace(f)
			r := regexp.MustCompile(`\s{2,}`)
			result := r.Split(trimmedString, -1)
			if len(result) == 2 {
				output.WriteString(fmt.Sprintf("| `%s` | %s |\n", result[0], result[1]))
			}
		}
		output.WriteString("\n\n")
	}

	if cmd.Example != "" {
		output.WriteString("**Examples**\n\n")
		output.WriteString(fmt.Sprintf("```shell\n%s\n```\n\n", strings.TrimRight(cmd.Example, "\n ")))
	}
}
