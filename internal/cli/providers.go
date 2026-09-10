package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sanjeevmalagi1/engage-ops/internal/provider"
)

func newProvidersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "providers",
		Short: "List registered providers and the resource types each supports",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			for _, name := range provider.Names() {
				p, err := provider.New(name)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "%s\n", name)
				for _, rt := range p.ResourceTypes() {
					fmt.Fprintf(out, "  - %s_%s\n", name, rt)
				}
			}
			return nil
		},
	}
}
