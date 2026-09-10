// Package core defines the provider and resource abstractions that every
// engagement-platform integration (Customer.io, WebEngage, MoEngage, Twilio,
// Mailchimp, ...) implements. It is the engage-ops equivalent of Terraform's
// provider SDK: a small, stable contract that the plan/apply engine drives
// without knowing anything about a specific platform's API.
package core

import (
	"context"
	"fmt"
	"strings"
)

// Attributes is the generic bag of resource fields as decoded from YAML or
// read back from a remote platform. Providers are responsible for
// interpreting their own keys.
type Attributes map[string]any

// Address uniquely identifies a resource instance within a config, e.g.
// "customerio_segment.active_users".
type Address struct {
	Provider string
	Type     string
	Name     string
}

func (a Address) String() string {
	return a.Provider + "_" + a.Type + "." + a.Name
}

// ParseAddress parses the "provider_type.name" form produced by
// Address.String back into its parts. Provider and resource type names must
// not themselves contain underscores for this round-trip to be unambiguous;
// engage-ops enforces that when providers register.
func ParseAddress(s string) (Address, error) {
	dot := strings.LastIndex(s, ".")
	if dot < 0 {
		return Address{}, fmt.Errorf("address %q missing \".\" separator before resource name", s)
	}
	provType, name := s[:dot], s[dot+1:]

	provider, typ, ok := strings.Cut(provType, "_")
	if !ok {
		return Address{}, fmt.Errorf("address %q missing \"_\" separator between provider and type", s)
	}

	return Address{Provider: provider, Type: typ, Name: name}, nil
}

// ChangeAction is the kind of change a plan step represents.
type ChangeAction string

const (
	ActionNoop   ChangeAction = "noop"
	ActionCreate ChangeAction = "create"
	ActionUpdate ChangeAction = "update"
	ActionDelete ChangeAction = "delete"
)

// Change describes one planned mutation to a single resource instance.
type Change struct {
	Address         Address
	Action          ChangeAction
	Before          Attributes
	After           Attributes
	ReplacedBecause string // non-empty if Update was upgraded to a delete+create
}

// Resource is implemented once per resource type by a provider (e.g. the
// customerio provider implements it for "segment", "campaign", ...).
type Resource interface {
	// Diff compares the desired attributes against the last-known state and
	// returns the change that would bring the remote platform in line.
	// state may be nil when the resource does not yet exist in state.
	Diff(ctx context.Context, addr Address, state, desired Attributes) (*Change, error)

	// Apply performs the change against the remote platform and returns the
	// resulting attributes to persist in state (including any
	// platform-assigned fields such as IDs).
	Apply(ctx context.Context, change *Change) (Attributes, error)

	// Read fetches the current remote attributes for drift detection.
	// Returns nil, nil if the resource no longer exists remotely.
	Read(ctx context.Context, state Attributes) (Attributes, error)
}

// Provider is a named collection of resource types backed by one engagement
// platform's API, plus the credentials/config needed to talk to it.
type Provider interface {
	// Name returns the provider's config key, e.g. "customerio".
	Name() string

	// Configure receives the provider block from YAML (API keys, region,
	// etc.) before any resource operations run.
	Configure(ctx context.Context, config Attributes) error

	// Resource returns the Resource implementation for a given resource
	// type name (e.g. "segment"), or false if unsupported.
	Resource(resourceType string) (Resource, bool)

	// ResourceTypes lists all resource type names this provider supports.
	ResourceTypes() []string
}
