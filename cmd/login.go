package cmd

import (
	"bytes"
	"fmt"
	"strings"

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
		Short: "Log in to your account",
		Long:  `You can login via device login or by using client credentials.`,
		Example: heredoc.Doc(`
      $ gwa login
      $ gwa login --client-id <YOUR_CLIENT_ID> --client-secret <YOUR_CLIENT_SECRET>
    `),
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, _ []string) error {
			var loginMethod string
			if loginFlags.IsClientCredential() {
				err := pkg.ClientCredentialsLogin(ctx, loginFlags.clientId, loginFlags.clientSecret)
				if err != nil {
					return err
				}
				loginMethod = "client credentials"
			} else {
				err := pkg.DeviceLogin(ctx)
				if err != nil {
					return err
				}
				loginMethod = "device login"
			}

			fmt.Printf("\n"+pkg.Checkmark()+" Successfully logged in using %s.\n\n", loginMethod)

			// List available gateways
			buf := new(bytes.Buffer)
			listCmd := GatewayListCmd(ctx, buf)
			if err := listCmd.Execute(); err != nil {
				return fmt.Errorf("Failed to list Gateways: %v", err)
			}

			if buf.Len() > 0 {
				fmt.Println("You have access to the following Gateways and can switch between them with 'gwa config set gateway <gateway-id>':")
				fmt.Println(buf.String())
			} else {
				fmt.Println("You don't have any Gateways. You can create one with 'gwa gateway create'.")
			}

			// Display current gateway
			currentCmd := GatewayCurrentCmd(ctx, buf)
			buf.Reset()
			if err := currentCmd.Execute(); err != nil {
				fmt.Println("No current Gateway set. You can set one with 'gwa config set gateway <gateway-id>'.")
			} else {
				lines := strings.Split(buf.String(), "\n")
				if len(lines) >= 2 {
					fmt.Printf("Using Gateway: \n%s\n", strings.TrimSpace(lines[1]))
				} else {
					fmt.Println("Unable to determine current Gateway.")
				}
			}

			return nil
		}),
	}

	loginCmd.Flags().StringVar(&loginFlags.clientId, "client-id", "", "Your gateway's client ID")
	loginCmd.Flags().StringVar(&loginFlags.clientSecret, "client-secret", "", "Your gateway's client secret")
	loginCmd.MarkFlagsRequiredTogether("client-id", "client-secret")

	return loginCmd
}
