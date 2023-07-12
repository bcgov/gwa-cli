package cmd

import (
	"fmt"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

func NewLoginCmd(ctx *pkg.AppContext) *cobra.Command {
	var clientId string
	var clientSecret string

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Log in to your IDIR account",
		Long: `You can login via device login or by using client credentials

To use device login, simply run the command like so:
    $ gwa login

To use your credentials you must supply both a client-id and client-secret:
    $ gwa login --client-id <YOUR_CLIENT_ID> --client-secret <YOUR_CLIENT_SECRET>
    `,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SilenceUsage = true
			if clientId != "" && clientSecret != "" {
				err := pkg.ClientCredentialsLogin(ctx, clientId, clientSecret)
				if err != nil {
					return err
				}
			} else {
				err := pkg.DeviceLogin(ctx)
				if err != nil {
					return err
				}
			}

			fmt.Println("Logged in")

			return nil
		},
	}

	loginCmd.Flags().StringVar(&clientId, "client-id", "", "Your gateway's client ID")
	loginCmd.Flags().StringVar(&clientSecret, "client-secret", "", "Your gateway's client secret")
	loginCmd.MarkFlagsRequiredTogether("client-id", "client-secret")

	return loginCmd
}
