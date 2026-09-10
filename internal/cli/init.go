package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const starterConfig = `# providers.yml — credentials for each engagement platform you manage.
# engageops does not expand ${VAR} placeholders (yet — see README.md
# roadmap), so keep real keys out of version control: use a gitignored
# local file or generate this file from a secrets manager before running.
providers:
  customerio:
    app_api_key: "REPLACE_ME"

# resources.yml — desired-state declarations. Run "engageops plan" to see
# what would change, then "engageops apply" to make it happen.
resources:
  - provider: customerio
    type: segment
    name: active_users
    attributes:
      name: "Active Users"
`

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Scaffold a starter engageops.yml config in the target directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(configDir, 0o755); err != nil {
				return fmt.Errorf("create config dir: %w", err)
			}

			path := filepath.Join(configDir, "engageops.yml")
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("%q already exists, remove it first if you want to regenerate", path)
			}

			if err := os.WriteFile(path, []byte(starterConfig), 0o644); err != nil {
				return fmt.Errorf("write %q: %w", path, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Wrote %s\n\nNext steps:\n  1. Fill in provider credentials\n  2. engageops validate\n  3. engageops plan\n  4. engageops apply\n", path)
			return nil
		},
	}
}
