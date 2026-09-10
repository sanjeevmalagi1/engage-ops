// Package cli wires up the engageops command-line interface: config/state
// file resolution shared by every subcommand, plus the terraform-style
// plan/apply/destroy/validate/providers commands.
package cli

import (
	"github.com/spf13/cobra"

	_ "github.com/sanjeevmalagi1/engage-ops/internal/provider/customerio"
	_ "github.com/sanjeevmalagi1/engage-ops/internal/provider/mailchimp"
	_ "github.com/sanjeevmalagi1/engage-ops/internal/provider/moengage"
	_ "github.com/sanjeevmalagi1/engage-ops/internal/provider/twilio"
	_ "github.com/sanjeevmalagi1/engage-ops/internal/provider/webengage"
)

var (
	configDir string
	statePath string
)

// NewRootCommand builds the engageops root cobra command with all
// subcommands attached.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "engageops",
		Short:         "Infrastructure-as-code for customer engagement platforms",
		Long:          "engageops manages resources across engagement platforms (Customer.io, WebEngage,\nMoEngage, Twilio, Mailchimp, ...) from declarative YAML, the way Terraform\nmanages cloud infrastructure.",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	root.PersistentFlags().StringVar(&configDir, "dir", ".", "directory containing engage-ops *.yml config files")
	root.PersistentFlags().StringVar(&statePath, "state", "engageops.tfstate.json", "path to the state file")

	root.AddCommand(
		newValidateCommand(),
		newPlanCommand(),
		newApplyCommand(),
		newDestroyCommand(),
		newProvidersCommand(),
		newInitCommand(),
	)
	return root
}
