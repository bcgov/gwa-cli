package cmd

import (
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

func NewInit(_ *pkg.AppContext) *cobra.Command {
	var initCmd = &cobra.Command{
		Deprecated: ".env files are no longer used, see config command",
		Use:        "init",
		Short:      "Generates a .env file in the current working directory.",
		Long: `
Generates a .env file in the current working directory.

To create and work with configurations you don't require CLIENT_ID or CLIENT_SECRET, but to make any API requests you will`,
	}

	return initCmd
}
