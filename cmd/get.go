package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Output interface {
	PrintTable() string
}

type OutputFlags struct {
	Json bool
	Yaml bool
}

func (o *OutputFlags) Print(data Container, buf *bytes.Buffer) error {
	if o.Json {
		output, err := data.JSON()
		if err != nil {
			return err
		}
		fmt.Println(output)
		return nil
	}
	if o.Yaml {
		output, err := data.YAML()
		if err != nil {
			return err
		}
		fmt.Println(output)
		return nil
	}
	tbl := table.New(data.GetHeaders()...)
	if buf != nil {
		tbl.WithWriter(buf)
	}
	data.GetRows(tbl)
	tbl.Print()
	return nil
}

func NewGetCmd(ctx *pkg.AppContext, buf *bytes.Buffer) *cobra.Command {
	var outputOptions = new(OutputFlags)
	var getCmd = &cobra.Command{
		Use:   "get [type] <flags>",
		Short: "Retrieve a table of a namespace's datasets, issuers and products",
		Example: `$ gwa get datasets
$ gwa get datasets --json
$ gwa get datasets --yaml`,
		ValidArgs: []string{"datasets", "issuers", "products"},
		Args:      cobra.OnlyValidArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Must provide an argument of datasets, issuers or products to get command")
			}
			if ctx.Namespace == "" {
				return fmt.Errorf("no namespace selected")
			}
			data, err := CreateAction(ctx, args[0])
			if err != nil {
				return err
			}
			outputOptions.Print(data, buf)
			return nil
		},
	}

	getCmd.Flags().BoolVar(&outputOptions.Json, "json", false, "Return output as JSON")
	getCmd.Flags().BoolVar(&outputOptions.Yaml, "yaml", false, "Return output as YAML")
	getCmd.MarkFlagsMutuallyExclusive("json", "yaml")
	return getCmd
}

func CreateAction(ctx *pkg.AppContext, operator string) (Container, error) {
	var path = fmt.Sprintf("/ds/api/v2/namespaces/%s/%s", ctx.Namespace, operator)
	if operator == "datasets" {
		path = fmt.Sprintf("/ds/api/v2/namespaces/%s/directory", ctx.Namespace)
	}
	url, _ := ctx.CreateUrl(path, nil)

	result, err := FetchData(ctx, url, operator)
	if err != nil {
		return result, err
	}

	return result, nil
}

type Container interface {
	GetHeaders() []interface{}
	GetRows(tbl table.Table)
	JSON() (string, error)
	YAML() (string, error)
}

type Datasets struct {
	Data []Dataset
}

func (d *Datasets) GetHeaders() []interface{} {
	return []interface{}{"Name", "Title"}
}

func (d *Datasets) GetRows(tbl table.Table) {
	for _, dataset := range d.Data {
		tbl.AddRow(dataset.Name, dataset.Title)
	}
}

func (d *Datasets) JSON() (string, error) {
	o, err := json.Marshal(d.Data)
	if err != nil {
		return "", err
	}

	return string(o), nil
}

func (d *Datasets) YAML() (string, error) {
	o, err := yaml.Marshal(d.Data)
	if err != nil {
		return "", err
	}

	return string(o), nil
}

type Issuers struct {
	Data []Issuer
}

func (d *Issuers) GetHeaders() []interface{} {
	return []interface{}{"Name", "Flow", "Mode", "Owner"}
}

func (d *Issuers) GetRows(tbl table.Table) {
	for _, issuer := range d.Data {
		tbl.AddRow(issuer.Name, issuer.Flow, issuer.Mode, issuer.Owner)
	}
}

func (d *Issuers) JSON() (string, error) {
	o, err := json.Marshal(d.Data)
	if err != nil {
		return "", err
	}

	return string(o), nil
}

func (d *Issuers) YAML() (string, error) {
	o, err := yaml.Marshal(d.Data)
	if err != nil {
		return "", err
	}

	return string(o), nil
}

