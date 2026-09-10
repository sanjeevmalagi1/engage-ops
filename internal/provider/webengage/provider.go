// Package webengage will implement the engage-ops provider for WebEngage
// (https://webengage.com). Not yet implemented — see customerio for the
// reference provider shape to follow.
package webengage

import (
	"context"
	"fmt"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
	"github.com/sanjeevmalagi1/engage-ops/internal/provider"
)

func init() {
	provider.Register("webengage", New)
}

type Provider struct {
	provider.Base
}

func New() core.Provider {
	p := &Provider{}
	p.Base = provider.NewBase("webengage", map[string]core.Resource{
		"campaign": provider.NewStubResource("webengage", "campaign"),
		"segment":  provider.NewStubResource("webengage", "segment"),
	})
	return p
}

func (p *Provider) Configure(_ context.Context, config core.Attributes) error {
	if _, ok := config["api_key"]; !ok {
		return fmt.Errorf("webengage provider: %q is required", "api_key")
	}
	return nil
}
