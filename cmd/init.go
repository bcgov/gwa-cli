package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/basvdlei/envfile"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

type initOptions struct {
	namespace    string
	clientId     string
	clientSecret string
	dataCenter   string
	apiVersion   int
	dsApiVersion int
	dev          bool
	prod         bool
	test         bool
	cwd          string
}

func (c *initOptions) getEnv() string {
	env := "dev"
	if c.dev == true {
		env = "dev"
	} else if c.prod == true {
		env = "prod"
	} else if c.test == true {
		env = "test"
	}

	return env
}

func (c *initOptions) validate() error {
	namespaceLength := len(c.namespace)
	if namespaceLength < 5 || namespaceLength > 15 {
		return fmt.Errorf("namespace must be between 5 and 15 characters long")
	}

	namespacePattern := "^[a-zA-Z0-9-]{5,15}$"
	regex := regexp.MustCompile(namespacePattern)
	if regex.MatchString(c.namespace) == false {
		return fmt.Errorf("namespace can only contain alphanumeric characters and -")
	}

	return nil
}

func NewInit(ctx *pkg.AppContext) *cobra.Command {
	opts := &initOptions{
		cwd: ctx.Cwd,
	}
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Generates a .env file in the current working directory.",
		Long: `Generates a .env file in the current working directory.

To create and work with configurations you don't require CLIENT_ID or CLIENT_SECRET, but to make any API requests you will`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := opts.validate(); err != nil {
				return err
			}
			err := runCmd(cmd, opts)

			if err != nil {
				return err
			}
			return nil
		},
	}

	initCmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "", "The namespace of you routes collection (required)")
	initCmd.Flags().StringVar(&opts.clientId, "client-id", "", "Namespace's Client ID from API Services Portal")
	initCmd.Flags().StringVar(&opts.clientSecret, "client-secret", "", "Namespace's Client Secret from API Services Portal")
	initCmd.Flags().StringVarP(&opts.dataCenter, "data-center", "d", "calgary", "Target a particular data centre")
	initCmd.Flags().BoolVarP(&opts.dev, "dev", "D", false, "Set the environment as development")
	initCmd.Flags().BoolVarP(&opts.prod, "prod", "P", false, "Set the environment as production")
	initCmd.Flags().BoolVarP(&opts.test, "test", "T", false, "Set the environment as test")
	initCmd.Flags().IntVar(&opts.apiVersion, "api-version", 2, "Set the API version")
	initCmd.Flags().IntVar(&opts.dsApiVersion, "ds-api-version", 2, "Set the Directory API version")
	initCmd.MarkFlagRequired("namespace")
	initCmd.MarkFlagRequired("client-id")
	initCmd.MarkFlagRequired("client-secret")
	initCmd.MarkFlagsMutuallyExclusive("dev", "prod", "test")

	return initCmd
}

func runCmd(cmd *cobra.Command, opts *initOptions) error {
	err := createConfig(opts)

	if err != nil {
		cmd.SetUsageTemplate("try something else")
		return err
	}

	fmt.Println(".env created")
	return nil
}

func createConfig(opts *initOptions) error {
	envPath := filepath.Join(opts.cwd, ".env")

	if _, err := os.Stat(envPath); err == nil {
		return fmt.Errorf(".env already exists")
	}

	settings := struct {
		NameSpace    string `env:"GWA_NAMESPACE"`
		ClientId     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
		Env          string `env:"GWA_ENV"`
		DataCenter   string `env:"DATA_CENTER"`
		ApiVersion   string `env:"API_VERSION"`
		DsApiVersion string `env:"DS_API_VERSION"`
	}{
		NameSpace:    opts.namespace,
		ClientId:     opts.clientId,
		ClientSecret: opts.clientId,
		Env:          opts.getEnv(),
		DataCenter:   opts.dataCenter,
		ApiVersion:   fmt.Sprint(opts.apiVersion),
		DsApiVersion: fmt.Sprint(opts.dsApiVersion),
	}

	out, err := envfile.Marshal(settings)
	if err != nil {
		return err
	}

	if err := os.WriteFile(envPath, out, 0666); err != nil {
		return err
	}

	return nil
}