type Products struct {
	Data []Product
}

func (d *Products) GetHeaders() []interface{} {
	return []interface{}{"Name", "AppId", "Environments"}
}

func (d *Products) GetRows(tbl table.Table) {
	for _, product := range d.Data {
		tbl.AddRow(product.Name, product.AppId, len(product.Environments))
	}
}

func (d *Products) JSON() (string, error) {
	o, err := json.Marshal(d.Data)
	if err != nil {
		return "", err
	}

	return string(o), nil
}

func (d *Products) YAML() (string, error) {
	o, err := yaml.Marshal(d.Data)
	if err != nil {
		return "", err
	}

	return string(o), nil
}

func MakeRequest[T any](ctx *pkg.AppContext, url string) (T, error) {
	var result T
	req, err := pkg.NewApiGet[T](ctx, url)
	if err != nil {
		return result, err
	}
	res, err := req.Do()
	if err != nil {
		return result, err
	}

	return res.Data, nil
}

func FetchData(ctx *pkg.AppContext, url string, operator string) (Container, error) {
	var result Container
	switch operator {
	case "datasets":
		data, err := MakeRequest[[]Dataset](ctx, url)
		if err != nil {
			return result, err
		}
		result = &Datasets{
			Data: data,
		}
		return result, nil
	case "issuers":
		data, err := MakeRequest[[]Issuer](ctx, url)
		if err != nil {
			return result, err
		}
		result = &Issuers{
			Data: data,
		}
		return result, nil
	case "products":
		data, err := MakeRequest[[]Product](ctx, url)
		if err != nil {
			return result, err
		}
		result = &Products{
			Data: data,
		}
		return result, nil
	}

	return result, nil
}

type Dataset struct {
	DownloadAudience  string   `json:"download_audience,omitempty" yaml:"download_audience,omitempty"`
	LicenseTitle      string   `json:"license_title,omitempty" yaml:"license_title,omitempty"`
	Name              string   `json:"name,omitempty" yaml:"name,omitempty"`
	Notes             string   `json:"notes,omitempty" yaml:"notes,omitempty"`
	Organization      string   `json:"organization,omitempty"`
	OrganizationUnit  string   `json:"organizationUnit,omitempty" yaml:"organizationUnit,omitempty"`
	RecordPublishDate string   `json:"record_publish_date,omitempty" yaml:"record_publish_date,omitempty"`
	SecurityClass     string   `json:"security_class,omitempty" yaml:"security_class,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	Title             string   `json:"title,omitempty"`
	ViewAudience      string   `json:"view_audience,omitempty" yaml:"view_audience,omitempty"`
}

type Product struct {
	AppId        string `json:"appId,omitempty" yaml:"appId"`
	Environments []struct {
		Active   bool   `json:"active,omitempty"`
		AppId    string `json:"appId,omitempty" yaml:"appId"`
		Approval bool   `json:"approval,omitempty"`
		Flow     string `json:"flow,omitempty"`
		Name     string `json:"name,omitempty"`
	} `json:"environments,omitempty"`
	Name string `json:"name,omitempty"`
}

type Issuer struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	Flow                string `json:"flow"`
	ClientAuthenticator string `json:"clientAuthenicator" yaml:"clientAuthenticator"`
	Mode                string `json:"mode"`
	EnvironmentDetails  []struct {
		ClientId           string `json:"clientId,omitempty" yaml:"clientId,omitempty"`
		ClientRegistration string `json:"clientRegistration,omitempty" yaml:"clientRegistration,omitempty"`
		ClientSecret       string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
		Environment        string `json:"environment,omitempty" yaml:"environment,omitempty"`
		IssuerUrl          string `json:"issuerUrl,omitempty" yaml:"issuerUrl,omitempty"`
	} `json:"environmentDetails" yaml:"environmentDetails"`
	Owner string `json:"owner"`
}
