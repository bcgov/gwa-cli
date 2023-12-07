package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type OutputFlags struct {
	Json bool
	Yaml bool
}

type GetFilters struct {
	Org string
}

func NewGetCmd(ctx *pkg.AppContext, buf *bytes.Buffer) *cobra.Command {
	var outputOptions = new(OutputFlags)
	var filters = new(GetFilters)
	var validArgs = []string{"datasets", "issuers", "organizations", "org-units", "products"}
	var getCmd = &cobra.Command{
		Use:   "get [type] <flags>",
		Short: fmt.Sprintf("Get gateway resources.  Retrieve a table of %s.", pkg.ArgumentsSliceToString(validArgs, "or")),
		Example: heredoc.Doc(`
      $ gwa get datasets
      $ gwa get datasets --json
      $ gwa get datasets --yaml
    `),
		ValidArgs: validArgs,
		Args:      cobra.OnlyValidArgs,
		RunE: pkg.WrapError(ctx, func(_ *cobra.Command, args []string) error {
			pkg.Info(fmt.Sprintf("Namespace: %s", ctx.Namespace))

			if len(args) == 0 {
				return fmt.Errorf("Must provide an argument of %s to get command", pkg.ArgumentsSliceToString(validArgs, "or"))
			}

			if ctx.Namespace == "" {
				return fmt.Errorf("no namespace selected")
			}

			req := NewRequest(ctx, args[0], filters)
			err := req.Fetch()
			if err != nil {
				return err
			}

			// JSON/YAML outputs
			if outputOptions.Json {
				json, err := json.Marshal(req.Data)
				if err != nil {
					return err
				}
				fmt.Println(string(json))
				return nil
			}

			if outputOptions.Yaml {
				yaml, err := yaml.Marshal(req.Data)
				if err != nil {
					return err
				}
				fmt.Println(string(yaml))
				return nil
			}

			// Pretty-print table
			cols := make([]interface{}, len(req.TableHeaders))
			for i, c := range req.TableHeaders {
				cols[i] = c
			}
			tbl := table.New(cols...)

			switch req.Layout {
			case Basic:
				PrintBasic(req.Data, tbl)
				break
			case OrgUnits:
				PrintOrgUnits(req.Data, tbl)
				break
			case Issuers:
				PrintIssuers(req.Data, tbl)
				break
			case Products:
				PrintProducts(req.Data, tbl)
				break
			}

			if buf != nil {
				tbl.WithWriter(buf)
			}
			tbl.Print()
			return nil
		}),
	}

	getCmd.Flags().StringVar(&filters.Org, "org", "", "Organization to filter results by")
	getCmd.Flags().BoolVar(&outputOptions.Json, "json", false, "Return output as JSON")
	getCmd.Flags().BoolVar(&outputOptions.Yaml, "yaml", false, "Return output as YAML")
	getCmd.MarkFlagsMutuallyExclusive("json", "yaml")
	return getCmd
}

// Enum to test against different accessors
type TableLayout int

const (
	Basic TableLayout = iota
	OrgUnits
	Issuers
	Products
)

type Getter struct {
	Ctx          *pkg.AppContext
	Data         interface{}
	Url          string
	Layout       TableLayout
	TableHeaders []string
}

func (g *Getter) Fetch() error {
	req, err := pkg.NewApiGet[interface{}](g.Ctx, g.Url)
	if err != nil {
		return err
	}
	data, err := req.Do()
	if err != nil {
		return err
	}
	g.Data = data.Data
	return nil
}

func NewRequest(ctx *pkg.AppContext, operator string, filters *GetFilters) *Getter {
	// Parse the URL type
	var path string
	switch operator {
	case "datasets":
		path = fmt.Sprintf("/ds/api/%s/namespaces/%s/directory", ctx.ApiVersion, ctx.Namespace)
		break
	case "organizations":
		path = fmt.Sprintf("/ds/api/%s/organizations", ctx.ApiVersion)
		break
	case "org-units":
		path = fmt.Sprintf("/ds/api/%s/organizations/%s", ctx.ApiVersion, filters.Org)
		break
	default:
		path = fmt.Sprintf("/ds/api/%s/namespaces/%s/%s", ctx.ApiVersion, ctx.Namespace, operator)
	}
	url, _ := ctx.CreateUrl(path, nil)

	// Populate the table headers and layout
	var tableHeaders []string
	var tableLayout TableLayout
	switch operator {
	case "datasets":
		fallthrough
	case "organizations":
		tableHeaders = []string{"Name", "Title"}
		tableLayout = TableLayout(Basic)
		break
	case "org-units":
		tableHeaders = []string{"Name", "Title"}
		tableLayout = TableLayout(OrgUnits)
		break
	case "issuers":
		tableHeaders = []string{"Name", "Flow", "Mode", "Owner"}
		tableLayout = TableLayout(Issuers)
		break
	case "products":
		tableHeaders = []string{"Name", "App ID", "Environments"}
		tableLayout = TableLayout(Products)
		break
	}

	return &Getter{
		Ctx:          ctx,
		Url:          url,
		TableHeaders: tableHeaders,
		Layout:       tableLayout,
	}
}

// TODO: Probably could move these into the struct some how
// NOTE: This is a bit verbose because the goal is to keep the request type as agnostic as possible, so for now the response traversing is "handled in post"
func PrintBasic(data interface{}, tbl table.Table) {
	switch d := data.(type) {
	case []interface{}:
		for _, item := range d {
			if x, ok := item.(map[string]interface{}); ok {
				tbl.AddRow(x["name"], x["title"])
			}
		}
		break
	}
}

func PrintOrgUnits(data interface{}, tbl table.Table) {
	if org, ok := data.(map[string]interface{}); ok {
		if orgUnits := org["orgUnits"].([]interface{}); ok {
			for _, item := range orgUnits {
				if orgUnit, ok := item.(map[string]interface{}); ok {
					tbl.AddRow(orgUnit["name"], orgUnit["title"])
				}
			}
		}
	}
}

func PrintIssuers(data interface{}, tbl table.Table) {
	switch d := data.(type) {
	case []interface{}:
		for _, issuer := range d {
			if i, ok := issuer.(map[string]interface{}); ok {
				tbl.AddRow(i["name"], i["flow"], i["mode"], i["owner"])
			}
		}
	}
}

func PrintProducts(data interface{}, tbl table.Table) {
	switch d := data.(type) {
	case []interface{}:
		for _, product := range d {
			var totalEnvironments int
			if p, ok := product.(map[string]interface{}); ok {
				if envs, ok := p["environments"].([]interface{}); ok {
					totalEnvironments = len(envs)
				} else {
					totalEnvironments = 0
				}
				tbl.AddRow(p["name"], p["appId"], totalEnvironments)
			}
		}
	}
}
