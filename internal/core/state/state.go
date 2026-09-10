// Package state persists the last-known-applied attributes for every
// resource engage-ops manages, so future plans can diff against reality
// instead of assuming a clean slate. It plays the same role as
// terraform.tfstate.
package state

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
)

const currentVersion = 1

// State is the on-disk (and in-memory working) representation of every
// resource instance engage-ops currently manages.
type State struct {
	Version   int                        `json:"version"`
	Resources map[string]core.Attributes `json:"resources"` // keyed by Address.String()
}

// New returns an empty state, used the first time engage-ops runs in a
// directory that has no state file yet.
func New() *State {
	return &State{Version: currentVersion, Resources: map[string]core.Attributes{}}
}

// Load reads state from path. A missing file is not an error: it returns a
// fresh empty State, matching Terraform's "no state yet" behavior.
func Load(path string) (*State, error) {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read state %q: %w", path, err)
	}

	var s State
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, fmt.Errorf("parse state %q: %w", path, err)
	}
	if s.Resources == nil {
		s.Resources = map[string]core.Attributes{}
	}
	return &s, nil
}

// Save writes state to path as pretty-printed JSON, creating or truncating
// the file as needed.
func (s *State) Save(path string) error {
	s.Version = currentVersion
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return fmt.Errorf("write state %q: %w", path, err)
	}
	return nil
}

// Get returns the stored attributes for addr, or nil if untracked.
func (s *State) Get(addr core.Address) core.Attributes {
	return s.Resources[addr.String()]
}

// Set records addr's attributes, or removes the entry when attrs is nil.
func (s *State) Set(addr core.Address, attrs core.Attributes) {
	if attrs == nil {
		delete(s.Resources, addr.String())
		return
	}
	s.Resources[addr.String()] = attrs
}

// Addresses returns every address currently tracked in state.
func (s *State) Addresses() []string {
	keys := make([]string, 0, len(s.Resources))
	for k := range s.Resources {
		keys = append(keys, k)
	}
	return keys
}
