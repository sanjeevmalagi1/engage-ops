package provider

import "github.com/sanjeevmalagi1/engage-ops/internal/core"

// Base implements the bookkeeping parts of core.Provider (name + resource
// type lookup) so individual providers only need to supply Configure and
// their Resource map.
type Base struct {
	name      string
	resources map[string]core.Resource
}

// NewBase constructs a Base with the given provider name and resource type
// implementations, keyed by resource type name (e.g. "segment").
func NewBase(name string, resources map[string]core.Resource) Base {
	return Base{name: name, resources: resources}
}

func (b Base) Name() string { return b.name }

func (b Base) Resource(resourceType string) (core.Resource, bool) {
	r, ok := b.resources[resourceType]
	return r, ok
}

func (b Base) ResourceTypes() []string {
	types := make([]string, 0, len(b.resources))
	for t := range b.resources {
		types = append(types, t)
	}
	return types
}
