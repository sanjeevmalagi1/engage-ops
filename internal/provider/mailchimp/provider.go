// Package mailchimp will implement the engage-ops provider for Mailchimp
// (https://mailchimp.com). Not yet implemented — see customerio for the
// reference provider shape to follow.
package mailchimp

import (
	"context"
	"fmt"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
	"github.com/sanjeevmalagi1/engage-ops/internal/provider"
)

func init() {
	provider.Register("mailchimp", New)
}

type Provider struct {
	provider.Base
}

func New() core.Provider {
	p := &Provider{}
	p.Base = provider.NewBase("mailchimp", map[string]core.Resource{
		"audience": provider.NewStubResource("mailchimp", "audience"),
		"campaign": provider.NewStubResource("mailchimp", "campaign"),
	})
	return p
}

func (p *Provider) Configure(_ context.Context, config core.Attributes) error {
	if _, ok := config["api_key"]; !ok {
		return fmt.Errorf("mailchimp provider: %q is required", "api_key")
	}
	return nil
}
