// Package diff provides the generic attribute-map comparison that most
// resource implementations can reuse instead of hand-rolling their own
// create/update/delete detection.
package diff

import (
	"reflect"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
)

// Compute derives a Change for addr by comparing state (nil if the resource
// is not yet tracked) against desired (nil if the resource was removed from
// config). It does not talk to any remote API.
func Compute(addr core.Address, state, desired core.Attributes) *core.Change {
	switch {
	case state == nil && desired == nil:
		return &core.Change{Address: addr, Action: core.ActionNoop}
	case state == nil:
		return &core.Change{Address: addr, Action: core.ActionCreate, After: desired}
	case desired == nil:
		return &core.Change{Address: addr, Action: core.ActionDelete, Before: state}
	case Equal(state, desired):
		return &core.Change{Address: addr, Action: core.ActionNoop, Before: state, After: desired}
	default:
		return &core.Change{Address: addr, Action: core.ActionUpdate, Before: state, After: desired}
	}
}

// Equal reports whether two attribute maps are semantically identical,
// ignoring keys present only in state that are conventionally
// provider-assigned (prefixed with "_", e.g. "_id", "_created_at").
func Equal(state, desired core.Attributes) bool {
	filtered := make(core.Attributes, len(state))
	for k, v := range state {
		if len(k) > 0 && k[0] == '_' {
			continue
		}
		filtered[k] = v
	}
	return reflect.DeepEqual(filtered, desired)
}
