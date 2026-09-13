package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// resolvePaths derives the config directory and state file path for a
// command invocation. --env selects environments/<env> and a per-environment
// state file by default, but an explicitly-set --dir or --state always wins,
// so single-environment usage is unaffected.
func resolvePaths(cmd *cobra.Command) (dir, statePathResolved string) {
	dir, statePathResolved = configDir, statePath
	if envName == "" {
		return dir, statePathResolved
	}
	if !cmd.Flags().Changed("dir") {
		dir = filepath.Join("environments", envName)
	}
	if !cmd.Flags().Changed("state") {
		statePathResolved = fmt.Sprintf("engageops.%s.tfstate.json", envName)
	}
	return dir, statePathResolved
}
