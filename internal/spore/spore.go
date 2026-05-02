package spore

type File struct {
	Name string `yaml:"name"`
	Role string `yaml:"role"`
}

type Core struct {
	Name        string `yaml:"name"`
	Kind        string `yaml:"kind"`
	Description string `yaml:"description"`
	Files       []File `yaml:"files"`
}

type Package struct {
	Name        string   `yaml:"name"`
	Parent      string   `yaml:"parent"`
	Kind        string   `yaml:"kind"`
	Description string   `yaml:"description"`
	Files       []File   `yaml:"files"`
	Imports     []string `yaml:"imports"`
}

type Spore struct {
	App         string
	Description string
	Language    string
	Core        Core
	Packages    []Package
}
