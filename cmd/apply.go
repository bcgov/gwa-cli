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

type parsedConfig map[string]interface{}

type ApplyOptions struct {
	input string
}

// Takes a dir to locate the input file and returns a slice of each doc contained in the YAML file
func (o *ApplyOptions) Parse(cwd string) ([][]byte, error) {
	filePath := filepath.Join(cwd, o.input)
	ext := filepath.Ext(filePath)
	if ext != ".yaml" && ext != ".yml" {
		return nil, fmt.Errorf("Invalid file type. %s is not a YAML file", o.input)
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	splitDocs, err := pkg.SplitYAML(file)
	if err != nil {
		return nil, err
	}

	return splitDocs, nil
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
	opts := &ApplyOptions{}
	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply gateway resources",
		Long:  "Apply your GatewayService, CredentialIssuer, DraftDataset, and Product resources.  Use the `generate-config` command to see examples of these resources.",
		Args:  cobra.OnlyValidArgs,
		Example: heredoc.Doc(`
$ gwa apply --input gw-config.yaml
    `),
		RunE: func(_ *cobra.Command, _ []string) error {
			yamlDocs, err := opts.Parse(ctx.Cwd)
			if err != nil {
				return err
			}

			counter := &PublishCounter{}
			for _, doc := range yamlDocs {
				// Step 1: Parse the yaml
				config, err := ExtractResourceConfig(doc)
				if err != nil {
					return err
				}
				// Step 2: Get the action type
				action := config.Action()
				// Step 3: Handle the action. Printing is done here
				switch action {
				case "publishGateway":
					fmt.Printf("↑ %s %s", config.Kind, config.Config["name"])
					err := PublishGatewayService(ctx, config.Config)
					if err != nil {
						counter.AddFailed()
						fmt.Print("\r")
						fmt.Printf("%s %s %s\n", pkg.Times(), config.Kind, config.Config["name"])
						break
					}

					counter.AddSuccess()
					fmt.Printf("%s %s %s\n", pkg.Checkmark(), config.Kind, config.Config["name"])
					fmt.Print("\r")
					break

				case "skip":
					counter.AddSkipped()
					fmt.Println(pkg.Indeterminate(), config.Config["name"])
					break

				default:
					fmt.Printf("↑ %s %s", config.Kind, config.Config["name"])
					result, err := PublishResource(ctx, config.Config, action)
					if err != nil {
						counter.AddFailed()
						fmt.Print("\r")
						fmt.Printf("%s %s %s\n", pkg.Times(), config.Kind, config.Config["name"])
						break
					}

					counter.AddSuccess()
					fmt.Print("\r")
					fmt.Printf("%s %s %s: %s\n", pkg.Checkmark(), config.Kind, config.Config["name"], result)
					break
				}
			}

			fmt.Println()
			fmt.Println(counter.Print())

			return nil
		},
	}

	applyCmd.Flags().StringVarP(&opts.input, "input", "i", "gw-config.yml", "YAML file containing your configuration")

	return applyCmd
}

type ResourceConfig struct {
	Config map[string]interface{}
	Kind   string
}

func (r *ResourceConfig) Action() string {
	var kindMapper = map[string]string{
		"CredentialIssuer": "issuer",
		"DraftDataset":     "dataset",
		"Product":          "product",
		"Environment":      "environment",
	}

	if slug, ok := kindMapper[r.Kind]; ok {
		return slug
	} else if r.Kind == "GatewayService" {
		return "publishGateway"
	}
	return "skip"
}

// doc is a single yaml document
func ExtractResourceConfig(doc []byte) (*ResourceConfig, error) {
	var result = &ResourceConfig{}
	err := yaml.Unmarshal(doc, &result.Config)
	if err != nil {
		return result, err
	}
	if result.Config["kind"] == nil {
		return result, fmt.Errorf("This config template is not supported")
	}
	result.Kind = result.Config["kind"].(string)
	delete(result.Config, "kind")
	return result, nil
}

type PutResponse struct {
	Status       int
	Result       string
	Reason       string
	Id           string
	OwnedBy      string
	ChildResults string
}

func PublishResource(ctx *pkg.AppContext, doc parsedConfig, arg string) (string, error) {
	body, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	route := fmt.Sprintf("/ds/api/v2/namespaces/%s/%ss", ctx.Namespace, arg)
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

func PublishGatewayService(ctx *pkg.AppContext, doc parsedConfig) error {
	var kongConfig = struct {
		Services []map[string]interface{} `json:"services"`
	}{}

	kongConfig.Services = append([]map[string]interface{}{}, doc)

	body, err := json.Marshal(kongConfig)
	if err != nil {
		return err
	}
	_, err = PublishToGateway(ctx, &PublishGatewayOptions{}, bytes.NewReader(body))
	if err != nil {
		return err
	}

	return err
}
