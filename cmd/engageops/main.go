// Command engageops is the CLI entrypoint. See internal/cli for command
// wiring and internal/core for the provider/plan/apply model.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sanjeevmalagi1/engage-ops/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
