package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var kindMapper = map[string]string{
	"CredentialIssuer": "issuer",
	"DraftDataset":     "dataset",
	"Product":          "product",
	"Environment":      "environment",
}

// After splitting a config, these are the possible types
type Resource struct {
	Kind   string
	Config map[string]interface{}
}

func (r *Resource) GetAction() string {
	if slug, ok := kindMapper[r.Kind]; ok {
		return slug
	}
	return ""
}

type GatewayService struct {
	Config []map[string]interface{}
}

type Skipped struct {
	Name string
	Kind string
}

// Input struct
type ApplyOptions struct {
	cwd    string
	input  string
	output []interface{}
}

// Takes a dir to locate the input file and returns a slice of each doc contained in the YAML file
func (o *ApplyOptions) Parse() error {
	var gatewayService = GatewayService{}

	filePath := filepath.Join(o.cwd, o.input)
	ext := filepath.Ext(filePath)
	if ext != ".yaml" && ext != ".yml" {
		return fmt.Errorf("Invalid file type. %s is not a YAML file", o.input)
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	splitDocs, err := pkg.SplitYAML(file)
	if err != nil {
		return err
	}

	// Inputs can have multiple configs, so collect separately here
	for _, doc := range splitDocs {
		var parsed map[string]interface{}
		err := yaml.Unmarshal(doc, &parsed)
		if err != nil {
			return err
		}
		if parsed["kind"] == nil {
			return fmt.Errorf("This config template is not supported")
		}
		kind := parsed["kind"].(string)
		delete(parsed, "kind")

		if kind == "GatewayService" {
			gatewayService.Config = append(gatewayService.Config, parsed)
		} else {
			if _, ok := kindMapper[kind]; ok {
				o.output = append(o.output, Resource{
					Kind:   kind,
					Config: parsed,
				})
			} else {
				skipped := Skipped{
					Kind: kind,
				}
				if name, ok := parsed["name"].(string); ok {
					skipped.Name = name
				}
				o.output = append(o.output, skipped)
			}
		}
	}
	o.output = append(o.output, gatewayService)
	return nil
}

type PublishCounter struct {
	Success int
	Failed  int
	Skipped int
}

func (p *PublishCounter) AddSkipped() {
	p.Skipped += 1
}
func (p *PublishCounter) AddFailed() {
	p.Failed += 1
}
func (p *PublishCounter) AddSuccess() {
	p.Success += 1
}

func (p *PublishCounter) Print() string {
	total := p.Success + p.Failed
	return fmt.Sprintf("%d/%d Published, %d Skipped", p.Success, total, p.Skipped)
}

func NewApplyCmd(ctx *pkg.AppContext) *cobra.Command {
	opts := &ApplyOptions{
		cwd: ctx.Cwd,
	}
	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply gateway resources",
		Long:  "Apply your GatewayService, CredentialIssuer, DraftDataset, and Product resources.  Use the `generate-config` command to see examples of these resources.",
		Args:  cobra.OnlyValidArgs,
		Example: heredoc.Doc(`
$ gwa apply --input gw-config.yaml
    `),
		RunE: func(_ *cobra.Command, _ []string) error {
			err := opts.Parse()
			if err != nil {
				return err
			}

			counter := &PublishCounter{}

			for _, config := range opts.output {
				switch c := config.(type) {
				case GatewayService:
					fmt.Println()
					fmt.Printf("↑ Publishing Gateway Services")
					res, err := PublishGatewayService(ctx, c.Config)
					if err != nil {
						counter.AddFailed()
						fmt.Print("\r")
						fmt.Printf("%s Gateway Services publish failed\n", pkg.Times())
						break
					}

					counter.AddSuccess()
					fmt.Printf("%s Gateway Services published\n", pkg.Checkmark())
					fmt.Println(res.Results)
					fmt.Print("\r")
					break

				case Skipped:
					counter.AddSkipped()
					fmt.Printf("%s [%s] %s\n", pkg.Indeterminate(), c.Kind, c.Name)
					break

				case Resource:
					fmt.Printf("↑ [%s] %s", c.Kind, c.Config["name"])
					result, err := PublishResource(ctx, c.Config, c.GetAction())
					if err != nil {
						counter.AddFailed()
						fmt.Print("\r")
						fmt.Printf("%s [%s] %s failed\n", pkg.Times(), c.Kind, c.Config["name"])
						break
					}

					counter.AddSuccess()
					fmt.Print("\r")
					fmt.Printf("%s [%s] %s: %s\n", pkg.Checkmark(), c.Kind, c.Config["name"], result)
					break
				}
			}

			fmt.Println()
			fmt.Println(counter.Print())

			return nil
		},
	}

	applyCmd.Flags().StringVarP(&opts.input, "input", "i", "gw-config.yml", "YAML file containing your configuration")
	applyCmd.MarkFlagRequired("input")

	return applyCmd
}

type PutResponse struct {
	Status       int
	Result       string
	Reason       string
	Id           string
	OwnedBy      string
	ChildResults string
}

func PublishResource(ctx *pkg.AppContext, doc map[string]interface{}, arg string) (string, error) {
	body, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	route := fmt.Sprintf("/ds/api/%s/namespaces/%s/%ss", ctx.ApiVersion, ctx.Namespace, arg)
	URL, _ := ctx.CreateUrl(route, nil)
	request, err := pkg.NewApiPut[PutResponse](ctx, URL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	res, err := request.Do()
	if err != nil {
		return "", err
	}

	return res.Data.Result, nil
}

func PublishGatewayService(ctx *pkg.AppContext, doc []map[string]interface{}) (PublishGatewayResponse, error) {
	var kongConfig = struct {
		Services []map[string]interface{} `json:"services"`
	}{}

	kongConfig.Services = doc

	body, err := json.Marshal(kongConfig)
	if err != nil {
		return PublishGatewayResponse{}, err
	}
	res, err := PublishToGateway(ctx, &PublishGatewayOptions{}, bytes.NewReader(body))
	if err != nil {
		return PublishGatewayResponse{}, err
	}

	return res, err
}
