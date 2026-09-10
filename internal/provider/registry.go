// Package provider holds the provider registry and the built-in engagement
// platform integrations. Each subpackage (customerio, webengage, moengage,
// twilio, mailchimp) implements core.Provider and registers itself here.
package provider

import (
	"fmt"
	"sort"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
)

// Factory constructs a fresh, unconfigured provider instance.
type Factory func() core.Provider

var registry = map[string]Factory{}

// Register makes a provider factory available under name. Called from each
// provider subpackage's init().
func Register(name string, factory Factory) {
	registry[name] = factory
}

// New instantiates a registered provider by name.
func New(name string) (core.Provider, error) {
	factory, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider %q (known providers: %v)", name, Names())
	}
	return factory(), nil
}

// Names returns all registered provider names, sorted.
func Names() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
