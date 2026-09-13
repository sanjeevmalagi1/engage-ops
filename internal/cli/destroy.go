package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sanjeevmalagi1/engage-ops/internal/config"
	"github.com/sanjeevmalagi1/engage-ops/internal/core/plan"
	"github.com/sanjeevmalagi1/engage-ops/internal/core/state"
)

func newDestroyCommand() *cobra.Command {
	var autoApprove bool

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Delete every resource engageops currently tracks in state",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, statePathResolved := resolvePaths(cmd)
			cfg, err := config.Load(dir)
			if err != nil {
				return err
			}
			st, err := state.Load(statePathResolved)
			if err != nil {
				return err
			}

			// A destroy plan is a normal plan against a config with no
			// resource declarations: every tracked resource is flagged
			// for deletion, but providers stay configured.
			emptyCfg := &config.Config{Providers: cfg.Providers}
			p, err := plan.Build(cmd.Context(), emptyCfg, st)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			printPlan(out, p)
			if !p.HasChanges() {
				return nil
			}

			if !autoApprove {
				fmt.Fprint(out, "\nDo you really want to destroy all resources?\n  engageops will delete all resources shown above.\n  Only 'yes' will be accepted to confirm.\n\nEnter a value: ")
				reader := bufio.NewReader(cmd.InOrStdin())
				line, _ := reader.ReadString('\n')
				if strings.TrimSpace(line) != "yes" {
					fmt.Fprintln(out, "Destroy cancelled.")
					return nil
				}
			}

			if err := p.Apply(cmd.Context(), st); err != nil {
				_ = st.Save(statePathResolved)
				return err
			}

			if err := st.Save(statePathResolved); err != nil {
				return err
			}
			fmt.Fprintln(out, "\nDestroy complete.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "skip interactive approval before destroying")
	return cmd
}
