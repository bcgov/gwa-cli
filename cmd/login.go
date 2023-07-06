package cmd

import (
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

func NewLoginCmd(ctx *pkg.AppContext) *cobra.Command {

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Log in to your IDIR account",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := pkg.DeviceLogin(ctx)

			if err != nil {
				return err
			}

			return nil
		},
	}
	return loginCmd
}
