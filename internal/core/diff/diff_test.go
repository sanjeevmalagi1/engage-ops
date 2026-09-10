package diff

import (
	"testing"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
)

func TestCompute(t *testing.T) {
	addr := core.Address{Provider: "customerio", Type: "segment", Name: "active_users"}

	tests := []struct {
		name    string
		state   core.Attributes
		desired core.Attributes
		want    core.ChangeAction
	}{
		{"both nil is noop", nil, nil, core.ActionNoop},
		{"no state means create", nil, core.Attributes{"name": "Active"}, core.ActionCreate},
		{"no desired means delete", core.Attributes{"name": "Active"}, nil, core.ActionDelete},
		{"identical attrs is noop", core.Attributes{"name": "Active"}, core.Attributes{"name": "Active"}, core.ActionNoop},
		{"changed attrs is update", core.Attributes{"name": "Active"}, core.Attributes{"name": "Inactive"}, core.ActionUpdate},
		{
			"state-only underscore keys are ignored",
			core.Attributes{"_id": 42, "name": "Active"},
			core.Attributes{"name": "Active"},
			core.ActionNoop,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compute(addr, tt.state, tt.desired)
			if got.Action != tt.want {
				t.Errorf("Compute() action = %s, want %s", got.Action, tt.want)
			}
			if got.Address != addr {
				t.Errorf("Compute() address = %v, want %v", got.Address, addr)
			}
		})
	}
}
