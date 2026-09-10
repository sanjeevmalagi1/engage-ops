// Package plan computes and applies the set of changes needed to reconcile
// engage-ops state with a desired configuration, mirroring the
// `terraform plan` / `terraform apply` split.
package plan

import (
	"context"
	"fmt"

	"github.com/sanjeevmalagi1/engage-ops/internal/config"
	"github.com/sanjeevmalagi1/engage-ops/internal/core"
	"github.com/sanjeevmalagi1/engage-ops/internal/core/state"
	"github.com/sanjeevmalagi1/engage-ops/internal/provider"
)

// Plan is an ordered list of changes plus the resolved providers/resources
// needed to Apply them, so Apply doesn't have to re-resolve anything.
type Plan struct {
	Changes []*core.Change

	resources map[core.Address]core.Resource
}

// Build configures every provider referenced in cfg and diffs each declared
// resource against st, plus flags any resource present in st but no longer
// declared in cfg for deletion.
func Build(ctx context.Context, cfg *config.Config, st *state.State) (*Plan, error) {
	providers := map[string]core.Provider{}
	for name, providerCfg := range cfg.Providers {
		p, err := provider.New(name)
		if err != nil {
			return nil, err
		}
		if err := p.Configure(ctx, providerCfg); err != nil {
			return nil, fmt.Errorf("configure provider %q: %w", name, err)
		}
		providers[name] = p
	}

	p := &Plan{resources: map[core.Address]core.Resource{}}
	declared := map[core.Address]bool{}

	for _, decl := range cfg.Resources {
		addr := decl.Address()
		declared[addr] = true

		prov, ok := providers[decl.Provider]
		if !ok {
			return nil, fmt.Errorf("%s: provider %q has no configuration block", addr, decl.Provider)
		}
		res, ok := prov.Resource(decl.Type)
		if !ok {
			return nil, fmt.Errorf("%s: provider %q has no resource type %q", addr, decl.Provider, decl.Type)
		}
		p.resources[addr] = res

		change, err := res.Diff(ctx, addr, st.Get(addr), decl.Attributes)
		if err != nil {
			return nil, fmt.Errorf("%s: diff: %w", addr, err)
		}
		p.Changes = append(p.Changes, change)
	}

	for _, key := range st.Addresses() {
		addr, err := core.ParseAddress(key)
		if err != nil {
			return nil, fmt.Errorf("state contains invalid address %q: %w", key, err)
		}
		if declared[addr] {
			continue
		}

		prov, ok := providers[addr.Provider]
		if !ok {
			return nil, fmt.Errorf("%s: tracked in state but provider %q has no configuration block (add it back to delete cleanly)", addr, addr.Provider)
		}
		res, ok := prov.Resource(addr.Type)
		if !ok {
			return nil, fmt.Errorf("%s: tracked in state but provider %q has no resource type %q", addr, addr.Provider, addr.Type)
		}
		p.resources[addr] = res

		change, err := res.Diff(ctx, addr, st.Get(addr), nil)
		if err != nil {
			return nil, fmt.Errorf("%s: diff: %w", addr, err)
		}
		p.Changes = append(p.Changes, change)
	}

	return p, nil
}

// HasChanges reports whether applying this plan would do anything.
func (p *Plan) HasChanges() bool {
	for _, c := range p.Changes {
		if c.Action != core.ActionNoop {
			return true
		}
	}
	return false
}

// Apply executes every non-noop change in order and updates st in place.
// It stops at the first error, leaving st reflecting everything applied so
// far so a re-run of plan/apply picks up where it left off.
func (p *Plan) Apply(ctx context.Context, st *state.State) error {
	for _, change := range p.Changes {
		if change.Action == core.ActionNoop {
			continue
		}

		res, ok := p.resources[change.Address]
		if !ok {
			return fmt.Errorf("%s: internal error: no resource resolved for this change", change.Address)
		}

		attrs, err := res.Apply(ctx, change)
		if err != nil {
			return fmt.Errorf("%s: apply %s: %w", change.Address, change.Action, err)
		}

		if change.Action == core.ActionDelete {
			st.Set(change.Address, nil)
		} else {
			st.Set(change.Address, attrs)
		}
	}
	return nil
}
