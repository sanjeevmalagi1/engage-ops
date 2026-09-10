package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newApplyCommand() *cobra.Command {
	var autoApprove bool

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the changes required to reach the desired config",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, st, err := buildPlan(cmd)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			printPlan(out, p)
			if !p.HasChanges() {
				return nil
			}

			if !autoApprove {
				fmt.Fprint(out, "\nDo you want to perform these actions? Only 'yes' will be accepted to approve.\n\nEnter a value: ")
				reader := bufio.NewReader(cmd.InOrStdin())
				line, _ := reader.ReadString('\n')
				if strings.TrimSpace(line) != "yes" {
					fmt.Fprintln(out, "Apply cancelled.")
					return nil
				}
			}

			if err := p.Apply(cmd.Context(), st); err != nil {
				// Persist whatever succeeded before the failure so state
				// stays in sync with reality and the next apply can resume.
				_ = st.Save(statePath)
				return err
			}

			if err := st.Save(statePath); err != nil {
				return err
			}
			fmt.Fprintln(out, "\nApply complete.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "skip interactive approval before applying")
	return cmd
}
