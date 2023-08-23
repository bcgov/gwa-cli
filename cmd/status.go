package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func NewStatusCmd(ctx *pkg.AppContext, buf *bytes.Buffer) *cobra.Command {
	var isJSON bool

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check the status of your configs",
		Example: `  $ gwa status
  $ gwa status --json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ctx.Namespace == "" {
				cmd.SetUsageTemplate(`
You can create a namespace by running:
    $ gwa namespace create`)
				return fmt.Errorf("no namespace has been defined")
			}
			data, err := FetchStatus(ctx)
			if err != nil {
				return err
			}

			if isJSON {
				str, err := json.Marshal(data)
				if err != nil {
					return err
				}
				fmt.Println(string(str))
				return nil
			}

			if len(data) > 0 {
				tbl := table.New("Status", "Name", "Reason", "Upstream")

				if buf != nil {
					tbl.WithWriter(buf)
				}

				for _, item := range data {
					var statusText = pkg.SuccessStyle.Render(item.Status)
					if item.Status == "DOWN" {
						statusText = pkg.ErrorStyle.Render(item.Status)
					}
					tbl.AddRow(statusText, item.Name, item.Reason, item.Upstream)
				}
				tbl.Print()
			} else {
				fmt.Println("You currently do not have any services")
			}

			return nil
		},
	}

	statusCmd.Flags().BoolVar(&isJSON, "json", false, "Output status as a JSON string")

	return statusCmd
}

type StatusJson struct {
	Name     string `json:"name"`
	Upstream string `json:"upstream"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
	Host     string `json:"host"`
	EnvHost  string `json:"env_host"`
}

func FetchStatus(ctx *pkg.AppContext) ([]StatusJson, error) {
	path := fmt.Sprintf("/gw/api/namespaces/%s/services", ctx.Namespace)
	URL, _ := ctx.CreateUrl(path, nil)
	request, err := pkg.NewApiGet[[]StatusJson](ctx, URL)
	if err != nil {
		return nil, err
	}
	response, err := request.Do()
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
