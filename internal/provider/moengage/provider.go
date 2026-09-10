// Package moengage will implement the engage-ops provider for MoEngage
// (https://moengage.com). Not yet implemented — see customerio for the
// reference provider shape to follow.
package moengage

import (
	"context"
	"fmt"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
	"github.com/sanjeevmalagi1/engage-ops/internal/provider"
)

func init() {
	provider.Register("moengage", New)
}

type Provider struct {
	provider.Base
}

func New() core.Provider {
	p := &Provider{}
	p.Base = provider.NewBase("moengage", map[string]core.Resource{
		"campaign": provider.NewStubResource("moengage", "campaign"),
		"segment":  provider.NewStubResource("moengage", "segment"),
	})
	return p
}

func (p *Provider) Configure(_ context.Context, config core.Attributes) error {
	if _, ok := config["data_api_id"]; !ok {
		return fmt.Errorf("moengage provider: %q is required", "data_api_id")
	}
	if _, ok := config["api_key"]; !ok {
		return fmt.Errorf("moengage provider: %q is required", "api_key")
	}
	return nil
}
