// Package config loads engage-ops YAML configuration: provider credentials
// and the desired-state resource declarations, analogous to a Terraform
// root module made of .tf files.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
)

// ResourceDecl is one `resources:` entry in YAML.
type ResourceDecl struct {
	Provider   string          `yaml:"provider"`
	Type       string          `yaml:"type"`
	Name       string          `yaml:"name"`
	Attributes core.Attributes `yaml:"attributes"`
}

func (r ResourceDecl) Address() core.Address {
	return core.Address{Provider: r.Provider, Type: r.Type, Name: r.Name}
}

// file mirrors the on-disk YAML shape of a single config file.
type file struct {
	Providers map[string]core.Attributes `yaml:"providers"`
	Resources []ResourceDecl             `yaml:"resources"`
}

// Config is the merged result of every *.yml/*.yaml file in a directory.
type Config struct {
	Providers map[string]core.Attributes
	Resources []ResourceDecl
}

// Load reads and merges every YAML file directly inside dir (non-recursive,
// mirroring how Terraform treats a module directory).
func Load(dir string) (*Config, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read config dir %q: %w", dir, err)
	}

	var paths []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext == ".yml" || ext == ".yaml" {
			paths = append(paths, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(paths)

	if len(paths) == 0 {
		return nil, fmt.Errorf("no .yml/.yaml files found in %q", dir)
	}

	cfg := &Config{Providers: map[string]core.Attributes{}}
	seen := map[core.Address]string{}

	for _, path := range paths {
		var f file
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %q: %w", path, err)
		}
		if err := yaml.Unmarshal(raw, &f); err != nil {
			return nil, fmt.Errorf("parse %q: %w", path, err)
		}

		for name, attrs := range f.Providers {
			cfg.Providers[name] = attrs
		}

		for _, r := range f.Resources {
			if r.Provider == "" || r.Type == "" || r.Name == "" {
				return nil, fmt.Errorf("%q: resource declarations require provider, type, and name", path)
			}
			addr := r.Address()
			if prev, ok := seen[addr]; ok {
				return nil, fmt.Errorf("%q: duplicate resource %s (first declared in %q)", path, addr, prev)
			}
			seen[addr] = path
			cfg.Resources = append(cfg.Resources, r)
		}
	}

	return cfg, nil
}
