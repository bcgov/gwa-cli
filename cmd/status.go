package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func NewStatusCmd(ctx *pkg.AppContext, buf *bytes.Buffer) *cobra.Command {
	var isJSON bool

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check the status of your services configured on the Kong gateway",
		Example: heredoc.Doc(`$ gwa status
  $ gwa status --json`),
		RunE: func(_ *cobra.Command, _ []string) error {
			if ctx.Gateway == "" {
				fmt.Println(heredoc.Doc(`
          You can create a gateway by running:
              $ gwa gateway create
          `),
				)
				return fmt.Errorf("no gateway has been defined")
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
	path := fmt.Sprintf("/gw/api/%s/gateways/%s/services", ctx.ApiVersion, ctx.Gateway)
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
