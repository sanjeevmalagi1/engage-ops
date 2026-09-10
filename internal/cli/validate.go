package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sanjeevmalagi1/engage-ops/internal/config"
	"github.com/sanjeevmalagi1/engage-ops/internal/core/plan"
	"github.com/sanjeevmalagi1/engage-ops/internal/core/state"
)

func newValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Check that the config directory parses, providers are configured, and resource types are known",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configDir)
			if err != nil {
				return err
			}
			// Diffing against empty state performs the same
			// provider/resource-type resolution `plan` would, without
			// requiring a state file or making any API calls.
			if _, err := plan.Build(cmd.Context(), cfg, state.New()); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "config is valid: %d provider(s), %d resource(s)\n", len(cfg.Providers), len(cfg.Resources))
			return nil
		},
	}
}
