// Package twilio will implement the engage-ops provider for Twilio
// (https://www.twilio.com), primarily Messaging/Verify services. Not yet
// implemented — see customerio for the reference provider shape to follow.
package twilio

import (
	"context"
	"fmt"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
	"github.com/sanjeevmalagi1/engage-ops/internal/provider"
)

func init() {
	provider.Register("twilio", New)
}

type Provider struct {
	provider.Base
}

func New() core.Provider {
	p := &Provider{}
	p.Base = provider.NewBase("twilio", map[string]core.Resource{
		"messaging_service": provider.NewStubResource("twilio", "messaging_service"),
	})
	return p
}

func (p *Provider) Configure(_ context.Context, config core.Attributes) error {
	if _, ok := config["account_sid"]; !ok {
		return fmt.Errorf("twilio provider: %q is required", "account_sid")
	}
	if _, ok := config["auth_token"]; !ok {
		return fmt.Errorf("twilio provider: %q is required", "auth_token")
	}
	return nil
}
