package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/bcgov/gwa-cli/pkg"
	"github.com/spf13/cobra"
)

func NewPublishCmd(_ *pkg.AppContext) *cobra.Command {
	var publishCmd = &cobra.Command{
		Deprecated: "Use apply instead.",
		Use:        "publish <type>",
		Short:      "Publish to DS API. Available commands are dataset, issuer and product",
		ValidArgs:  []string{"dataset", "product", "issuer"},
		Args:       cobra.OnlyValidArgs,
		Example: heredoc.Doc(`
    $ gwa publish dataset --input content.yaml
    $ gwa publish product --input content.yaml
    $ gwa publish issuer --input content.yaml
    `),
	}

	return publishCmd
}
