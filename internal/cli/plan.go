package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/sanjeevmalagi1/engage-ops/internal/config"
	"github.com/sanjeevmalagi1/engage-ops/internal/core"
	"github.com/sanjeevmalagi1/engage-ops/internal/core/plan"
	"github.com/sanjeevmalagi1/engage-ops/internal/core/state"
)

func newPlanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Show what engageops apply would do, without changing anything",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, _, err := buildPlan(cmd)
			if err != nil {
				return err
			}
			printPlan(cmd.OutOrStdout(), p)
			return nil
		},
	}
}

// buildPlan loads config + state from the shared --dir/--state flags and
// computes a plan, returning the state alongside it so callers (apply,
// destroy) can reuse the same load without reading twice.
func buildPlan(cmd *cobra.Command) (*plan.Plan, *state.State, error) {
	cfg, err := config.Load(configDir)
	if err != nil {
		return nil, nil, err
	}
	st, err := state.Load(statePath)
	if err != nil {
		return nil, nil, err
	}
	p, err := plan.Build(cmd.Context(), cfg, st)
	if err != nil {
		return nil, nil, err
	}
	return p, st, nil
}

func printPlan(w io.Writer, p *plan.Plan) {
	if !p.HasChanges() {
		fmt.Fprintln(w, "No changes. Infrastructure matches the configuration.")
		return
	}

	var creates, updates, deletes int
	for _, c := range p.Changes {
		switch c.Action {
		case core.ActionCreate:
			fmt.Fprintf(w, "  + %s will be created\n", c.Address)
			creates++
		case core.ActionUpdate:
			fmt.Fprintf(w, "  ~ %s will be updated\n", c.Address)
			updates++
		case core.ActionDelete:
			fmt.Fprintf(w, "  - %s will be deleted\n", c.Address)
			deletes++
		}
	}
	fmt.Fprintf(w, "\nPlan: %d to add, %d to change, %d to destroy.\n", creates, updates, deletes)
}
