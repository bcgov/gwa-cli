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

type payload map[string]interface{}

type OutputFlags struct {
	Json bool
	Yaml bool
}

type RequestFilters struct {
	Org string
}

// requests feature 4 different types of URLs
// 3 different table column

func NewGetCmd(ctx *pkg.AppContext, buf *bytes.Buffer) *cobra.Command {
	var outputOptions = new(OutputFlags)
	var filters = new(RequestFilters)
	var validArgs = []string{"datasets", "issuers", "organizations", "organization", "products"}
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
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Must provide an argument of %s to get command", pkg.ArgumentsSliceToString(validArgs, "or"))
			}

			if ctx.Namespace == "" {
				return fmt.Errorf("no namespace selected")
			}

			getter := NewRequest(ctx, args[0], filters)
			err := getter.Fetch()
			if err != nil {
				return err
			}

			// JSON/YAML outputs
			if outputOptions.Json {
				json, err := json.Marshal(getter.Data)
				if err != nil {
					return err
				}
				fmt.Println(string(json))
				return nil
			}

			if outputOptions.Yaml {
				yaml, err := yaml.Marshal(getter.Data)
				if err != nil {
					return err
				}
				fmt.Println(string(yaml))
				return nil
			}

			// Pretty-print table
			cols := make([]interface{}, len(getter.TableHeaders))
			for i, c := range getter.TableHeaders {
				cols[i] = c
			}
			tbl := table.New(cols...)
			switch getter.Type {
			case Basic:
				PrintBasic(getter.Data, tbl)
				break
			case OrgUnits:
				PrintOrgUnits(getter.Data, tbl)
				break
			case Issuers:
				PrintIssuers(getter.Data, tbl)
				break
			case Products:
				PrintProducts(getter.Data, tbl)
				break
			}

			if buf != nil {
				tbl.WithWriter(buf)
			}
			tbl.Print()
			return nil
		},
	}

	getCmd.Flags().StringVar(&filters.Org, "name", "", "Organization to filter results by")
	getCmd.Flags().BoolVar(&outputOptions.Json, "json", false, "Return output as JSON")
	getCmd.Flags().BoolVar(&outputOptions.Yaml, "yaml", false, "Return output as YAML")
	getCmd.MarkFlagsMutuallyExclusive("json", "yaml")
	return getCmd
}

// Enum to test against different accessors
type TableType int

const (
	Basic TableType = iota
	OrgUnits
	Issuers
	Products
)

type Getter struct {
	Ctx          *pkg.AppContext
	Data         interface{}
	Url          string
	Type         TableType
	TableHeaders []string
	TableColumns []string
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

func NewRequest(ctx *pkg.AppContext, operator string, filters *RequestFilters) *Getter {
	// Parse the URL type
	var path string
	switch operator {
	case "datasets":
		path = fmt.Sprintf("/ds/api/%s/namespaces/%s/directory", ctx.ApiVersion, ctx.Namespace)
		break
	case "organizations":
		path = fmt.Sprintf("/ds/api/%s/organizations", ctx.ApiVersion)
		break
	case "organization":
		path = fmt.Sprintf("/ds/api/%s/organizations/%s", ctx.ApiVersion, filters.Org)
		break
	default:
		path = fmt.Sprintf("/ds/api/%s/namespaces/%s/%s", ctx.ApiVersion, ctx.Namespace, operator)
	}
	url, _ := ctx.CreateUrl(path, nil)

	// Populate the table headers and columns
	var tableHeaders []string
	var tableType TableType
	switch operator {
	case "datasets":
		fallthrough
	case "organizations":
		tableHeaders = []string{"Name", "Title"}
		tableType = TableType(Basic)
		break
	case "organization":
		tableHeaders = []string{"Name", "Title"}
		tableType = TableType(OrgUnits)
		break
	case "issuers":
		tableHeaders = []string{"Name", "Flow", "Mode", "Owner"}
		tableType = TableType(Issuers)
		break
	case "products":
		tableHeaders = []string{"Name", "App ID", "Environments"}
		tableType = TableType(Products)
		break
	}

	return &Getter{
		Ctx:          ctx,
		Url:          url,
		TableHeaders: tableHeaders,
		Type:         tableType,
	}
}

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
