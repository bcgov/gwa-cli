package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewNamespaceCmd(ctx *pkg.AppContext) *cobra.Command {
	var namespaceCmd = &cobra.Command{
		Use:   "namespace",
		Short: "Manage your namespaces",
		Long:  `Namespaces are used to organize your services.`,
	}
	namespaceCmd.AddCommand(NamespaceListCmd(ctx))
	namespaceCmd.AddCommand(NamespaceCreateCmd(ctx))
	namespaceCmd.AddCommand(NamespaceDestroyCmd(ctx))
	namespaceCmd.AddCommand(NamespaceCurrentCmd(ctx))
	return namespaceCmd
}

type NamespaceFormData struct {
	Name        string `json:"name,omitempty" url:"name,omitempty"`
	Description string `json:"displayName,omitempty" url:"description,omitempty"`
}

func NamespaceListCmd(ctx *pkg.AppContext) *cobra.Command {
	var listCommand = &cobra.Command{
		Use:   "list",
		Short: "List all your managed namespaces",
		RunE: func(_ *cobra.Command, _ []string) error {
			URL, _ := ctx.CreateUrl("/ds/api/v2/namespaces", nil)
			r, err := pkg.NewApiGet[[]string](ctx, URL)
			if err != nil {
				return err
			}
			loader := pkg.NewSpinner()
			loader.Start()
			response, err := r.Do()
			if err != nil {
				if response.StatusCode == http.StatusUnauthorized {
					fmt.Println()
					fmt.Println(
						heredoc.Doc(`
              Next Steps:
              Run gwa login to obtain another auth token
            `),
					)
				}
				return err
			}
			loader.Stop()

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
		Example: heredoc.Doc(`
    $ gwa namespace create
    $ gwa namespace create --name my-namespace --description="This is my namespace"
    `),
		RunE: func(_ *cobra.Command, _ []string) error {
			namespace, err := createNamespace(ctx, &namespaceFormData)
			if err != nil {
				return err
			}

			err = setCurrentNamespace(namespace)
			if err != nil {
				return err
			}

			fmt.Println(namespace)
			return nil
		},
	}
	createCommand.Flags().StringVarP(&namespaceFormData.Name, "name", "n", "", "optionally define your own namespace")
	createCommand.Flags().StringVarP(&namespaceFormData.Description, "description", "d", "", "optionally add a description")

	return createCommand
}

type NamespaceResult struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func createNamespace(ctx *pkg.AppContext, data *NamespaceFormData) (string, error) {
	URL, err := ctx.CreateUrl("/ds/api/v2/namespaces", nil)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	r, err := pkg.NewApiPost[NamespaceResult](ctx, URL, bytes.NewBuffer(body))
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
		RunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Namespace == "" {
				fmt.Println(heredoc.Doc(`
You can create a namespace by running:
    $ gwa namespace create
`))
				return fmt.Errorf("no namespace has been defined")
			}

			fmt.Println(ctx.Namespace)
			return nil
		},
	}
	return currentCmd
}

func setCurrentNamespace(ns string) error {
	viper.Set("namespace", ns)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

type NamespaceDestroyOptions struct {
	Force bool `url:"force"`
}

func NamespaceDestroyCmd(ctx *pkg.AppContext) *cobra.Command {
	var destroyOptions NamespaceDestroyOptions
	var destroyCommand = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the current namespace",
		RunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Namespace == "" {
				fmt.Println(heredoc.Doc(`
          A namespace must be set via the config command

          Example:
              $ gwa config set namespace YOUR_NAMESPACE_NAME
          `),
				)
				return fmt.Errorf("No namespace has been set")
			}

			loader := pkg.NewSpinner()
			loader.Start()
			err := destroyNamespace(ctx, &destroyOptions)
			loader.Stop()
			if err != nil {
				return err
			}

			err = setCurrentNamespace("")
			if err != nil {
				return err
			}

			fmt.Println("Namespace destroyed:", ctx.Namespace)
			return nil
		},
	}

	destroyCommand.Flags().BoolVar(&destroyOptions.Force, "force", false, "force deletion")

	return destroyCommand
}

func destroyNamespace(ctx *pkg.AppContext, destroyOptions *NamespaceDestroyOptions) error {
	pathname := fmt.Sprintf("/ds/api/v2/namespaces/%s", ctx.Namespace)
	URL, err := ctx.CreateUrl(pathname, destroyOptions)
	if err != nil {
		return err
	}
	r, err := pkg.NewApiDelete[NamespaceResult](ctx, URL)
	if err != nil {
		return err
	}

	_, err = r.Do()
	if err != nil {
		return err
	}

	return nil
}
