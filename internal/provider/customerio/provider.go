// Package customerio implements the engage-ops provider for Customer.io
// (https://customer.io). It is the reference provider implementation:
// other platforms (WebEngage, MoEngage, Twilio, Mailchimp) should follow the
// same shape.
package customerio

import (
	"context"
	"fmt"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
	"github.com/sanjeevmalagi1/engage-ops/internal/provider"
)

func init() {
	provider.Register("customerio", New)
}

const (
	defaultAppAPIURL = "https://api.customer.io/v1"
)

type Provider struct {
	provider.Base

	siteID string
	apiKey string
	appKey string
	apiURL string
	client httpClient
}

func New() core.Provider {
	p := &Provider{apiURL: defaultAppAPIURL, client: defaultHTTPClient()}
	p.Base = provider.NewBase("customerio", map[string]core.Resource{
		"segment": &segmentResource{provider: p},
	})
	return p
}

func (p *Provider) Configure(_ context.Context, config core.Attributes) error {
	appKey, _ := config["app_api_key"].(string)
	if appKey == "" {
		return fmt.Errorf("customerio provider: %q is required", "app_api_key")
	}
	p.appKey = appKey

	if siteID, ok := config["site_id"].(string); ok {
		p.siteID = siteID
	}
	if apiKey, ok := config["api_key"].(string); ok {
		p.apiKey = apiKey
	}
	if apiURL, ok := config["app_api_url"].(string); ok && apiURL != "" {
		p.apiURL = apiURL
	}
	return nil
}
