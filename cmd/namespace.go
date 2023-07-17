package cmd

import (
	"fmt"
	"net/http"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

func NewNamespaceCmd(ctx *pkg.AppContext) *cobra.Command {
	var namespaceCmd = &cobra.Command{
		Use:   "namespace",
		Short: "Manage your namespaces",
		Long:  `Longer explanation to come...`,
	}
	namespaceCmd.AddCommand(NamespaceListCmd(ctx))
	namespaceCmd.AddCommand(NamespaceCreateCmd(ctx))
	namespaceCmd.AddCommand(NamespaceCurrentCmd(ctx))
	return namespaceCmd
}

type NamespaceFormData struct {
	name        string
	description string
}

func NamespaceListCmd(ctx *pkg.AppContext) *cobra.Command {
	var listCommand = &cobra.Command{
		Use:   "list",
		Short: "List all your managed namespaces",
		RunE: func(cmd *cobra.Command, _ []string) error {
			URL, _ := ctx.CreateUrl("/ds/api/v2/namespaces", nil)
			r, err := pkg.NewApiGet[[]string](ctx, URL)
			if err != nil {
				return err
			}
			response, err := r.Do()
			if err != nil {
				if response.StatusCode == http.StatusUnauthorized {
					cmd.SetUsageTemplate("\nNext Steps:\nRun gwa login to obtain another auth token")
				}
				return err
			}

			if len(response.Data) <= 0 {
				fmt.Println("You have no namespaces")
			}

			for _, n := range response.Data {
				fmt.Println(n)
			}

			return nil
		},
	}

	return listCommand
}

func NamespaceCreateCmd(ctx *pkg.AppContext) *cobra.Command {
	var namespaceFormData NamespaceFormData
	var createCommand = &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace",
		RunE: func(cmd *cobra.Command, _ []string) error {
			namespace, err := createNamespace(ctx, &namespaceFormData)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}

			// TODO: just returning the name, but determine if a URL would be better
			fmt.Println(namespace)
			return nil
		},
	}
	createCommand.Flags().StringVarP(&namespaceFormData.name, "name", "n", "", "optionally define your own namespace")
	createCommand.Flags().StringVarP(&namespaceFormData.description, "description", "d", "", "optionally add a description")

	return createCommand
}

type NamespaceResult struct {
	Name string `json:"name"`
}

func createNamespace(ctx *pkg.AppContext, data *NamespaceFormData) (string, error) {
	URL, err := ctx.CreateUrl("/ds/api/v2/namespaces", data)
	if err != nil {
		return "", err
	}
	r, err := pkg.NewApiPost[NamespaceResult](ctx, URL, nil)
	if err != nil {
		return "", err
	}

	response, err := r.Do()
	if err != nil {
		return "", err
	}

	return response.Data.Name, nil
}

func NamespaceCurrentCmd(ctx *pkg.AppContext) *cobra.Command {
	var currentCmd = &cobra.Command{
		Use:   "current",
		Short: "Display the current namespace",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ctx.Namespace == "" {
				cmd.SetUsageTemplate(`
You can create a namespace by running:
    $ gwa namespace create`)
				return fmt.Errorf("no namespace has been defined")
			}

			fmt.Println(ctx.Namespace)
			return nil
		},
	}
	return currentCmd
}
