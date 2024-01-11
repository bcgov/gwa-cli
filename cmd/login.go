package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

type LoginFlags struct {
	clientId     string
	clientSecret string
}

func (l *LoginFlags) IsClientCredential() bool {
	return l.clientId != "" && l.clientSecret != ""
}

// TODO: Instead of printing from the auth service's methods, use a goroutine and
// post back status updates to this function to keep in line with other methods
func NewLoginCmd(ctx *pkg.AppContext) *cobra.Command {
	loginFlags := &LoginFlags{}

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Log in to your IDIR account",
		Long:  `You can login via device login or by using client credentials.`,
		Example: heredoc.Doc(`
      $ gwa login
      $ gwa login --client-id <YOUR_CLIENT_ID> --client-secret <YOUR_CLIENT_SECRET>
    `),
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, _ []string) error {
			if loginFlags.IsClientCredential() {
				err := pkg.ClientCredentialsLogin(ctx, loginFlags.clientId, loginFlags.clientSecret)
				if err != nil {
					return err
				}
			} else {
				err := pkg.DeviceLogin(ctx)
				if err != nil {
					return err
				}
			}

			fmt.Println(pkg.Checkmark(), pkg.PrintSuccess("Successfully logged in"))

			return nil
		}),
	}

	loginCmd.Flags().StringVar(&loginFlags.clientId, "client-id", "", "Your gateway's client ID")
	loginCmd.Flags().StringVar(&loginFlags.clientSecret, "client-secret", "", "Your gateway's client secret")
	loginCmd.MarkFlagsRequiredTogether("client-id", "client-secret")

	return loginCmd
}
