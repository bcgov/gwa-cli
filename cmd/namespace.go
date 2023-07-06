package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

func NewNamespaceCmd(ctx *pkg.AppContext) *cobra.Command {
	var namespaceCmd = &cobra.Command{
		Use:   "namespace",
		Short: "Manage your namespaces",
		Long:  `Longer explanation to come...`,
	}
	namespaceCmd.AddCommand(NamespaceCreateCmd(ctx))
	return namespaceCmd
}

type NamespaceFormData struct {
	namespace   string
	description string
}

func NamespaceCreateCmd(ctx *pkg.AppContext) *cobra.Command {
	var namespaceFormData NamespaceFormData
	var createCommand = &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace",
		RunE: func(cmd *cobra.Command, _ []string) error {
			namespace, err := createNamespace(ctx, &namespaceFormData)
			if err != nil {
				cmd.SetUsageTemplate("\n")
				return err
			}

			successMessage := fmt.Sprintf("namespace %s created", namespace)
			fmt.Println(successMessage)
			return nil
		},
	}
	createCommand.Flags().StringVarP(&namespaceFormData.namespace, "namespace", "n", "", "optionally define your own namespace")
	createCommand.Flags().StringVarP(&namespaceFormData.description, "description", "d", "", "optionally add a description")
	return createCommand
}

type NamespaceResult struct {
	Name string `json:"name"`
}

func createNamespace(ctx *pkg.AppContext, data *NamespaceFormData) (string, error) {
	client := http.Client{}

	URL, err := ctx.CreateUrl("/ds/api/v2/namespaces", data)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest(http.MethodPost, URL, nil)
	if err != nil {
		return "", err
	}
	bearer := fmt.Sprintf("bearer %s", ctx.ApiKey)
	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusOK {
		var result NamespaceResult
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		json.Unmarshal(body, &result)
		return result.Name, err
	}

	var result ApiErrorResponse
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	json.Unmarshal(body, &result)
	errMessage := strings.Join([]string{result.Message, result.Details.Item.Message}, " ")
	return "", errors.New(errMessage)
}

type ApiErrorResponse struct {
	Message string `json:"message"`
	Details struct {
		Item struct {
			Message string `json:"message"`
		} `json:"d0"`
	} `json:"details"`
}
