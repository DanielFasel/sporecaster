package loader

import (
	"fmt"
	"os"

	"github.com/DanielFasel/sporecaster/internal/spore"
	"gopkg.in/yaml.v3"
)

func Load(path string) (*spore.Spore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var raw struct {
		App         string          `yaml:"app"`
		Description string          `yaml:"description"`
		Language    string          `yaml:"language"`
		Core        spore.Core      `yaml:"core"`
		Packages    []spore.Package `yaml:"packages"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if raw.App == "" {
		return nil, fmt.Errorf("%s: app is required", path)
	}
	if raw.Language == "" {
		return nil, fmt.Errorf("%s: language is required", path)
	}

	return &spore.Spore{
		App:         raw.App,
		Description: raw.Description,
		Language:    raw.Language,
		Core:        raw.Core,
		Packages:    raw.Packages,
	}, nil
}
