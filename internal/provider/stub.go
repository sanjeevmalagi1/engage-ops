package provider

import (
	"context"
	"fmt"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
)

// stubResource satisfies core.Resource for provider/resource-type pairs that
// are registered (so `engageops providers` lists them and config validates)
// but not yet wired up to a real API. Diff still works against local state
// so `plan` is meaningful; Apply and Read fail loudly instead of silently
// no-opping.
type stubResource struct {
	providerName, resourceType string
}

// NewStubResource is exported so provider packages that haven't implemented
// a given resource type yet can register a clear placeholder.
func NewStubResource(providerName, resourceType string) core.Resource {
	return &stubResource{providerName: providerName, resourceType: resourceType}
}

func (s *stubResource) Diff(_ context.Context, addr core.Address, state, desired core.Attributes) (*core.Change, error) {
	action := core.ActionNoop
	switch {
	case state == nil && desired != nil:
		action = core.ActionCreate
	case state != nil && desired == nil:
		action = core.ActionDelete
	case state != nil && desired != nil:
		action = core.ActionUpdate
	}
	return &core.Change{Address: addr, Action: action, Before: state, After: desired}, nil
}

func (s *stubResource) Apply(context.Context, *core.Change) (core.Attributes, error) {
	return nil, fmt.Errorf("%s_%s: not yet implemented, see CONTRIBUTING.md to add a provider", s.providerName, s.resourceType)
}

func (s *stubResource) Read(context.Context, core.Attributes) (core.Attributes, error) {
	return nil, fmt.Errorf("%s_%s: not yet implemented, see CONTRIBUTING.md to add a provider", s.providerName, s.resourceType)
}
