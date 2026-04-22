package loader

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type File struct {
	Name string `yaml:"name"`
	Role string `yaml:"role"`
}

type CoreSpec struct {
	Name        string `yaml:"name"`
	Kind        string `yaml:"kind"`
	Description string `yaml:"description"`
	Files       []File `yaml:"files"`
}

type PackageSpec struct {
	Name        string `yaml:"name"`
	Kind        string `yaml:"kind"`
	Description string `yaml:"description"`
	Files       []File `yaml:"files"`
}

// SubPackageSpec represents a zoom-level-2 package declared with a "/" key (e.g. "verify/golang").
type SubPackageSpec struct {
	Key         string
	Parent      string `yaml:"parent"`
	Kind        string `yaml:"kind"`
	Description string `yaml:"description"`
	Files       []File `yaml:"files"`
}

type Spore struct {
	App         string
	Description string
	Language    string
	Core        CoreSpec
	Packages    []PackageSpec
	SubPackages []SubPackageSpec
}

func Load(path string) (*Spore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parsing spore.yaml: %w", err)
	}
	if len(root.Content) == 0 {
		return nil, fmt.Errorf("empty spore.yaml")
	}

	mapping := root.Content[0]

	// Decode known top-level fields via struct tags.
	type rawSpore struct {
		App         string        `yaml:"app"`
		Description string        `yaml:"description"`
		Language    string        `yaml:"language"`
		Core        CoreSpec      `yaml:"core"`
		Packages    []PackageSpec `yaml:"packages"`
	}
	var raw rawSpore
	if err := mapping.Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding spore.yaml: %w", err)
	}

	spore := &Spore{
		App:         raw.App,
		Description: raw.Description,
		Language:    raw.Language,
		Core:        raw.Core,
		Packages:    raw.Packages,
	}

	// Sub-packages are top-level keys that contain "/" (e.g. "verify/golang").
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		key := mapping.Content[i].Value
		if !strings.Contains(key, "/") {
			continue
		}
		var sub SubPackageSpec
		if err := mapping.Content[i+1].Decode(&sub); err != nil {
			return nil, fmt.Errorf("decoding sub-package %s: %w", key, err)
		}
		sub.Key = key
		spore.SubPackages = append(spore.SubPackages, sub)
	}

	return spore, nil
}
